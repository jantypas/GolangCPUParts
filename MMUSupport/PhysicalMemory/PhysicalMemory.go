package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
	"bytes"
	"container/list"
	"errors"
	"go/ast"
	"strconv"
)

const (
	MemoryTypeEmpty       = 0
	MemoryTypeKernel      = 1
	MemoryTypeIO          = 2
	MemoryTypeROM         = 3
	MemoryTypeVirtualRAM  = 4
	MemoryTypePhysicalRAM = 5

	PageSize = 4096
)

// PhysicalPage
// For every physical page we manage, we keep one of these structures
type PhysicalPage struct {
	Buffer     []byte
	MemoryType int
}

// PhysicalMemory
// The physical memory container keeps track of all our physical pages
type PhysicalMemory struct {
	PhysicalPages       []PhysicalPage // The pages themselves
	VirtualPageList     *list.List     // The list of pages we can give to virtual memory
	KernelPageList      *list.List     // The list of pages that are kernel pages
	IOPageList          *list.List     // The list of pages for I/O
	ROMPageList         *list.List     // The list of pages that we use for ROM
	PhysicalRAMPageList *list.List     // Physical pages (can't be virtualized)
	EmtpyPage           *list.List     // Empty pages
	FreeRAMPages        *list.List     //
	UsedRAMPages        *list.List
}

// PhysicalMemoryRegion
// For every block of memory in a memory map, we define a region.
// It tells use the type of memory and the number of pages for that type
type PhysicalMemoryRegion struct {
	Comment    string
	NumPages   uint32
	MemoryType int
}

// PhysicalMemory_Initialize
// Given a memory map, creates the physical memory structure and returns a pointer to it.
func PhysicalMemory_Initialize(mapname string) *PhysicalMemory {
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Initialize", "Initializing physical memory")
	pmr, ok := MemoryMapTable[mapname]
	if !ok {
		RemoteLogging.LogEvent("ERROR",
			"PhysicalMemory_Initialize", "Unable to find memory map with name "+mapname)
		return nil
	}
	// Compute total page size
	totalSize := 0
	for i := 0; i < len(pmr); i++ {
		totalSize += int(pmr[i].NumPages)
	}
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Initialize",
		"Total size of physical memory is "+strconv.Itoa(int(totalSize)))
	// Build the page table
	pm := PhysicalMemory{
		PhysicalPages:       make([]PhysicalPage, totalSize),
		VirtualPageList:     list.New(),
		KernelPageList:      list.New(),
		IOPageList:          list.New(),
		ROMPageList:         list.New(),
		PhysicalRAMPageList: list.New(),
		FreeRAMPages:        list.New(),
		UsedRAMPages:        list.New(),
	}
	// Create buffers for each non-empty page
	currentPage := 0
	for i := 0; i < len(pmr); i++ {
		switch pm.PhysicalPages[i].MemoryType {
		case MemoryTypeKernel:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.KernelPageList.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].Buffer = make([]byte, PageSize)
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypeKernel
			}
			break
		case MemoryTypeIO:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.IOPageList.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypeIO
			}
			break
		case MemoryTypeROM:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.ROMPageList.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].Buffer = make([]byte, PageSize)
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypeROM
			}
			break
		case MemoryTypeVirtualRAM:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.VirtualPageList.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].Buffer = make([]byte, PageSize)
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypeVirtualRAM
			}
			break
		case MemoryTypePhysicalRAM:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.PhysicalRAMPageList.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].Buffer = make([]byte, PageSize)
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypePhysicalRAM
				pm.FreeRAMPages.PushBack(uint64(currentPage))
			}
			break
		case MemoryTypeEmpty:
			for j := 0; j < int(pmr[i].NumPages); j++ {
				pm.FreeRAMPages.PushBack(uint64(currentPage))
				pm.PhysicalPages[currentPage].MemoryType = MemoryTypeEmpty
			}
			currentPage++
			break
		}
	}
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Initialize", "Initialization completed")
	return &pm
}

// GetPages
// Return the list of pages given a type
func (pm *PhysicalMemory) GetPages(pagetype int) *list.List {
	switch pagetype {
	case MemoryTypeKernel:
		return pm.KernelPageList
	case MemoryTypeIO:
		return pm.IOPageList
	case MemoryTypeROM:
		return pm.ROMPageList
	case MemoryTypeVirtualRAM:
		return pm.VirtualPageList
	case MemoryTypePhysicalRAM:
		return pm.PhysicalRAMPageList
	case MemoryTypeEmpty:
		return nil
	}
	return nil
}

// Terminate
// Terminate the physical memory system
func (pm *PhysicalMemory) Terminate() {
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Terminate", "Terminating physical memory")
	for i := 0; i < len(pm.PhysicalPages); i++ {
		if pm.PhysicalPages[i].MemoryType != MemoryTypeEmpty {
			pm.PhysicalPages[i].Buffer = nil
		}
		pm.IOPageList = nil
		pm.KernelPageList = nil
		pm.ROMPageList = nil
		pm.PhysicalRAMPageList = nil
		pm.FreeRAMPages = nil
		pm.UsedRAMPages = nil
		pm.VirtualPageList = nil
		pm.PhysicalPages = nil
	}
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Terminate", "Termination completed")
}

func (pm *PhysicalMemory) AllocateVirtualPages(numPage uint32) (*list.List, error) {
	RemoteLogging.LogEvent("INFO",
		"AllocateVirtualPages",
		"Allocating "+strconv.Itoa(int(numPage))+" virtual pages")
	if numPage > uint32(pm.FreeVirtualPages.Len()) {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPages", "Not enough physical free pages")
		return &list.List{}, errors.New("Not enough physical free pages")
	}
	lst := list.New()
	for i := 0; i < int(numPage); i++ {
		lst.PushBack(pm.FreeVirtualPages.Remove(pm.FreeVirtualPages.Front()).(uint64))
	}
	RemoteLogging.LogEvent("INFO", "AllocateVirtualPages", "Allocation completed")
	return lst, nil
}

func (pm *PhysicalMemory) ReturnVirtualPages(lst *list.List) {
	RemoteLogging.LogEvent("INFO",
		"ReturnVirtualPages",
		"Returning "+strconv.Itoa(lst.Len())+" virtual pages")
	for e := lst.Front(); e != nil; e = e.Next() {
		pm.FreeVirtualPages.PushBack(e.Value)
		pm.UsedVirtualPages.Remove(e)
	}
	RemoteLogging.LogEvent("INFO", "ReturnVirtualPages", "Return completed")
}

func (pm *PhysicalMemory) ReadAddress(addr uint64) (byte, error) {
	RemoteLogging.LogEvent("INFO",
		"Physical ReadAddress",
		"Reading address "+strconv.Itoa(int(addr)))
	page := addr / PageSize
	offset := addr % PageSize
	RemoteLogging.LogEvent("INFO",
		"Physical ReadAddress",
		"Page is "+strconv.Itoa(int(page))+" and offset is "+strconv.Itoa(int(offset)))
	if page >= uint64(len(pm.PhysicalPages)) {
		RemoteLogging.LogEvent("ERROR",
			"Physical ReadAddress",
			"Invalid physical address")
		return 0, errors.New("Invalid physical address")
	}
	switch pm.PhysicalPages[page].MemoryType {
	case MemoryTypeIO:
		RemoteLogging.LogEvent("INFO", "Physical ReadAddress", "Reading from IO not implemented")
		return 0, errors.New("Reading from IO not implemented")
	case MemoryTypeROM:
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeKernel:
		RemoteLogging.LogEvent("INFO", "Physical ReadAddress", "Reading from kernel not implemented")
		return pm.PhysicalPages[page].Buffer[offset], nil
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent("ERROR", "Physical eadAddress", "Reading from empty memory")
		return 0, errors.New("Reading from empty memory")
	}
	RemoteLogging.LogEvent("ERROR", "Physical ReadAddress", "Unknown memory type")
	return 0, errors.New("Unknown memory type")
}

func (pm *PhysicalMemory) WriteAddress(addr uint64, data byte) error {
	RemoteLogging.LogEvent("INFO",
		"Physical WriteAddress",
		"Reading address "+strconv.Itoa(int(addr)))
	page := addr / PageSize
	offset := addr % PageSize
	RemoteLogging.LogEvent("INFO",
		"Physical WriteAddress",
		"Page is "+strconv.Itoa(int(page))+" and offset is "+strconv.Itoa(int(offset)))
	if page >= uint64(len(pm.PhysicalPages)) {
		RemoteLogging.LogEvent("ERROR", "Physical WriteAddress", "Invalid physical address")
		return errors.New("Invalid physical address")
	}
	switch pm.PhysicalPages[page].MemoryType {
	case MemoryTypeIO:
		RemoteLogging.LogEvent("INFO", "Physical WriteAddress", "Reading from IO not implemented")
		return errors.New("Reading from IO not implemented")
	case MemoryTypeROM:
		RemoteLogging.LogEvent("ERROR", "Physical WriteAddress", "Can't write to ROM")
		return errors.New("Can't write to ROM")
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeKernel:
		RemoteLogging.LogEvent("INFO", "Physical WriteAddress", "Writing completed")
		return nil
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent("ERROR", "WriteAddress", "Writing to empty memory")
		return errors.New("Reading from empty memory")
	}
	RemoteLogging.LogEvent("ERROR", "WriteAddress", "Unknown memory type")
	return errors.New("Unknown memory type")
}
