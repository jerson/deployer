package deployer

import (
	"context"
	"strings"
	"testing"

	"github.com/jerson/deployer/pkg/entities"

	"github.com/stretchr/testify/assert"
)

func TestNewDeployer(t *testing.T) {

	var func1 = func(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
		installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {
			health = func(ctx context.Context) (context.Context, error) {
				return ctx, nil
			}
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
	deploy := NewDeployer(WithDeepLimit(10))
	deploy.AddDeployments(ctx, deployments...)
	err := deploy.InstallAll(ctx)
	err = deploy.UpgradeDeployed(ctx)
	err = deploy.UpgradeAll(ctx)
	err = deploy.HealthAll(ctx)
	err = deploy.HealthDeployed(ctx)
	err = deploy.UninstallAll(ctx)
	err = deploy.UninstallDeployed(ctx)

	printDeployer(ctx, deploy)

	assert.NoError(t, err)
}

func deploymentNames(deployments []entities.Deployment) []string {
	var names []string
	for _, deployment := range deployments {
		names = append(names, FunctionName(deployment))
	}
	return names
}

func printDeployer(ctx context.Context, deploy *Deployer) {

	log := Log(ctx)

	if len(deploy.deployments) > 0 {
		log.Info("")
		log.Info("deployments sorted:")
		for key, deployment := range deploymentNames(deploy.deployments) {
			log.WithField("index", key).Info(deployment)
		}
	}

	if len(deploy.usedBy) > 0 {
		log.Info("")
		log.Info("dependencies map:")
		for key, dependencies := range deploy.usedBy {
			names := deploymentNames(dependencies)
			log.WithField("required by", strings.Join(names, ", ")).
				Info(key)
		}
	}

	if len(deploy.deploymentsStarted) > 0 {
		log.Info("")
		log.Info("deploymentsStarted:")
		for key, uninstaller := range deploy.deploymentsStarted {
			log.WithField("index", key).Info(FunctionName(uninstaller))
		}
	}

	if len(deploy.deploymentStatus) > 0 {
		log.Info("")
		log.Info("deployment status:")
		for key, status := range deploy.deploymentStatus {
			log.WithField("install_error", status.Install.Error).
				WithField("install_duration", status.Install.Took()).
				WithField("uninstall_error", status.Uninstall.Error).
				WithField("uninstall_duration", status.Uninstall.Took()).
				Info(key)
		}
	}
}
