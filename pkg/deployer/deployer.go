package deployer

import (
	"context"
	"sync"

	"github.com/enescakir/emoji"
	"github.com/jerson/deployer/pkg/entities"
)

// ContextKey ...
type ContextKey string

const (
	// ContextKeyLog ...
	ContextKeyLog ContextKey = "log"
	// ContextKeyLogLevel ...
	ContextKeyLogLevel ContextKey = "log_level"
	// ContextKeyDisableSTDOUT ...
	ContextKeyDisableSTDOUT ContextKey = "disable_stdout"
	// ContextKeyTimeout ...
	ContextKeyTimeout ContextKey = "timeout"
)

// LogField ...
type LogField string

var (
	// LogFieldMethod ...
	LogFieldMethod = string(emoji.FlyingDisc)
	// LogFieldIndex ...
	LogFieldIndex = string(emoji.Star)
	// LogFieldRunner ...
	LogFieldRunner = string(emoji.Rocket)
	// LogFieldParent ...
	LogFieldParent = string(emoji.Satellite)
	// LogFieldDeployment ...
	LogFieldDeployment = string(emoji.LightBulb)
	// LogFieldDependency ...
	LogFieldDependency = string(emoji.Eyes)
)

// Deployer ...
type Deployer struct {
	deployments        []entities.Deployment
	deploymentsStarted []entities.Deployment
	deploymentStatus   entities.DeploymentStatus
	usedBy             entities.DependencyMap
	dependencyDeep     entities.DependencyDeep
	deepLimit          int
	sync.RWMutex
}

// NewDeployer ...
func NewDeployer(opts ...Option) *Deployer {

	instance := &Deployer{
		deployments:        []entities.Deployment{},
		deploymentsStarted: []entities.Deployment{},
		usedBy:             entities.DependencyMap{},
		dependencyDeep:     entities.DependencyDeep{},
		deploymentStatus:   entities.DeploymentStatus{},
		deepLimit:          10,
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

// Deployments ...
func (d *Deployer) Deployments() entities.Deployments {
	return d.deployments
}

// DeploymentStatusAll ...
func (d *Deployer) DeploymentStatusAll() entities.DeploymentStatus {
	return d.deploymentStatus
}

// DeploymentStatus ...
func (d *Deployer) DeploymentStatus(deployment entities.Deployment) *entities.Status {
	d.RLock()
	defer d.RUnlock()
	return d.deploymentStatus[FunctionName(deployment)]
}

// UsedByAll ...
func (d *Deployer) UsedByAll() entities.DependencyMap {
	return d.usedBy
}

// UsedBy ...
func (d *Deployer) UsedBy(deployment entities.Deployment) []entities.Deployment {
	d.RLock()
	defer d.RUnlock()
	return d.usedBy[FunctionName(deployment)]
}

// TODO implement this in the future for deep deps
//// DependsOn ...
//func (d *Deployer) DependsOn(deployment entities.Deployment) []entities.Deployment {
//	d.RLock()
//	defer d.RUnlock()
//	return d.usedBy[FunctionName(deployment)]
//}

// HealthAll ...
func (d *Deployer) HealthAll(ctx context.Context) error {
	for _, deployment := range d.deployments {
		_ = d.Health(ctx, deployment)
	}
	//var wg sync.WaitGroup
	//for _, deployment := range d.deploy {
	//	wg.Add(1)
	//	go func() {
	//		_ = d.Health(ctx,deployment)
	//		wg.Done()
	//	}()
	//}
	//wg.Wait()
	return nil
}

// HealthDeployed ...
func (d *Deployer) HealthDeployed(ctx context.Context) error {
	for _, deployment := range d.deploymentsStarted {
		err := d.Health(ctx, deployment)
		if err != nil {
			continue
		}
	}
	return nil
}

// Health ...
func (d *Deployer) Health(ctx context.Context, deployment entities.Deployment) (err error) {
	key := FunctionName(deployment)
	status := d.status(key)
	status.Health.Reset()

	log := Log(ctx)
	log = log.WithField(LogFieldDeployment, key)
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	_, installer := deployment(ctx)
	_, _, _, health := installer(ctx)
	d.deploymentsStarted = append(d.deploymentsStarted, deployment)

	if health == nil {
		log.Warn("health method missing")
		return nil
	}

	log = log.WithField(LogFieldRunner, "health")
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	status.Health.Start()

	//PrintTitle(ctx,key)
	log.Info("start")

	ctx, err = health(ctx)
	if err != nil {
		log.Error(err)
	}
	status.Health.End(err)
	log.Info("end")
	return err
}

// UpgradeAll ...
func (d *Deployer) UpgradeAll(ctx context.Context) error {

	d.print(ctx)
	for _, deployment := range d.deployments {
		err := d.Upgrade(ctx, deployment)
		if err != nil {
			continue
		}
	}
	return nil
}

// UpgradeDeployed ...
func (d *Deployer) UpgradeDeployed(ctx context.Context) error {

	d.print(ctx)
	for _, deployment := range d.deploymentsStarted {
		err := d.Upgrade(ctx, deployment)
		if err != nil {
			continue
		}
	}
	return nil
}

// Upgrade ...
func (d *Deployer) Upgrade(ctx context.Context, deployment entities.Deployment) (err error) {
	key := FunctionName(deployment)
	status := d.status(key)
	status.Upgrade.Reset()

	log := Log(ctx)
	log = log.WithField(LogFieldDeployment, key)
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	_, installer := deployment(ctx)
	_, upgrade, _, _ := installer(ctx)
	d.deploymentsStarted = append(d.deploymentsStarted, deployment)

	log = log.WithField(LogFieldRunner, "upgrade")
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	status.Upgrade.Start()

	PrintTitle(ctx, key)
	log.Info("start")

	ctx, err = upgrade(ctx)
	if err != nil {
		log.Error(err)
	}

	status.Upgrade.End(err)
	log.Info("end")
	return err
}

// InstallAll ...
func (d *Deployer) InstallAll(ctx context.Context) error {

	d.print(ctx)
	for _, deployment := range d.deployments {
		err := d.Install(ctx, deployment)
		if err != nil {
			return err
		}
	}
	return nil
}

// Install ...
func (d *Deployer) Install(ctx context.Context, deployment entities.Deployment) (err error) {
	key := FunctionName(deployment)
	status := d.status(key)
	status.Install.Reset()

	log := Log(ctx)
	log = log.WithField(LogFieldDeployment, key)
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	_, installer := deployment(ctx)
	install, _, _, _ := installer(ctx)
	d.deploymentsStarted = append(d.deploymentsStarted, deployment)

	log = log.WithField(LogFieldRunner, "install")
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	status.Install.Start()

	PrintTitle(ctx, key)
	log.Info("start")
	ctx, err = install(ctx)
	if err != nil {
		log.Error(err)
	}

	status.Install.End(err)
	log.Info("end")
	return err
}

// UninstallAll ...
func (d *Deployer) UninstallAll(ctx context.Context) error {

	for i := len(d.deployments) - 1; i > -1; i-- {
		deployment := d.deployments[i]
		err := d.Uninstall(ctx, deployment)
		if err != nil {
			continue
		}
	}
	return nil
}

// UninstallDeployed ...
func (d *Deployer) UninstallDeployed(ctx context.Context) error {

	for i := len(d.deploymentsStarted) - 1; i > -1; i-- {
		deployment := d.deploymentsStarted[i]
		err := d.Uninstall(ctx, deployment)
		if err != nil {
			continue
		}
	}
	return nil
}

// Uninstall ...
func (d *Deployer) Uninstall(ctx context.Context, deployment entities.Deployment) (err error) {
	key := FunctionName(deployment)
	status := d.status(key)
	status.Uninstall.Reset()

	log := Log(ctx)
	log = log.WithField(LogFieldDeployment, key)
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	_, installer := deployment(ctx)
	_, _, uninstall, _ := installer(ctx)

	log = log.WithField(LogFieldRunner, "uninstall")
	ctx = context.WithValue(ctx, ContextKeyLog, log)

	status.Uninstall.Start()

	PrintTitle(ctx, key)
	log.Info("start")

	ctx, err = uninstall(ctx)
	if err != nil {
		log.Warn(err)
	}
	status.Uninstall.End(err)
	log.Info("end")
	return err
}

// AddDeployments ...
func (d *Deployer) AddDeployments(ctx context.Context, deployments ...entities.Deployment) {
	for _, deployment := range deployments {
		d.addDeployments(ctx, deployment, nil)
	}
}

func (d *Deployer) addDeployments(ctx context.Context, deployment entities.Deployment, parent entities.Deployment) {

	log := Log(ctx)
	log = log.WithField(LogFieldDeployment, FunctionName(deployment))
	ctx = context.WithValue(ctx, ContextKeyLog, log)
	dependencies, _ := deployment(ctx)
	for _, dependency := range dependencies {
		if !d.allowDeep(dependency) {
			log.WithField(LogFieldDependency, FunctionName(dependency)).
				Warn("max deep limit reached:", d.deepLimit)
			continue
		}
		d.addDeployments(ctx, dependency, deployment)
	}

	d.addDeployment(ctx, deployment, parent)
}

func (d *Deployer) allowDeep(dependency entities.Deployment) bool {
	keyDependency := FunctionName(dependency)
	d.dependencyDeep[keyDependency]++
	return d.dependencyDeep[keyDependency] < d.deepLimit
}

func (d *Deployer) addDeployment(ctx context.Context, deployment entities.Deployment, parent entities.Deployment) {

	log := Log(ctx)
	log = log.
		WithField(LogFieldDeployment, FunctionName(deployment)).
		WithField(LogFieldParent, FunctionName(parent)).
		WithField(LogFieldMethod, "addDeployment")

	key := FunctionName(deployment)
	if key == FunctionName(parent) {
		log.Warn("duplicated")
		return
	}

	if d.usedBy[key] == nil {
		d.usedBy[key] = []entities.Deployment{}
	}

	for _, dependency := range d.usedBy[key] {
		if FunctionName(dependency) == FunctionName(parent) {
			log.Trace("skipped duplicated dependency")
			return
		}
	}

	if parent != nil {
		d.usedBy[key] = append(d.usedBy[key], parent)
	}

	for _, deploymentItem := range d.deployments {
		if key == FunctionName(deploymentItem) {
			log.Debug("skipped duplicated")
			return
		}
	}

	d.deployments = append(d.deployments, deployment)
}

func (d *Deployer) status(key string) *entities.Status {

	d.Lock()
	defer d.Unlock()
	if d.deploymentStatus[key] == nil {
		d.deploymentStatus[key] = entities.NewStatus()
	}

	return d.deploymentStatus[key]
}

func (d *Deployer) print(ctx context.Context) {
	log := Log(ctx)
	log = log.WithField(LogFieldMethod, "InstallAll")
	PrintTitle(ctx, "Console:")
	for key, deployment := range d.deployments {
		log.WithField(LogFieldIndex, key).Info(FunctionName(deployment))
	}
}
