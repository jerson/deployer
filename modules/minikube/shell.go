package minikube

import (
	"context"
	"fmt"
	"os"

	"github.com/jerson/deployer/modules/shell"
)

// ContextKeyBinary ...
const ContextKeyBinary = "minikube_binary"

// Binary ...
var Binary = "minikube"

func binary(ctx context.Context) string {
	binary := ctx.Value(ContextKeyBinary)
	if binary != nil {
		return binary.(string)
	}
	return Binary
}

// Run ...
func Run(ctx context.Context, commands ...string) (output string, err error) {
	return shell.Run(ctx, binary(ctx), commands...)
}

// Settings ...
type Settings struct {
	Name     string
	Driver   string
	Memory   uint
	CPUS     uint
	DiskSize string
}

// Config ...
func Config(ctx context.Context, properties Settings) error {

	if properties.Driver != "none" {
		_, _ = Run(ctx, "config", "set", "profile", properties.Name)
	}
	if properties.CPUS > 0 {
		_, _ = Run(ctx, "config", "set", "cpus", fmt.Sprint(properties.CPUS))
	}
	if properties.Memory > 0 {
		_, _ = Run(ctx, "config", "set", "memory", fmt.Sprint(properties.Memory))
	}
	if properties.Driver != "" {
		_, _ = Run(ctx, "config", "set", "vm-driver", properties.Driver)
	}
	if properties.DiskSize != "" {
		_, _ = Run(ctx, "config", "set", "disk-size", properties.DiskSize)
	}

	return nil
}

// Env ...
func Env(ctx context.Context) (err error) {
	err = os.Setenv("DOCKER_TLS_VERIFY", getEnv(ctx, "DOCKER_TLS_VERIFY"))
	if err != nil {
		return err
	}
	err = os.Setenv("DOCKER_HOST", getEnv(ctx, "DOCKER_HOST"))
	if err != nil {
		return err
	}
	err = os.Setenv("DOCKER_CERT_PATH", getEnv(ctx, "DOCKER_CERT_PATH"))
	if err != nil {
		return err
	}
	err = os.Setenv("MINIKUBE_ACTIVE_DOCKERD", getEnv(ctx, "MINIKUBE_ACTIVE_DOCKERD"))
	if err != nil {
		return err
	}
	return err
}

func getEnv(ctx context.Context, name string) string {
	output, _ := shell.Run(ctx, "sh", "-c", fmt.Sprintf("eval $(%s docker-env) && echo $%s", binary(ctx), name))
	return output
}
