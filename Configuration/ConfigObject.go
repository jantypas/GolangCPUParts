package Configuration

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CPUDescriptor struct {
	CPUType    uint64 `json:"cpu_type"`
	FeatureA   uint64 `json:"feature_a"`
	FeatureB   uint64 `json:"feature_b"`
	Parameters string `json:"parameters"`
}

type MemoryDescriptor struct {
	Key          int    `json:"key"`
	Comment      string `json:"comment"`
	StartAddress uint64 `json:"start_address"`
	EndAddress   uint64 `json:"end_address"`
	MemoryType   string `json:"memory_type"`
	Parameters   string `json:"parameters"`
}

type IODescriptor struct {
	Class      string `json:"class"`
	Subclass   string `json:"subclass"`
	Model      string `json:"model"`
	MountPoint string `json:"mountPoint"`
	Parameters string `json:"parameters"`
}

type ConfigurationDescriptor struct {
	CPU    CPUDescriptor      `json:"cpu"`
	Memory []MemoryDescriptor `json:"memory"`
	IO     []IODescriptor     `json:"IO"`
}

type ConfigSettings struct {
	SwapFileName string `json:"swap_file_name"`
}

type SystemConfigs struct {
	Name        string                  `json:"name"`
	Description ConfigurationDescriptor `json:"description"`
}

type ConfigObject struct {
	Version       int             `json:"version"`
	Settings      ConfigSettings  `json:"settings"`
	Conifguration []SystemConfigs `json:"configuration"`
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
			SwapFileName: "/tmp/swap.swp",
		},
		Conifguration: []SystemConfigs{
			SystemConfigs{
				Name: "Old-IBM-Mainframe",
				Description: ConfigurationDescriptor{
					CPU: CPUDescriptor{
						CPUType:    1000_0000_0000_0000,
						FeatureA:   0000_0000_0000_0000,
						FeatureB:   0000_0000_0000_0000,
						Parameters: "",
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "2MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_001F_FFFF,
							MemoryType:   "Physical-RAM",
							Parameters:   "",
						},
						{
							Key:          1,
							Comment:      "Kernel RAM",
							StartAddress: 0x0000_0000_0002_0000,
							EndAddress:   0x0000_0000_0002_FFFF,
							MemoryType:   "Kernel-RAM",
							Parameters:   "",
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "Legacy",
							Subclass:   "Panel",
							Model:      "ButtonBox",
							MountPoint: "/dev/button",
							Parameters: "",
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card/0",
							Parameters: "mode=read",
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-reader/1",
							Parameters: "mode=read",
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-punch/0",
							Parameters: "mode=punch",
						},
						{
							Class:      "Legacy",
							Subclass:   "CardService",
							Model:      "CardALot",
							MountPoint: "/dev/card-punch/1",
							Parameters: "mode=punch",
						},
						{
							Class:      "Legacy",
							Subclass:   "Printer",
							Model:      "Printer",
							MountPoint: "/dev/printer/0",
							Parameters: "mode=text",
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/0",
							Parameters: "mode=800bpi",
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/1",
							Parameters: "mode=800bpi",
						},
						{
							Class:      "Legacy",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/2",
							Parameters: "mode=800bpi",
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "120MB",
							MountPoint: "/dev/disk/0",
							Parameters: "size=120MB",
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "250MB",
							MountPoint: "/dev/disk/1",
							Parameters: "size=250MB",
						},
						{
							Class:      "Legacy",
							Subclass:   "Disk",
							Model:      "500MB",
							MountPoint: "/dev/disk/2",
							Parameters: "size=500MB",
						},
						{
							Class:      "Legacy",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
							Parameters: "",
						},
						{
							Class:      "Legacy",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/1",
							Parameters: "",
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
						Parameters: "",
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "64MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_0000_FFFF,
							MemoryType:   "Physical-RAM",
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "80s-CPM",
							Subclass:   "Keyboard",
							Model:      "ASCII",
							MountPoint: "/dev/keyboard",
							Parameters: "",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Printer",
							Model:      "ASCII",
							MountPoint: "/dev/printer/0",
							Parameters: "mode=text",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Display",
							Model:      "ASCII",
							MountPoint: "/dev/display",
							Parameters: "mode=text",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Disk",
							Model:      "120MB",
							MountPoint: "/dev/disk/0",
							Parameters: "size=120MB",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Floppy",
							Model:      "80Track",
							MountPoint: "/dev/floppy/0",
							Parameters: "size=80",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "Floppy",
							Model:      "120Track",
							MountPoint: "/dev/floppy/1",
							Parameters: "size=120",
						},
						{
							Class:      "80s-CPM",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
							Parameters: "size=144",
						},
					},
				},
			},
			{
				Name: "Vax-11/780r-64MB",
				Description: ConfigurationDescriptor{
					CPU: CPUDescriptor{
						CPUType:    0x1000_0000_0000_0000,
						FeatureA:   0x0000_0000_0000_0000,
						FeatureB:   0x0000_0000_0000_0000,
						Parameters: "",
					},
					Memory: []MemoryDescriptor{
						MemoryDescriptor{
							Key:          0,
							Comment:      "64MB of RAM",
							StartAddress: 0x0000_0000_0000_0000,
							EndAddress:   0x0000_0000_03FFF_FFFF,
							MemoryType:   "Virtual-RAM",
							Parameters:   "",
						},
						{
							Key:          1,
							Comment:      "Kernel 16MB RAM",
							StartAddress: 0x0000_0000_0400_0000,
							EndAddress:   0x0000_0000_04FF_FFFF,
							MemoryType:   "Kernel-RAM",
							Parameters:   "",
						},
						{
							Key:          2,
							Comment:      "I/O RAM 16MB",
							StartAddress: 0x0000_0000_0500_0000,
							EndAddress:   0x0000_0000_05FF_FFFF,
							MemoryType:   "I/O-RAM",
							Parameters:   "",
						},
						{
							Key:          3,
							Comment:      "512KB Buffer RAM",
							StartAddress: 0x0000_0000_0600_0000,
							EndAddress:   0x0000_0000_0607_FFFF,
							MemoryType:   "Buffer-RAM",
						},
					},
					IO: []IODescriptor{
						IODescriptor{
							Class:      "VAX",
							Subclass:   "Console",
							Model:      "ASCII",
							MountPoint: "/dev/console",
							Parameters: "",
						},
						{
							Class:      "VAX",
							Subclass:   "Printer",
							Model:      "ASCII",
							MountPoint: "/dev/printer/0",
						},
						{
							Class:      "VAX",
							Subclass:   "Network",
							Model:      "Ethernet",
							MountPoint: "/dev/ethernet/0",
							Parameters: "mode=10Mb",
						},
						{
							Class:      "VAX",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/0",
							Parameters: "mode=1600bpi",
						},
						{
							Class:      "VAX",
							Subclass:   "Tape",
							Model:      "TK",
							MountPoint: "/dev/tape/1",
							Parameters: "mode=1600bpi",
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "120MB",
							MountPoint: "/dev/disk/0",
							Parameters: "size=120MB",
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "250MB",
							MountPoint: "/dev/disk/1",
							Parameters: "size=250Mb",
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "500MB",
							MountPoint: "/dev/disk/2",
							Parameters: "size=500MB",
						},
						{
							Class:      "VAX",
							Subclass:   "Disk",
							Model:      "500MB",
							MountPoint: "/dev/disk/3",
						},
						{
							Class:      "VAX",
							Subclass:   "COM",
							Model:      "Modem",
							MountPoint: "/dev/modem/0",
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
	for _, v := range cfg.Conifguration {
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
	for _, v := range cfg.Conifguration {
		if v.Name == s {
			return &v
		}
	}
	return nil
}

func (cfg *ConfigObject) GetConfigurationSettings() *ConfigSettings {
	return &cfg.Settings
}
