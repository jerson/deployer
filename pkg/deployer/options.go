package deployer

// Option ...
type Option func(instance *Deployer)

// WithDeepLimit ...
func WithDeepLimit(limit int) Option {
	return func(instance *Deployer) {
		instance.deepLimit = limit
	}
}
