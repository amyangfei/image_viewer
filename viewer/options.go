package viewer

type Options struct {
	Headless bool `flag:"headless"`
}

func NewOptions() *Options {
	return &Options{
		Headless: false,
	}
}
