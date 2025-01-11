package VirtualMemory

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/MMUSupport/PhysicalMemory"
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
	"strconv"
)

type VirtualPage struct {
	PageFlags    uint64
	PhysicalPage uint32
}

type VMContainer struct {
	Swapper        MMUSupport.SwapperInterface
	PhysicalMemory PhysicalMemory.PhysicalMemory
	VirtualPages   []VirtualPage
	NumPages       uint32
	FreePages      *list.List
	UsedPages      *list.List
}

func VirtualMemory_Initiailize(
	pm PhysicalMemory.PhysicalMemory,
	swapper MMUSupport.SwapperInterface,
	numVirtPages uint32) (*VMContainer, error) {
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Initialize", "Initializing virtual memory")
	vmc := VMContainer{
		PhysicalMemory: pm,
		Swapper:        swapper,
		VirtualPages:   make([]VirtualPage, numVirtPages),
		NumPages:       numVirtPages,
		FreePages:      list.New(),
		UsedPages:      list.New(),
	}
	for i := uint32(0); i < numVirtPages; i++ {
		vmc.FreePages.PushBack(i)
	}
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Initialize", "Initialization completed")
	return &vmc, nil
}

func (vmc *VMContainer) Terminate() {
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Terminate", "Terminating virtual memory")
	vmc.FreePages = nil
	vmc.UsedPages = nil
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Terminate", "Termination completed")
}

func (vmc *VMContainer) AllocateVirtualPage() (uint32, error) {
	RemoteLogging.LogEvent("INFO", "AllocateVirtualPage", "Allocating virtual page")
	if vmc.FreePages.Len() == 0 {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPage", "No free pages")
		return 0, errors.New("No free virtual pages")
	}
	page := vmc.FreePages.Remove(vmc.FreePages.Front()).(uint32)
	vmc.VirtualPages[page].PhysicalPage = 0
	vmc.VirtualPages[page].PageFlags |= MMUSupport.PageIsActive | MMUSupport.PageIsOnDisk
	vmc.UsedPages.PushBack(page)
	RemoteLogging.LogEvent("INFO",
		"AllocateVirtualPage",
		"Allocation completed.  Page = "+strconv.Itoa(int(page)))
	return page, nil
}

func (vmc *VMContainer) ReturnVirtualPage(page uint32) error {
	RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Returning virtual page")
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsActive == 0 {
		RemoteLogging.LogEvent("ERROR", "ReturnVirtualPage", "Page is not active")
		return errors.New("Page is not active")
	}
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsOnDisk == MMUSupport.PageIsOnDisk {
		RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Page is on disk -- no physical page to free")
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsOnDisk
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsActive
		vmc.VirtualPages[page].PhysicalPage = 0
		vmc.FreePages.PushBack(page)
		vmc.UsedPages.Remove(vmc.UsedPages.Front())
		return nil
	}
	RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Return completed")
	return nil
}
