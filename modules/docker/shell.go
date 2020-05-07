package docker

import (
	"context"
	"errors"

	"github.com/jerson/deployer/modules/shell"
	"github.com/jerson/deployer/pkg/deployer"
)

// ContextKeyBinary ...
const ContextKeyBinary = "docker_binary"

// Binary ...
var Binary = "docker"

func binary(ctx context.Context) string {
	binary := ctx.Value(ContextKeyBinary)
	if binary != nil {
		return binary.(string)
	}
	return Binary
}

// IsRunning ...
func IsRunning(ctx context.Context, container string) (running bool, err error) {

	output, err := Run(ctx, "container", "inspect", container, "-f", "{{.State.Running}}")
	if err != nil {
		return false, err
	}
	if output == "false" {
		return false, errors.New("container not running")

	}
	return true, nil
}

// WaitIsRunning ...
func WaitIsRunning(ctx context.Context, container string) (running bool, err error) {
	return shell.WaitIsTrue(ctx, func(ctx context.Context) (ready bool, err error) {
		return IsRunning(ctx, container)
	})
}

// Remove ...
func Remove(ctx context.Context, container string) (err error) {

	_, err = Run(context.WithValue(ctx, deployer.ContextKeyDisableSTDOUT, true), "container", "inspect", container, "-f", "{{.State.Running}}")
	if err != nil {
		return nil
	}

	_, err = Run(ctx, "rm", container, "-f")

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
