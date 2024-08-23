package config

type AppOption struct {
	DisableCLIParse bool
	IsStaked        bool
	Bootnodes       []string
	RandomIdentity  bool
}

type Option func(*AppOption)

var DisableCLIParse = func(o *AppOption) {
	o.DisableCLIParse = true
}

var EnableStaked = func(o *AppOption) {
	o.IsStaked = true
}

var EnableRandomIdentity = func(o *AppOption) {
	o.RandomIdentity = true
}

func (a *AppOption) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(a)
	}
}

func WithBootNodes(bootnodes ...string) Option {
	return func(o *AppOption) {
		o.Bootnodes = append(o.Bootnodes, bootnodes...)
	}
}
