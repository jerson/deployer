package deployer

import (
	"context"
	"errors"
	"testing"

	"github.com/jerson/deployer/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewDeployerStart(t *testing.T) {

	var func1 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func2 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func1}, installer
	}

	var func3 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func4 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func3, func2}, installer
	}

	deployments := []entities.Deployment{func4}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	err := deploy.InstallAll(ctx)

	printDeployer(ctx, deploy)

	assert.NoError(t, err)

	assert.Equal(t, 4, len(deploy.deployments))
	assert.Equal(t, FunctionName(func3), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func4)]))
}

func TestNewDeployerStartError(t *testing.T) {

	errInstall := errors.New("error install")
	var func1 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func2 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, errInstall
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func1}, installer
	}

	var func3 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func4 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, errInstall
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func3, func2}, installer
	}

	deployments := []entities.Deployment{func4}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	err := deploy.InstallAll(ctx)

	printDeployer(ctx, deploy)

	assert.Error(t, err)

	assert.Equal(t, deploy.deploymentStatus, deploy.DeploymentStatusAll())
	assert.Equal(t, errInstall, deploy.deploymentStatus[FunctionName(func2)].Install.Error)

	assert.Equal(t, 4, len(deploy.deployments))
	assert.Equal(t, FunctionName(func3), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func4)]))
}

func TestNewDeployerStartUninstallError(t *testing.T) {

	errInstall := errors.New("error install")
	errUninstall := errors.New("error uninstall")
	var func1 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, errUninstall
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func2 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, errInstall
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func1}, installer
	}

	var func3 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{}, installer
	}

	var func4 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			uninstall = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			upgrade = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			install = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
			return install, upgrade, uninstall, health
		}
		return []entities.Deployment{func3, func2}, installer
	}

	deployments := []entities.Deployment{func4}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	err := deploy.InstallAll(ctx)

	printDeployer(ctx, deploy)

	assert.Error(t, err)
	assert.Equal(t, deploy.deploymentStatus, deploy.DeploymentStatusAll())
	assert.Equal(t, errInstall, deploy.deploymentStatus[FunctionName(func2)].Install.Error)

	err = deploy.UninstallDeployed(ctx)
	assert.NoError(t, err)
	assert.Equal(t, errUninstall, deploy.deploymentStatus[FunctionName(func1)].Uninstall.Error)

	assert.Equal(t, 4, len(deploy.deployments))
	assert.Equal(t, FunctionName(func3), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func4)]))
}
