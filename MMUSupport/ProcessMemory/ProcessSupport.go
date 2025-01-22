package ProcessMemory

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/MMUSupport/PhysicalMemory"
	"GolangCPUParts/MMUSupport/VirtualMemory"
	"GolangCPUParts/RemoteLogging"
	"debug/elf"
	"errors"
	"os"
	"time"
)

// ProcessObject
// Each active process has one of these structures associated with it
type ProcessObject struct {
	Name      string    // Process name -- ex: onyxmon
	Args      []string  // Any arguments associated with this process
	PID       uint16    // Process ID
	PPID      uint16    // Parent process ID
	UID       uint32    // User ID
	GID       uint32    // Group ID
	CreatedOn time.Time // When was this process started
	System    bool      // Does this process have system privileges
	State     uint64    // Process state flags
	Memory    []uint32  // Memory pages for this process
}

// ProcessTable
// All processes exist in this table
type ProcessTable struct {
	Processes      map[uint16]ProcessObject               // All process objects
	NextPID        uint16                                 // Next process ID
	Swapper        MMUSupport.SwapperInterface            // Reference our swapper
	PhysicalMemory PhysicalMemory.PhysicalMemoryContainer // Reference our physical memory
	VirtualMemory  VirtualMemory.VMContainer              // Reference our virtual memory
}

// ProcessTable_Initialize
// Create the process table itself.  NOTE, this does not actually create processes within the table
//
// name -- Name of the physical memory map
func ProcessTable_Initialize(name string) (*ProcessTable, error) {
	RemoteLogging.LogEvent("INFO", "ProcessTable_Initialize", "Initializing process table")
	pt := ProcessTable{}
	// Make the map for our ProcessObjects
	pt.Processes = make(map[uint16]ProcessObject)
	pt.NextPID = 0
	// Hook up our swapper
	pt.Swapper.Filename = "/tmp/swap.swp"
	// Hook up physical memory
	pmem, err := PhysicalMemory.PhysicalMemory_Initialize(name)
	if err != nil {
		return nil, err
	}
	pt.PhysicalMemory = *pmem
	// Hook up virtual memory
	numPages := pmem.ReturnListOfPageType(PhysicalMemory.MemoryTypeVirtualRAM)
	if numPages.Len() == 0 {
		return nil, errors.New("No virtual memory pages found")
	}
	// Start the swapper up
	err = pt.Swapper.Initialize()
	if err != nil {
		return nil, err
	}
	vmem, err := VirtualMemory.VirtualMemory_Initiailize(*pmem, pt.Swapper, uint32(numPages.Len()))
	if err != nil {
		return nil, err
	}
	pt.VirtualMemory = *vmem
	RemoteLogging.LogEvent("INFO", "ProcessTable_Initialize", "Process table initialized")
	return &pt, nil
}

func (pt *ProcessTable) Terminate() {
	RemoteLogging.LogEvent("INFO", "ProcessTable_Terminate", "Terminating process table")
	pt.PhysicalMemory.Terminate()
	pt.VirtualMemory.Terminate()
	err := pt.Swapper.Terminate()
	if err != nil {
		return
	}
	RemoteLogging.LogEvent("INFO", "ProcessTable_Terminate", "Process table terminated")
}

// BuildMemoryMap
// A process has a unique memory map based on what it needs.
// Given a requested number of virtual pages, and the other sections it needs,
// build a list of pages for that process.
func (pt *ProcessTable) BuildMemoryMap(sz uint32, mask uint) ([]uint32, error) {
	RemoteLogging.LogEvent("INFO", "ProcessTable_BuildMemoryMap", "Building memory map")
	// First make sure our region list only has ONE virtual region and it's the first one
	if pt.PhysicalMemory.Regions[0].MemoryType != PhysicalMemory.MemoryTypeVirtualRAM {
		RemoteLogging.LogEvent("ERROR", "ProcessTable_BuildMemoryMap", "First region is not virtual memory")
		return nil, errors.New("First region is not virtual memory")
	}
	if len(pt.PhysicalMemory.Regions) > 1 {
		for i := 1; i < len(pt.PhysicalMemory.Regions); i++ {
			if pt.PhysicalMemory.Regions[i].MemoryType == PhysicalMemory.MemoryTypeVirtualRAM {
				RemoteLogging.LogEvent("ERROR", "ProcessTable_BuildMemoryMap", "More than one virtual memory region found")
				return nil, errors.New("More than one virtual memory region found")
			}
		}
	}
	// First try to get the requested number of virtual pages
	lst, err := pt.VirtualMemory.AllocateNVirtualPages(sz)
	if err != nil {
		return nil, errors.New("Unable to allocate requested number of virtual memory pages")
	}
	// We received our virtual pages, place them at the bottom of our map
	// Look at the mask and determine what additional pages to add.
	// If we see a region that is not being asked for, ex: the region is a ROM region we don't need,
	// fill it with empty pages
}
