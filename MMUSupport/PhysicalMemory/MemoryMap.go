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
			NumPages:   0x4_000,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "16MB Physical RAM",
			NumPages:   0x1_000,
			MemoryType: MemoryTypePhysicalRAM,
		},
	},
	"LINUX-X64": {
		{
			Comment:    "4GB Virtual RAM",
			NumPages:   0x100_000,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "512MB IO RAM",
			NumPages:   0x20_000,
			MemoryType: MemoryTypeIORAM,
		},
		{
			Comment:    "Buffer RAM 8MB",
			NumPages:   0x800,
			MemoryType: MemoryTypeBufferRAM,
		},
		{
			Comment:    "1GB ROM Space",
			NumPages:   0x40_000,
			MemoryType: MemoryTypePhysicalROM,
		},
		{
			Comment:    "1GB Kernel Space",
			NumPages:   0x40_000,
			MemoryType: MemoryTypeKernelRAM,
		},
	},
	"TEST": {
		{
			Comment:    "1MB Virtual RAM",
			NumPages:   4,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "1MB Physical RAM",
			NumPages:   4,
			MemoryType: MemoryTypePhysicalRAM,
		},
		{
			Comment:    "1MB I/O RAM",
			NumPages:   4,
			MemoryType: MemoryTypeIORAM,
		},
		{
			Comment:    "1MB Buffer RAM",
			NumPages:   4,
			MemoryType: MemoryTypeBufferRAM,
		},
		{
			Comment:    "1MB Physical ROM",
			NumPages:   4,
			MemoryType: MemoryTypePhysicalROM,
		},
		{
			Comment:    "1MB Empty space",
			NumPages:   4,
			MemoryType: MemoryTypeEmpty,
		},
		{
			Comment:    "1MB Kernel ROM",
			NumPages:   4,
			MemoryType: MemoryTypeKernelROM,
		},
		{
			Comment:    "1MB IO ROM",
			NumPages:   4,
			MemoryType: MemoryTypeIOROM,
		},
	},
}
