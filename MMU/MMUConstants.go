package MMU

import "time"

const (
	PageProtectionUserCanRead      = 0x1
	PageProtectionUserCanWrite     = 0x2
	PageProtectionUserCanExecute   = 0x4
	PageProtectionKUserNeedSystem  = 0x8
	PageProtectionKGroupCanRead    = 0x10
	PageProtectionKGroupCanWrite   = 0x20
	PageProtectionKGroupCanExecute = 0x40
	PageProtectionKWorldCanRead    = 0x80
	PageProtectionKWorldCanWrite   = 0x100
	PageProtectionKWorldCanExecute = 0x200

	PageProtectionMaskUser  = 0x000F
	PageProtectionMaskGroup = 0x00F0
	PageProtectionMaskWorld = 0x0F00

	PageSize = 4096

	PageIsActive = 0x1
	PageIsDirty  = 0x2
	PageIsOnDisk = 0x4

	ProtectionNeedRead    = 0x1
	ProtectionNeedWrite   = 0x2
	ProtectionNeedExecute = 0x4
	ProtectionHaveSystem  = 0x8

	VirtualErrorNoPages = 0x1

	ProcessStateSleeping     = 0
	ProcessStateWaitingForIO = 1
	ProcessStateWaitingToRun = 2
	ProcessStateRunning      = 3
)

type ProcessObject struct {
	Name      string
	Args      []string
	PID       int
	PPID      int
	GID       int
	System    bool
	State     int
	Memory    []int
	CreatedOn time.Time
}

type ProcessTable struct {
	ProcessList map[int]ProcessObject
	NextPID     int
	MMU         MMUStruct
}
type MMUConfig struct {
	Swapper          SwapperInterface // The swapper that swaps pages in and out from disk
	NumVirtualPages  int              // Number of virtual memory pages
	NumPhysicalPages int              // Number of physical memory pages
}

type VirtualPage struct {
	PhysicalPageID int
	Protection     int
	Flags          int
	ProcessID      int
	GroupID        int
	SegmentID      int
}

type MMUStruct struct {
	MMUConfig         MMUConfig
	TLB               []MMUTLB
	PhysicalMem       []byte
	VirtualMemory     []VirtualPage
	FreeVirtualPages  []int
	FreePhysicalPages []int
	UsedVirtualPages  []int
	UsedPhysicalPages []int
	LRUCache          []int
}
