package workers

type WorkerOption struct {
	isTwitterWorker        bool
	isWebScraperWorker     bool
	isDiscordScraperWorker bool
	masaDir                string
}

type WorkerOptionFunc func(*WorkerOption)

var EnableTwitterWorker = func(o *WorkerOption) {
	o.isTwitterWorker = true
}

var EnableWebScraperWorker = func(o *WorkerOption) {
	o.isWebScraperWorker = true
}

var EnableDiscordScraperWorker = func(o *WorkerOption) {
	o.isDiscordScraperWorker = true
}

func WithMasaDir(dir string) WorkerOptionFunc {
	return func(o *WorkerOption) {
		o.masaDir = dir
	}
}

func (a *WorkerOption) Apply(opts ...WorkerOptionFunc) {
	for _, opt := range opts {
		opt(a)
	}
}
