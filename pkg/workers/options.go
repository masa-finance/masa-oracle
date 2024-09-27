package workers

type WorkerOption struct {
	isTwitterWorker        bool
	isWebScraperWorker     bool
	isLLMServerWorker      bool
	isDiscordScraperWorker bool
}

type WorkerOptionFunc func(*WorkerOption)

var EnableTwitterWorker = func(o *WorkerOption) {
	o.isTwitterWorker = true
}

var EnableWebScraperWorker = func(o *WorkerOption) {
	o.isWebScraperWorker = true
}

var EnableLLMServerWorker = func(o *WorkerOption) {
	o.isLLMServerWorker = true
}

var EnableDiscordScraperWorker = func(o *WorkerOption) {
	o.isDiscordScraperWorker = true
}

func (a *WorkerOption) Apply(opts ...WorkerOptionFunc) {
	for _, opt := range opts {
		opt(a)
	}
}
