package skaffold

import (
	"context"

	"github.com/jerson/deployer/modules/shell"
)

// ContextKeyBinary ...
const ContextKeyBinary = "skaffold_binary"

// Binary ...
var Binary = "skaffold"

func binary(ctx context.Context) string {
	binary := ctx.Value(ContextKeyBinary)
	if binary != nil {
		return binary.(string)
	}
	return Binary
}

// Deploy ...
func Deploy(ctx context.Context, projectDir string) (err error) {
	_, err = RunInDir(ctx, projectDir, "run")
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
