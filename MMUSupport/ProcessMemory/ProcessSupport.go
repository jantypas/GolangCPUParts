package ProcessMemory

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/MMUSupport/PhysicalMemory"
	"GolangCPUParts/MMUSupport/VirtualMemory"
	"time"
)

type ProcessObject struct {
	Name      string
	Args      []string
	PID       uint16
	PPID      uint16
	UID       uint32
	GID       uint32
	CreatedOn time.Time
	System    bool
	Memory    []uint32
}

type ProcessTable struct {
	Processes      map[uint16]ProcessObject
	NextPID        uint16
	Swapper        MMUSupport.SwapperInterface
	PhysicalMemory PhysicalMemory.PhysicalMemoryContainer
	VirtualMemory  VirtualMemory.VMContainer
}

func ProcessTable_Initialize(name string, sz uint32) (*ProcessTable, error) {
	pt := ProcessTable{}
	pt.Processes = make(map[uint16]ProcessObject)
	pt.NextPID = 0
	pt.Swapper.Filename = "/tmp/swap.swp"
	pmem, err := PhysicalMemory.PhysicalMemory_Initialize(name)
	if err != nil {
		return nil, err
	}
	pt.PhysicalMemory = *pmem
	vmem, err := VirtualMemory.VirtualMemory_Initiailize(*pmem, pt.Swapper, sz)
	pt.VirtualMemory = *vmem
	return &pt, nil
}

func (pt *ProcessTable) Terminate() {
	pt.PhysicalMemory.Terminate()
	pt.VirtualMemory.Terminate()
}
