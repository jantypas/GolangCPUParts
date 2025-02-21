package Configuration

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CPUDescriptor struct {
	CPUType    uint64            `json:"cpu_type"`
	FeatureA   uint64            `json:"feature_a"`
	FeatureB   uint64            `json:"feature_b"`
	Parameters map[string]string `json:"parameters"`
}

type MemoryDescriptor struct {
	Key          int               `json:"key"`
	Comment      string            `json:"comment"`
	StartAddress uint64            `json:"start_address"`
	EndAddress   uint64            `json:"end_address"`
	MemoryType   string            `json:"memory_type"`
	Parameters   map[string]string `json:"parameters"`
}

type IODescriptor struct {
	Class      string            `json:"class"`
	Subclass   string            `json:"subclass"`
	Model      string            `json:"model"`
	MountPoint string            `json:"mountPoint"`
	Parameters map[string]string `json:"parameters"`
}

type ConfigurationDescriptor struct {
	CPU    CPUDescriptor      `json:"cpu"`
	Memory []MemoryDescriptor `json:"memory"`
	IO     []IODescriptor     `json:"IO"`
}

type ConfigSettings struct {
	SwapFileName   string
	HostVolumePath string
}

type SystemConfigs struct {
	Name        string                  `json:"name"`
	Description ConfigurationDescriptor `json:"description"`
}

type ConfigObject struct {
	Version       int             `json:"version"`
	Settings      ConfigSettings  `json:"settings"`
	Configuration []SystemConfigs `json:"configuration"`
}

var MemoryTypeNames = []string{
	"Empty",
	"Virtual-RAM",
	"Physical-RAM",
	"Buffer-RAM",
	"Kernel-RAM",
	"I/O-RAM",
	"Physical-ROM",
}

var MemoryTypeValues = map[string]uint64{}

func MockConfig() ([]byte, error) {
	cfg := ConfigObject{
		Version: 1,
		Settings: ConfigSettings{
			SwapFileName:   "/tmp/swap.swp",
			HostVolumePath: "/tmp/host/volumes",
		},
		Configuration: []SystemConfigs{
			SystemConfigs{
				Name: "Old-IBM-Mainframe",
				Description: ConfigurationDescriptor{
					CPU: CPUDescriptor{
						CPUType:    1000_0000_0000_0000,
						FeatureA:   0000_0000_0000_0000,
						FeatureB:   0000_0000_0000_0000,
						Parameters: map[string]string{},
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "2MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_001F_FFFF,
							MemoryType:   "Physical-RAM",
							Parameters:   map[string]string{},
						},
						{
							Key:          1,
							Comment:      "Kernel RAM",
							StartAddress: 0x0000_0000_0002_0000,
							EndAddress:   0x0000_0000_0002_FFFF,
							MemoryType:   "Kernel-RAM",
							Parameters: map[string]string{
								"preload": "mon:winiloader.bin",
							},
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "Legacy",
							Subclass:   "Panel",
							Model:      "ButtonBox",
							MountPoint: "/dev/button",
							Parameters: map[string]string{},
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card/0",
							Parameters: map[string]string{
								"mode": "read",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-reader/1",
							Parameters: map[string]string{
								"mode": "read",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-punch/0",
							Parameters: map[string]string{
								"mode": "punch",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-punch/1",
							Parameters: map[string]string{
								"mode": "punch",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Printer",
							Model:      "Printer",
							MountPoint: "/dev/printer/0",
							Parameters: map[string]string{
								"mode": "text",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/0",
							Parameters: map[string]string{
								"mode": "800bpi",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/1",
							Parameters: map[string]string{
								"mode": "800bpi",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/2",
							Parameters: map[string]string{
								"mode": "800bpi",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "Wini",
							MountPoint: "/dev/disk/0",
							Parameters: map[string]string{
								"size":  "30MB",
								"model": "Wini3030",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "Wini",
							MountPoint: "/dev/disk/1",
							Parameters: map[string]string{
								"size":  "30MB",
								"model": "Wini3030",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "Wini",
							MountPoint: "/dev/disk/2",
							Parameters: map[string]string{
								"size":  "30MB",
								"model": "Wini3030",
							},
						},
						{
							Class:      "Legacy",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
							Parameters: map[string]string{},
						},
						{
							Class:      "Legacy",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/1",
							Parameters: map[string]string{},
						},
					},
				},
			},
			{
				Name: "Kaypro-CPM-64KB",
				Description: ConfigurationDescriptor{
					CPU: CPUDescriptor{
						CPUType:    1000_0000_0000_0000,
						FeatureA:   0000_0000_0000_0000,
						FeatureB:   0000_0000_0000_0000,
						Parameters: map[string]string{},
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "64MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_0000_FFFF,
							MemoryType:   "Physical-RAM",
							Parameters:   map[string]string{},
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "80s-CPM",
							Subclass:   "Keyboard",
							Model:      "ASCII",
							MountPoint: "/dev/keyboard",
							Parameters: map[string]string{},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Printer",
							Model:      "ASCII",
							MountPoint: "/dev/printer/0",
							Parameters: map[string]string{
								"mode": "text",
							},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Display",
							Model:      "ASCII",
							MountPoint: "/dev/display",
							Parameters: map[string]string{
								"mode": "text",
							},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Disk",
							Model:      "Shugart",
							MountPoint: "/dev/disk/0",
							Parameters: map[string]string{
								"size": "20MB",
							},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Floppy",
							Model:      "Flimsiwrite",
							MountPoint: "/dev/floppy/0",
							Parameters: map[string]string{
								"tracks": "80",
							},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Floppy",
							Model:      "Flimsiwrite",
							MountPoint: "/dev/floppy/1",
							Parameters: map[string]string{
								"tracks": "80",
							},
						},
						{
							Class:      "80s-CPM",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
							Parameters: map[string]string{
								"mode":  "hayes",
								"speed": "57600",
							},
						},
					},
				},
			},
			{
				Name: "Vax-11/780-64MB",
				Description: ConfigurationDescriptor{
					CPU: CPUDescriptor{
						CPUType:    1000_0000_0000_0000,
						FeatureA:   0000_0000_0000_0000,
						FeatureB:   0000_0000_0000_0000,
						Parameters: map[string]string{},
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "64MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_03FF_FFFF,
							MemoryType:   "Virtual-RAM",
							Parameters:   map[string]string{},
						},
						{
							Key:          1,
							Comment:      "Kernel 16MB RAM",
							StartAddress: 0x0000_0000_0400_0000,
							EndAddress:   0x0000_0000_04FF_FFFF,
							MemoryType:   "Kernel-RAM",
							Parameters:   map[string]string{},
						},
						{
							Key:          2,
							Comment:      "I/O RAM 16MB",
							StartAddress: 0x0000_0000_0500_0000,
							EndAddress:   0x0000_0000_05FF_FFFF,
							MemoryType:   "I/O-RAM",
							Parameters:   map[string]string{},
						},
						{
							Key:          3,
							Comment:      "512KB Buffer RAM",
							StartAddress: 0x0000_0000_0600_0000,
							EndAddress:   0x0000_0000_0607_FFFF,
							MemoryType:   "Buffer-RAM",
							Parameters:   map[string]string{},
						},
						{
							Key:          4,
							Comment:      "512KB System RAM",
							StartAddress: 0x0000_0000_0608_0000,
							EndAddress:   0x0000_0000_060F_FFFF,
							MemoryType:   "Physical-RAM",
							Parameters:   map[string]string{},
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "VAX",
							Subclass:   "Console",
							Model:      "ASCII",
							MountPoint: "/dev/console",
							Parameters: map[string]string{},
						},
						{
							Class:      "VAX",
							Subclass:   "Printer",
							Model:      "ASCII",
							MountPoint: "/dev/printer/0",
							Parameters: map[string]string{
								"mode": "text",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Network",
							Model:      "Ethernet",
							MountPoint: "/dev/ethernet/0",
							Parameters: map[string]string{
								"mode": "10base-T",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/0",
							Parameters: map[string]string{
								"mode": "1600bpi",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/1",
							Parameters: map[string]string{
								"mode": "1600bpi",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "RU0K",
							MountPoint: "/dev/disk/0",
							Parameters: map[string]string{
								"size": "120MB",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "RU0K",
							MountPoint: "/dev/disk/1",
							Parameters: map[string]string{
								"size": "250MB",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "RU0K",
							MountPoint: "/dev/disk/2",
							Parameters: map[string]string{
								"size": "500MB",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "RU0K",
							MountPoint: "/dev/disk/3",
							Parameters: map[string]string{
								"size": "500MB",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
							Parameters: map[string]string{
								"mode":  "hayes",
								"speed": "57600",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "PTY",
							Model:      "NA",
							MountPoint: "/dev/pty/0",
							Parameters: map[string]string{
								"num": "64",
							},
						},
						{
							Class:      "VAX",
							Subclass:   "TTY",
							Model:      "NA",
							MountPoint: "/dev/tty/0",
							Parameters: map[string]string{
								"num": "64",
							},
						},
					},
				},
			},
		},
	}

	s, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func checkMemoryType(s string) bool {
	for _, v := range MemoryTypeNames {
		if v == s {
			return true
		}
	}
	return false
}

func LoadConfiguration(s []byte) (*ConfigObject, error) {
	cfg := ConfigObject{}
	err := json.Unmarshal([]byte(s), &cfg)
	if err != nil {
		return nil, err
	}
	for _, v := range cfg.Configuration {
		// Verify the memory type names are correct
		s := v.Description.Memory[0].MemoryType
		if !checkMemoryType(s) {
			return nil, errors.New("Invalid memory type " + s + " in configuration")
		}
	}
	return &cfg, nil
}

func (cfg *ConfigObject) Save() string {
	s, _ := json.Marshal(cfg)
	return string(s)
}

func (cfg *ConfigObject) Dump() {
	s, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(s))
}

func (cfg *ConfigObject) GetConfigByName(s string) *SystemConfigs {
	for _, v := range cfg.Configuration {
		if v.Name == s {
			return &v
		}
	}
	return nil
}

func (cfg *ConfigObject) GetConfigurationSettings() *ConfigSettings {
	return &cfg.Settings
}
