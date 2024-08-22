package config

type AppOption struct {
	DisableCLIParse bool
	IsStaked        bool
}

type Option func(*AppOption)

var DisableCLIParse = func(o *AppOption) {
	o.DisableCLIParse = true
}

var EnableStaked = func(o *AppOption) {
	o.IsStaked = true
}

func (a *AppOption) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(a)
	}
}
