package PortIO

type PortIOConfigObject struct {
	Name            string
	HandleInByte    func(port uint64) (byte, error)
	HandleInWOrd    func(port uint64) (uint16, error)
	HandleInDouble  func(port uint64) (uint32, error)
	HandleInQuad    func(port uint64) (uint64, error)
	HandleOutByte   func(port uint64, value byte) error
	HandleOutWOrd   func(port uint64, value uint16) error
	HandleOutDouble func(port uint64, value uint32) error
	HandleOutQuad   func(port uint64, value uint64) error
}

const (
	LegacyBeeperPort         = 0x1000_0000_0000_0000
	LegacyLedPanelPort       = 0x1000_0000_0000_0001
	LegacyClockControlPort   = 0x1000_0000_0001_0000
	LegacyClockDataPort      = 0x1000_0000_0001_0001
	LegacyClockInterruptPort = 0x1000_0000_0001_0002
)

var PortIOConfig = map[uint64]PortIOConfigObject{
	LegacyBeeperPort: PortIOConfigObject{
		Name:          "LegacyBeeper",
		HandleOutByte: func(port uint64, value byte) error { return nil },
	},
	LegacyLedPanelPort: PortIOConfigObject{
		Name:          "LegacyLedPanel",
		HandleOutQuad: func(port uint64, value uint64) error { return nil },
	},
	LegacyClockControlPort: PortIOConfigObject{
		Name:          "LegacyClockControl",
		HandleOutByte: func(port uint64, value byte) error { return nil },
	},
	LegacyClockDataPort: PortIOConfigObject{
		Name:          "LegacyClockData",
		HandleInQuad:  func(port uint64) (uint64, error) { return 0, nil },
		HandleOutQuad: func(port uint64, value uint64) error { return nil },
	},
	LegacyClockInterruptPort: PortIOConfigObject{
		Name:          "LegacyClockInterrupt",
		HandleOutQuad: func(port uint64, value uint64) error { return nil },
	},
}
