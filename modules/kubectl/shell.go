package kubectl

import (
	"context"
	"errors"
	"fmt"

	"github.com/jerson/deployer/modules/shell"
	"github.com/jerson/deployer/pkg/deployer"

	"github.com/sirupsen/logrus"
)

// ContextKeyBinary ...
const ContextKeyBinary = "kubectl_binary"

// Binary ...
var Binary = "kubectl"

func binary(ctx context.Context) string {
	binary := ctx.Value(ContextKeyBinary)
	if binary != nil {
		return binary.(string)
	}
	return Binary
}

// ContextKeyProject ...
const ContextKeyProject = "kubectl_project"

// Project ...
var Project = "default"

func project(ctx context.Context) string {
	binary := ctx.Value(ContextKeyProject)
	if binary != nil {
		return binary.(string)
	}
	return Project
}

// CreateNamespace ...
func CreateNamespace(ctx context.Context, name string) (err error) {
	_, err = Run(ctx, "create", "namespace", name)
	return err
}

// Exec ...
func Exec(ctx context.Context, pod, command string) (name string, err error) {

	output, err := shell.Run(ctx, "sh", "-c", fmt.Sprintf(`
%s --namespace %s exec %s -- %s
`, binary(ctx), project(ctx), pod, command))

	return output, err
}

// ExecContainer ...
func ExecContainer(ctx context.Context, pod, container, command string) (name string, err error) {

	output, err := shell.Run(ctx, "sh", "-c", fmt.Sprintf(`
%s --namespace %s exec %s --container %s -- %s
`, binary(ctx), project(ctx), pod, container, command))

	return output, err
}

// GetPodReady ...
func GetPodReady(ctx context.Context, labels ...string) (name string, err error) {

	_, err = WaitPodReady(ctx, labels...)
	if err != nil {
		return "", err
	}

	return GetPod(ctx, labels...)
}

// GetPod ...
func GetPod(ctx context.Context, labels ...string) (name string, err error) {

	log := deployer.Log(ctx)
	var arguments []string
	arguments = append(arguments, "--namespace", project(ctx))
	arguments = append(arguments, "get", "pods")

	for _, label := range labels {
		arguments = append(arguments, "-l")
		arguments = append(arguments, label)
	}

	arguments = append(arguments, "-o", "jsonpath={.items[0].metadata.name}")
	output, err := Run(ctx, arguments...)
	if log.Logger.Level >= logrus.DebugLevel {
		// this will add extra break line for better log
		_, _ = log.Logger.Out.Write([]byte("\n"))
	}

	return output, err
}

// WaitPodReady ...
func WaitPodReady(ctx context.Context, labels ...string) (ready bool, err error) {
	return WaitPodReadyWithNamespace(ctx, project(ctx), labels...)
}

// WaitPodReadyWithNamespace ...
func WaitPodReadyWithNamespace(ctx context.Context, namespace string, labels ...string) (ready bool, err error) {
	return shell.WaitIsTrue(ctx, func(ctx context.Context) (ready bool, err error) {
		return PodReadyWithNamespace(ctx, namespace, labels...)
	})

}

// PodReady ...
func PodReady(ctx context.Context, labels ...string) (ready bool, err error) {
	return PodReadyWithNamespace(ctx, project(ctx), labels...)
}

// PodReadyWithNamespace ...
func PodReadyWithNamespace(ctx context.Context, namespace string, labels ...string) (ready bool, err error) {

	log := deployer.Log(ctx)

	var arguments []string
	arguments = append(arguments, "--namespace", namespace)
	arguments = append(arguments, "get", "pods")

	for _, label := range labels {
		arguments = append(arguments, "-l")
		arguments = append(arguments, label)
	}

	arguments = append(arguments, "-o", "jsonpath={..status.conditions[?(@.type==\"Ready\")].status}")
	output, err := Run(ctx, arguments...)
	if log.Logger.Level >= logrus.DebugLevel {
		// this will add extra break line for better log
		_, _ = log.Logger.Out.Write([]byte("\n"))
	}
	if output == "True" {
		return true, nil
	}

	return false, errors.New("pod not ready")
}

// Proxy ...
func Proxy(ctx context.Context, namespace, deployment string, portTarget, portSource int) (err error) {

	var arguments []string
	arguments = append(arguments, "--namespace", namespace, "port-forward")
	arguments = append(arguments, fmt.Sprintf("deployment/%s", deployment))
	arguments = append(arguments, fmt.Sprintf("%d:%d", portTarget, portSource))
	_, err = Run(ctx, arguments...)

	return err
}

// Restart ...
func Restart(ctx context.Context, deployment string) (err error) {

	var arguments []string
	arguments = append(arguments, "--namespace", project(ctx), "rollout", "restart")
	arguments = append(arguments, fmt.Sprintf("deployment/%s", deployment))
	_, err = Run(ctx, arguments...)

	return err
}

// Run ...
func Run(ctx context.Context, commands ...string) (output string, err error) {
	return shell.Run(ctx, binary(ctx), commands...)
}

// RunInDir ...
func RunInDir(ctx context.Context, dir string, commands ...string) (output string, err error) {
	return shell.RunInDir(ctx, dir, binary(ctx), commands...)
}
