package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
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

type PhysicalPage struct {
	Buffer     [PageSize]byte
	MemoryType int
}

type PhysicalMemory struct {
	PhysicalPages    []PhysicalPage
	FreeVirtualPages *list.List
	UsedVirtualPages *list.List
}

type PhysicalMemoryRegion struct {
	Comment    string
	NumPages   uint32
	MemoryType int
}

func PhysicalMemory_Initialize(pmr []PhysicalMemoryRegion) *PhysicalMemory {
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Initialize", "Initializing physical memory")
	totalSize := pmr[len(pmr)-1].EndAddress
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Initialize",
		"Total size of physical memory is "+strconv.Itoa(int(totalSize)))
	pm := PhysicalMemory{
		PhysicalPages: make([]PhysicalPage, totalSize),
	}
	pm.FreeVirtualPages = list.New()
	pm.UsedVirtualPages = list.New()
	for i := 0; i < len(pmr); i++ {
		for j := pmr[i].StartPage; j <= pmr[i].EndAddress; j++ {
			pm.PhysicalPages[j].MemoryType = pmr[i].MemoryType
			pm.FreeVirtualPages.PushBack(j)
		}
	}
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Initialize", "Initialization completed")
	return &pm
}

func (pm *PhysicalMemory) Terminate() {
	RemoteLogging.LogEvent("INFO", "PhysicalMemory_Terminate", "Terminating physical memory")
	pm.FreeVirtualPages = nil
	pm.UsedVirtualPages = nil
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
