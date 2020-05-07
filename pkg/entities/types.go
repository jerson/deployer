package entities

import (
	"context"
)

// Install ...
type Install func(ctx context.Context) (context.Context, error)

// Upgrade ...
type Upgrade func(ctx context.Context) (context.Context, error)

// Health ...
type Health func(ctx context.Context) (context.Context, error)

// Uninstall ...
type Uninstall func(ctx context.Context) (context.Context, error)

// Installer ...
type Installer func(ctx context.Context) (install Install, upgrade Upgrade, uninstall Uninstall, health Health)

// Deployment ...
type Deployment func(ctx context.Context) (dependencies []Deployment, installer Installer)

// Deployments ...
type Deployments []Deployment

// DeploymentStatus ...
type DeploymentStatus map[string]*Status

// DependencyMap ...
type DependencyMap map[string][]Deployment

// DependencyDeep ...
type DependencyDeep map[string]int
