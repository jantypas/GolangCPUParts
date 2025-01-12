package PhysicalMemory

var MemoryMapTable = map[string][]PhysicalMemoryRegion{
	"OLD-IBM-MAINFRAME": {
		{
			Comment:    "2MB of physical RAM",
			NumPages:   0x200,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"KAYPRO": {
		{
			Comment:    "64KB Physical RAM",
			NumPages:   0x10,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"VAX": {
		{
			Comment:    "64MB Virtual RAM",
			NumPages:   0x4000,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "16MB Physical RAM",
			NumPages:   0x1000,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"LINUX-X64": {
		{
			Comment:    "4GB Virtual RAM",
			NumPages:   0x10000,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "512MB IO RAM",
			NumPages:   0x20000,
			MemoryType: MemoryTypeIO,
		},
		{
			Comment:    "1GB ROM Space",
			NumPages:   0x40000,
			MemoryType: MemoryTypeROM,
		},
		{
			Comment:    "1GB Kernel Space",
			NumPages:   0x40000,
			MemoryType: MemoryTypeKernel,
		},
	},
}
