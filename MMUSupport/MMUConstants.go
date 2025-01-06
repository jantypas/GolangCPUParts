package MMUSupport

import "container/list"

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
	PageIsLocked = 0x8

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

type MMUConfig struct {
	Swapper          SwapperInterface // The swapper that swaps pages in and out from disk
	NumVirtualPages  int              // Number of virtual memory pages
	NumPhysicalPages int              // Number of physical memory pages
}

type MMUStruct struct {
	MMUConfig         MMUConfig
	PhysicalMem       []byte
	VirtualMemory     []VirtualPage
	FreeVirtualPages  *list.List
	FreePhysicalPages *list.List
	UsedVirtualPages  *list.List
	UsedPhysicalPages *list.L
	LRUCache          *list.List
}
