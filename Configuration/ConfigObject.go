package Configuration

type ConfigObject struct {
	SwapFileNames string
}

func (cfg *ConfigObject) Initialize(cfgsrc string) {
	cfg.SwapFileNames = "/tmp/swap.swp"
}
