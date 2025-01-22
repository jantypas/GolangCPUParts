package VirtualMemory

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/MMUSupport/PhysicalMemory"
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
	"strconv"
)

const (
	MinLRUCache      = 6
	LRUCacheTakeRate = 3
)

// VirtualPage
// Each virtual memory page is contained in one these structures
type VirtualPage struct {
	PageFlags    uint64 // The various page state flags (active, on disk etc)
	PhysicalPage uint32 // The physical page we're referring to
}

// VMContainer
// The container for all our virtual pages
type VMContainer struct {
	Swapper        MMUSupport.SwapperInterface            // A reference to our swapper service
	PhysicalMemory PhysicalMemory.PhysicalMemoryContainer // The physical memory service
	VirtualPages   []VirtualPage                          // Table of virtual pages
	NumPages       uint32                                 // # of virtual pages
	FreePages      *list.List                             // Free page list
	UsedPages      *list.List                             // Used page list
	LRUCache       *list.List                             // Least recently used page cache
}

// VirtualMemory_Initialize
// Initializes the virtual memory system.  Takes a reference to the
// physical memory container, our swapper interface, the number of virtual pages we want
func VirtualMemory_Initiailize(
	pm PhysicalMemory.PhysicalMemoryContainer,
	swapper MMUSupport.SwapperInterface,
	numVirtPages uint32) (*VMContainer, error) {
	// First create the virtual memory structure itself
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Initialize", "Initializing virtual memory")
	vmc := VMContainer{
		PhysicalMemory: pm,
		Swapper:        swapper,
		VirtualPages:   make([]VirtualPage, numVirtPages),
		NumPages:       numVirtPages,
		FreePages:      list.New(),
		UsedPages:      list.New(),
		LRUCache:       list.New(),
	}
	// For each virtual page, put it on the free page list
	for i := uint32(0); i < numVirtPages; i++ {
		vmc.FreePages.PushBack(i)
	}
	// Done
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Initialize", "Initialization completed")
	return &vmc, nil
}

// Terminate
// Termiante the virtual memory system
func (vmc *VMContainer) Terminate() {
	// Free our page lists
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Terminate", "Terminating virtual memory")
	vmc.FreePages = nil
	vmc.UsedPages = nil
	// Even though we're done, it's the caller's responsability to shut down the physical memory
	// system and stop the swapper
	RemoteLogging.LogEvent("INFO", "VirtualMemory_Terminate", "Termination completed")
}

// AllocateVirtualPage
// Allocate one page out of virtual memory -- returns the page numbger
func (vmc *VMContainer) AllocateVirtualPage() (uint32, error) {
	// If we're out of free pages, return the error
	RemoteLogging.LogEvent("INFO", "AllocateVirtualPage", "Allocating virtual page")
	if vmc.FreePages.Len() == 0 {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPage", "No free pages")
		return 0, errors.New("No free virtual pages")
	}
	// Find the free page off the free page list
	page := vmc.FreePages.Remove(vmc.FreePages.Front()).(uint32)
	vmc.VirtualPages[page].PhysicalPage = 0
	// Mark the retrieved page as active, and on disk
	vmc.VirtualPages[page].PageFlags |= MMUSupport.PageIsActive | MMUSupport.PageIsOnDisk
	// Mark the page on the used list
	vmc.UsedPages.PushBack(page)
	RemoteLogging.LogEvent("INFO",
		"AllocateVirtualPage",
		"Allocation completed.  Page = "+strconv.Itoa(int(page)))
	return page, nil
}

// ReturnVirtualPage
// Return a specific page
func (vmc *VMContainer) ReturnVirtualPage(page uint32) error {
	// If page is not active, this is an error
	RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Returning virtual page")
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsActive == 0 {
		RemoteLogging.LogEvent("ERROR", "ReturnVirtualPage", "Page is not active")
		return errors.New("Page is not active")
	}
	// If page is on disk, mark it as off-disk and non-active
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsOnDisk == MMUSupport.PageIsOnDisk {
		RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Page is on disk -- no physical page to free")
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsOnDisk
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsActive
		vmc.VirtualPages[page].PhysicalPage = 0
		// Put back on the free list
		vmc.FreePages.PushBack(page)
		vmc.UsedPages.Remove(vmc.UsedPages.Front())
		return nil
	} else {
		// Return the physical back to the physical page free pool
		page := vmc.VirtualPages[page].PhysicalPage
		RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Page is not on disk -- freeing physical page")
		err := vmc.PhysicalMemory.ReturnVirtualPage(page)
		// If we fail, return an error
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "ReturnVirtualPage", "Unable to return physical page")
			return errors.New("Unable to return physical page")
		}
		// Mark the virtual page as inactive and not on disk
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsOnDisk
		vmc.VirtualPages[page].PageFlags &= ^MMUSupport.PageIsActive
		vmc.VirtualPages[page].PhysicalPage = 0
		// Put back on the free list
		vmc.FreePages.PushBack(page)
		vmc.UsedPages.Remove(vmc.UsedPages.Front())
		RemoteLogging.LogEvent("INFO", "ReturnVirtualPage", "Physical page returned")
		return nil
	}
}

// SwapOutPage
// Swap a virtual page to disk
func (vnc *VMContainer) SwapOutPage(page uint32) error {
	// Get our physical page if any
	phyPage := vnc.VirtualPages[page].PhysicalPage
	// If the virtual page is on disk, we can't swap it out
	if vnc.VirtualPages[page].PageFlags&MMUSupport.PageIsOnDisk == MMUSupport.PageIsOnDisk {
		RemoteLogging.LogEvent("INFO", "SwapOutPage", "Page is on disk -- no physical page to swap out")
		return nil
	} else {
		// Swap the page out
		RemoteLogging.LogEvent("INFO", "SwapOutPage", "Page is not on disk -- swapping out physical page")
		err := vnc.Swapper.SwapOut(vnc.PhysicalMemory, phyPage)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Swap failed: "+err.Error())
			return err
		}
		// Return the physical page to the free pool
		lst := list.New()
		lst.PushBack(phyPage)
		vnc.PhysicalMemory.ReturnVirtualPages(lst)
		// Mark the page as on disk
		vnc.VirtualPages[page].PageFlags |= MMUSupport.PageIsOnDisk
		RemoteLogging.LogEvent("INFO", "SwapOutPage", "Swap completed")
		return nil
	}
}

func (vmc *VMContainer) SwapOutOldPages() error {
	if vmc.LRUCache.Len() > MinLRUCache {
		for i := 0; i < LRUCacheTakeRate; i++ {
			err := vmc.SwapOutPage(vmc.LRUCache.Back().Value.(uint32))
			if err != nil {
				return err
			}
		}
		return nil
	}
	if vmc.LRUCache.Len() > LRUCacheTakeRate {
		err := vmc.SwapOutPage(vmc.LRUCache.Back().Value.(uint32))
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("No pages to swap out")
}

func (vnc *VMContainer) SwapInPage(page uint32) error {
	// If we don't have room to swap a page in, swap some out
	if vnc.FreePages.Len() == 0 {
		RemoteLogging.LogEvent("INFO", "SwapInPage", "No free pages -- swapping out some pages")
		err := vnc.SwapOutOldPages()
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapInPage", "Unable to swap out old pages")
			return errors.New("Unable to swap out old pages")
		}
		RemoteLogging.LogEvent("INFO", "SwapInPage", "Swapping out completed")
		err = vnc.SwapInPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapInPage", "Unable to swap in page")
			return errors.New("Unable to swap in page")
		}
	}
	return nil
}

func (vmc *VMContainer) ReadAddress(addr uint64) (byte, error) {
	RemoteLogging.LogEvent("INFO", "ReadAddress", "Reading address "+strconv.Itoa(int(addr)))
	page := uint32(addr / PhysicalMemory.PageSize)
	offset := addr % PhysicalMemory.PageSize
	RemoteLogging.LogEvent("INFO",
		"ReadAddress",
		"Page is "+strconv.Itoa(int(page))+" and offset is "+strconv.Itoa(int(offset)))
	if page >= uint32(len(vmc.VirtualPages)) {
		RemoteLogging.LogEvent("ERROR", "ReadAddress", "Invalid virtual address")
		return 0, errors.New("Invalid virtual address")
	}
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsActive == 0 {
		RemoteLogging.LogEvent("ERROR", "ReadAddress", "Page is not active")
		return 0, errors.New("Page is not active")
	}
	if vmc.VirtualPages[page].PageFlags&MMUSupport.PageIsOnDisk == MMUSupport.PageIsOnDisk {
		RemoteLogging.LogEvent("INFO", "ReadAddress", "Reading from on disk page")
		err := vmc.SwapInPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "ReadAddress", "Unable to swap in page")
			return 0, errors.New("Unable to swap in page")
		}
		RemoteLogging.LogEvent("INFO", "ReadAddress", "Page swapped in")
		vmc.LRUCache.PushBack(page)
		return vmc.PhysicalMemory.PhysicalPages[vmc.VirtualPages[page].PhysicalPage].Buffer[offset], nil
	}
	return 0, errors.New("Unknown error")
}

func (vnc *VMContainer) WriteAddress(addr uint64, data byte) error {
	RemoteLogging.LogEvent("INFO",
		"WriteAddress",
		"Writing "+strconv.Itoa(int(data))+" to address "+strconv.Itoa(int(addr)))
	page := uint32(addr / PhysicalMemory.PageSize)
	offset := addr % MMUSupport.PageSize
	if vnc.VirtualPages[page].PageFlags&MMUSupport.PageIsActive == 0 {
		RemoteLogging.LogEvent("ERROR", "WriteAddress", "Page is not active")
		return errors.New("Page is not active")
	}
	if vnc.VirtualPages[page].PageFlags&MMUSupport.PageIsOnDisk == MMUSupport.PageIsOnDisk {
		err := vnc.SwapInPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "WriteAddress", "Unable to swap in page")
			return errors.New("Unable to swap in page")
		}
		vnc.LRUCache.PushBack(page)
		vnc.PhysicalMemory.PhysicalPages[vnc.VirtualPages[page].PhysicalPage].Buffer[offset] = data
		return nil
	}
	return errors.New("Unknown error")
}
