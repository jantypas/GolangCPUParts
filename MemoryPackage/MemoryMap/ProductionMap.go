package MemoryMap

type MemoryRegionList []MemoryMapRegion

var ProductionMap = map[string]MemoryRegionList{
	"OLD-IBM-MAINFRAME": []MemoryMapRegion{
		{
			Key:          0,
			Comment:      "2MB Physical RAM",
			Tag:          "PHYSICAL-RAM",
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_001F_FFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable,
			SegmentType:  SegmentTypePhysicalRAM,
		},
		{
			Key:          1,
			Comment:      "1MB Physical RAM",
			Tag:          "KERNEL-RAM",
			StartAddress: 0x0000_0000_0020_0000,
			EndAddress:   0x0000_0000_002F_FFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable | ProtectionSystem,
			SegmentType:  SegmentTypePhysicalRAM,
		},
	},
	"KAYPRO-CPM": {
		{
			Key:          0,
			Comment:      "64KB Physical RAM",
			Tag:          "PHYSICAL-RAM",
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable,
			SegmentType:  SegmentTypePhysicalRAM,
		},
	},
	"VAX-11/780": {
		{
			Key:          0,
			Comment:      "64MB Virtual RAM",
			Tag:          "VIRTUAL-RAM",
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_03FF_FFFF,
			Permissions:  ProtectionExecutable | ProtectionWritable,
			SegmentType:  SegmentTypeVirtualRAM,
		},
		{
			Key:          1,
			Comment:      "16MB Physical RAM",
			Tag:          "KERNEL-RAM",
			StartAddress: 0x0000_0000_0450_0000,
			EndAddress:   0x0000_0000_04FF_FFFF,
			Permissions:  ProtectionExecutable | ProtectionWritable | ProtectionSystem,
			SegmentType:  SegmentTypePhysicalRAM,
		},
	},
	"Linux-8GB": {
		{
			Key:          0,
			Comment:      "8GB Virtual RAM",
			Tag:          "VIRTUAL-RAM",
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0001_FFFF_FFFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable,
			SegmentType:  SegmentTypeVirtualRAM,
		},
		{
			Key:          1,
			Comment:      "1GB ROM",
			Tag:          "ROM",
			StartAddress: 0x0000_0002_8000_0000,
			EndAddress:   0x0000_0002_4000_0000,
			Permissions:  ProtectionExecutable,
			SegmentType:  SegmentTypePhysicalRAM,
		},
		{
			Key:          2,
			Comment:      "1GB IO RAM",
			StartAddress: 0x0000_0002_C000_0000,
			EndAddress:   0x0000_0002_FFFF_FFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable,
			SegmentType:  SegmentTypePhysicalIO,
		},
		{
			Key:          3,
			Comment:      "2GB Kernel RAM",
			Tag:          "KERNEL-RAM",
			StartAddress: 0x0000_0003_0000_0000,
			EndAddress:   0x0000_0003_7FFF_FFFF,
			Permissions:  ProtectionWritable | ProtectionExecutable | ProtectionSystem | SegmentLocked,
		},
	},
}
