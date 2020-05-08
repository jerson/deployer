package deploy

import (
	"context"
	"github.com/jerson/deployer/pkg/entities"
)

// Sample2 ...
func Sample2(_ context.Context) (dependencies []entities.Deployment, installer entities.Installer) {
	installer = func(ctx context.Context) (install entities.Install, upgrade entities.Upgrade, uninstall entities.Uninstall, health entities.Health) {

		health = func(ctx context.Context) (context.Context, error) {
			return delay(ctx,1)
		}
		uninstall = func(ctx context.Context) (context.Context, error) {
			return delay(ctx,5)
		}
		upgrade = func(ctx context.Context) (context.Context, error) {
			return delay(ctx,3)
		}
		install = func(ctx context.Context) (context.Context, error) {
			return delay(ctx,10)
		}
		return install, upgrade, uninstall, health
	}
	return []entities.Deployment{}, installer
}
