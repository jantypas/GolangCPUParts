package PhysicalMemory

const (
	MemoryTypeEmpty       = 0x0001
	MemoryTypeVirtualRAM  = 0x0002
	MemoryTypePhysicalRAM = 0x0004
	MemoryTypePhysicalROM = 0x0008
	MemoryTypeKernelRAM   = 0x0010
	MemoryTypeKernelROM   = 0x0020
	MemoryTypeIORAM       = 0x0040
	MemoryTypeIOROM       = 0x0080
	MemoryTypeBufferRAM   = 0x0100

	PageSize = 4096
)

type PhysicalMemoryRegion struct {
	Comment    string
	NumPages   uint32
	MemoryType int
}

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
			NumPages:   0x100,
			MemoryType: MemoryTypeVirtualRAM,
		},
		{
			Comment:    "1MB Physical RAM",
			NumPages:   0x100,
			MemoryType: MemoryTypePhysicalRAM,
		},
		{
			Comment:    "1MB I/O RAM",
			NumPages:   0x100,
			MemoryType: MemoryTypeIORAM,
		},
		{
			Comment:    "1MB Physical ROM",
			NumPages:   0x100,
			MemoryType: MemoryTypePhysicalROM,
		},
	},
}
