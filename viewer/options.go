package viewer

type Options struct {
	Headless   bool `flag:"headless"`
	DriverPort int  `flag:"driver-port"`
}

func NewOptions() *Options {
	return &Options{
		Headless:   false,
		DriverPort: 9515,
	}
}
