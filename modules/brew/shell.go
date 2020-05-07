package brew

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jerson/deployer/modules/shell"
)

// ContextKeyBinary ...
const ContextKeyBinary = "brew_binary"

// Binary ...
var Binary = "brew"

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

// InstallWithLook ...
func InstallWithLook(ctx context.Context, name, binary string) (output string, err error) {

	path := envPath(name)
	err = addToPath(path)
	if err != nil {
		return path, err
	}
	foundPath, _ := exec.LookPath(binary)
	if foundPath == "" {
		err = addToBashEnv(ctx, path)
		if err != nil {
			return path, err
		}
		return Run(ctx, "install", name)
	}

	return path, nil
}

// Install ...
func Install(ctx context.Context, name string) (output string, err error) {
	return InstallWithLook(ctx, name, name)
}

func envPath(name string) string {
	if runtime.GOOS == "darwin" {
		return fmt.Sprintf("/usr/local/opt/%s/bin", name)
	} else if runtime.GOOS == "linux" {
		return fmt.Sprintf("/home/linuxbrew/.linuxbrew/opt/%s/bin", name)
	}
	return ""
}

func addToBashEnv(ctx context.Context, path string) (err error) {
	rcName := bashFileName()
	_, err = shell.Run(ctx, "sh", "-c", fmt.Sprintf(`echo 'export PATH="%s:$PATH"' >> ~/%s`, path, rcName))
	return err
}

func addToPath(dir string) error {
	return os.Setenv("PATH", fmt.Sprintf("%s:%s", dir, os.Getenv("PATH")))
}

func bashFileName() string {
	name := ".bashrc"

	shellName := os.Getenv("SHELL")

	if strings.Contains(shellName, "/zsh") {
		name = ".zshrc"
	}

	return name
}
