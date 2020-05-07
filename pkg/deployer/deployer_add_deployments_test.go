package deployer

import (
	"context"
	"testing"

	"github.com/jerson/deployer/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestAddAllDeploymentsDefault(t *testing.T) {

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
		return []entities.Deployment{func3, func2, func1}, installer
	}

	deployments := []entities.Deployment{func4}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	printDeployer(ctx, deploy)
	assert.Equal(t, 4, len(deploy.deployments))
	assert.Equal(t, FunctionName(func3), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 2, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func4)]))
}

func TestAddAllDeploymentsAlternative(t *testing.T) {

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
		return []entities.Deployment{func2, func3}, installer
	}

	deployments := []entities.Deployment{func4}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	printDeployer(ctx, deploy)
	assert.Equal(t, 4, len(deploy.deployments))
	assert.Equal(t, FunctionName(func1), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func4)]))

}

func TestAddAllDeploymentsLoop(t *testing.T) {
	deployments := []entities.Deployment{func3, func1}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	printDeployer(ctx, deploy)
	assert.Equal(t, 3, len(deploy.deployments))
	assert.Equal(t, FunctionName(func2), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func1), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 2, len(deploy.usedBy[FunctionName(func1)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func2)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func3)]))

}
func func1(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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
	return []entities.Deployment{func2, func3}, installer
}
func func2(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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
func func3(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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

func TestAddAllDeploymentsLoopSingle(t *testing.T) {

	deployments := []entities.Deployment{func6, func5}
	ctx := context.Background()
	deploy := NewDeployer()
	deploy.AddDeployments(ctx, deployments...)
	printDeployer(ctx, deploy)
	assert.Equal(t, 3, len(deploy.deployments))
	assert.Equal(t, FunctionName(func4), FunctionName(deploy.deployments[0]))
	assert.Equal(t, FunctionName(func5), FunctionName(deploy.deployments[len(deploy.deployments)-1]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func4)]))
	assert.Equal(t, 0, len(deploy.usedBy[FunctionName(func5)]))
	assert.Equal(t, 1, len(deploy.usedBy[FunctionName(func6)]))

}
func func4(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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
func func5(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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
	return []entities.Deployment{func6}, installer
}
func func6(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
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
	return []entities.Deployment{func4}, installer
}
