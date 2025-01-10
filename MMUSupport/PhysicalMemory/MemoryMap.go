package PhysicalMemory

var MemoryMapTable = map[string][]PhysicalMemoryRegion{
	"OLD-IBM-MAINFRAME": {
		{
			Comment:    "2MB of physical RAM",
			StartPage:  0x0000_0000,
			EndAddress: 0x001F_FFFF,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"KAYPRO": {
		{
			Comment:    "64KB Physical RAM",
			StartPage:  0x0000,
			EndAddress: 0xFFFF,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"VAX": {
		{
			Comment:    "64MB Virtual RAM",
			StartPage:  0x0000_0000,
			EndAddress: 0x03FF_FFFF,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "16MB Physical RAM",
			StartPage:  0x0400_0000,
			EndAddress: 0x04FF_FFFF,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"LINUX-X64": {
		{
			Comment:    "4GB Virtual RAM",
			StartPage:  0x0000_0000,
			EndAddress: 0xFFFF_FFFF,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "512MB IO RAM",
			StartPage:  0x1_0000_0000,
			EndAddress: 0x1_1FFF_0000,
			MemoryType: MemoryTypeIO,
		},
		{
			Comment:    "1GB ROM Space",
			StartPage:  0x1_2000_0000,
			EndAddress: 0x1_5FFF_FFFF,
			MemoryType: MemoryTypeROM,
		},
		{
			Comment:    "1GB Kernel Space",
			StartPage:  0x1_8000_0000,
			EndAddress: 0x1_9FFF_FFFF,
			MemoryType: MemoryTypeKernel,
		},
	},
}
