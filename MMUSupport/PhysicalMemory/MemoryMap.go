package PhysicalMemory

const (
	MemoryTypeEmpty       = 0
	MemoryTypeVirtualRAM  = 1
	MemoryTypePhysicalRAM = 2
	MemoryTypePhysicalROM = 3
	MemoryTypeKernelRAM   = 4
	MemoryTypeKernelROM   = 5
	MemoryTypeIORAM       = 6
	MemoryTypeIOROM       = 7
	MemoryTypeBufferRAM   = 8
	PageSize              = 4096
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
}
