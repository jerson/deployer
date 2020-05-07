package helm

import (
	"context"
	"strings"

	"github.com/jerson/deployer/modules/shell"
	"github.com/jerson/deployer/pkg/deployer"
)

// ContextKeyBinary ...
const ContextKeyBinary = "helm_binary"

// Binary ...
var Binary = "helm"

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

// Init ...
func Init(ctx context.Context) (err error) {
	log := deployer.Log(ctx)
	version, _ := Version(ctx)
	if version < 3 {
		// this is only for helm2
		log.Info("Helm tiler init")
		_, err = Run(ctx, "init")
		if err != nil {
			return nil
		}
	}
	return err
}

// Install ...
func Install(ctx context.Context, name string, commands ...string) (output string, err error) {

	var arguments []string
	arguments = append(arguments, "install")

	version, _ := Version(ctx)
	if version < 3 {
		arguments = append(arguments, []string{"--name", name}...)
	} else {
		arguments = append(arguments, name)
	}

	arguments = append(arguments, commands...)
	return Run(ctx, arguments...)
}

// Version please fix this implementation
func Version(ctx context.Context) (int, error) {

	output, err := Run(ctx, "version")
	if err != nil {
		return 0, err
	}
	if strings.Contains(output, "v2.") {
		return 2, nil
	}
	return 3, nil
}
