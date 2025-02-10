package VirtualMemory

import (
	"GolangCPUParts/Configuration"
	"GolangCPUParts/MemoryPackage/PhysicalMemory"
	"GolangCPUParts/MemoryPackage/Swapper"
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
	"fmt"
	"strconv"
)

const (
	PageStatus_Active = uint64(0x0000_0000_0000_0001)
	PageStatus_OnDisk = uint64(0x0000_0000_0000_0002)
	PageStatus_Locked = uint64(0x0000_0000_0000_0004)

	MaxSwapPages = 4
)

type VMContainer struct {
	MemoryPages        map[uint32]VMPage
	SystemDescriptor   Configuration.ConfigurationDescriptor
	Swapper            *Swapper.SwapperContainer
	PhysicalPMemory    *PhysicalMemory.PhysicalMemoryManager
	FreePhysicalMemory *list.List
	UsedPhysicalMemory *list.List
	FreeVirtualPages   *list.List
	UsedVirtualPages   *list.List
	LRUCache           *list.List
}

type VMPage struct {
	VirutalPage  uint32
	PhysicalPage uint32
	Status       uint64
}

func (vmc *VMContainer) IsPageActive(page uint32) bool {
	return vmc.MemoryPages[page].Status&PageStatus_Active != 0
}
func (vmc *VMContainer) IsPageOnDisk(page uint32) bool {
	return vmc.MemoryPages[page].Status&PageStatus_OnDisk != 0
}
func (vmc *VMContainer) IsPageLocked(page uint32) bool {
	return vmc.MemoryPages[page].Status&PageStatus_Locked != 0
}

func (vmc *VMContainer) SetPageActive(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status |= PageStatus_Active
	vmc.MemoryPages[page] = s
}
func (vmc *VMContainer) SetPageIsOnDisk(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status |= PageStatus_OnDisk
	vmc.MemoryPages[page] = s
}
func (vmc *VMContainer) SetPageIsLocked(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status |= PageStatus_Locked
	vmc.MemoryPages[page] = s
}
func (vmc *VMContainer) SetPageIsNotOnDisk(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status &= ^PageStatus_OnDisk
	vmc.MemoryPages[page] = s
}
func (vmc *VMContainer) SetPageIsNotActive(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status &= (^PageStatus_Active)
	vmc.MemoryPages[page] = s
}
func (vmc *VMContainer) PageIsNotLocked(page uint32) {
	s := vmc.MemoryPages[page]
	s.Status &= ^PageStatus_Locked
	vmc.MemoryPages[page] = s
}
func ListFindUint32(l *list.List, v uint32) *list.Element {
	for l := l.Front(); l != nil; l = l.Next() {
		if l.Value.(uint32) == v {
			return l
		}
	}
	return nil
}

func DebugList(name string, l *list.List) {
	fmt.Println("List: " + name)
	for l := l.Front(); l != nil; l = l.Next() {
		fmt.Println(l.Value.(uint32))
	}
}

func MoveFreeToUsed(freelst *list.List, usedlst *list.List, pg uint32) {
	elm := ListFindUint32(freelst, pg)
	if elm == nil {
		panic("Can't find page in free list")
		return
	}
	freelst.Remove(elm)
	usedlst.PushBack(pg)
}

func MoveUsedToFree(usedlist *list.List, freelst *list.List, pg uint32) {
	elm := ListFindUint32(usedlist, pg)
	if elm == nil {
		panic("Can't find page in used list")
		return
	}
	usedlist.Remove(elm)
	freelst.PushBack(pg)
}

func VirtualMemoryInitialize(
	cfg Configuration.ConfigObject,
	name string) (*VMContainer, error) {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Initializing virtual memory")
	sd := cfg.GetConfigByName(name)
	if sd == nil {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "Failed to get system descriptor")
		return nil, errors.New("Failed to get system descriptor by that name")
	}
	// Get our physical memory ranges
	mr := sd.Description.Memory
	if mr == nil {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "Failed to get memory ranges")
		return nil, errors.New("Failed to get memory ranges")
	}
	// Make our virtual memory container
	vmc := VMContainer{}
	// Try to start up physical memory
	pmc, err := PhysicalMemory.PhysicalMemoryInitialize(cfg, name)
	if err != nil {
		return nil, err
	}
	vmc.PhysicalPMemory = pmc
	// Set up our page menagement lists
	vmc.UsedPhysicalMemory = list.New()
	vmc.FreePhysicalMemory = list.New()
	vmc.UsedVirtualPages = list.New()
	vmc.FreeVirtualPages = list.New()
	vmc.LRUCache = list.New()
	vmc.SystemDescriptor = sd.Description
	vmc.MemoryPages = make(map[uint32]VMPage)
	// Get the virtual memory suitable pages from Physical Memory
	byType, err := pmc.GetBlockByType(PhysicalMemory.MemoryType_VirtualRAM)
	if err != nil {
		return nil, err
	}
	numVPages := byType.NumPages
	if numVPages == 0 {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "No suitable virtual memory pages found")
		return nil, errors.New("No suitable virtual memory pages found")
	}
	if numVPages > 1024*1024 {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "Too many virtual memory pages found")
		return nil, errors.New("Too many virtual memory pages found")
	}
	// Build the virtual paages into the lower 4GB (20-bits) of the map
	for i := 0; i < numVPages; i++ {
		vmc.MemoryPages[uint32(i)] =
			VMPage{uint32(i), i, PageStatus_Active | PageStatus_OnDisk}
	}
	// Finally start the swapper
	vmc.Swapper, err =
		Swapper.Swapper_Initialize(cfg.GetConfigurationSettings().SwapFileName)
	return &vmc, nil
}

func (vmc *VMContainer) Terminate() error {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryTerminate", "Terminating virtual memory")
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryTerminate", "VMContainer is nil")
		return errors.New("VMContainer is nil")
	}
	vmc.Swapper.Terminate()
	vmc.PhysicalPMemory.Terminate()
	vmc.MemoryPages = nil
	vmc.MemoryPages = nil
	vmc.Swapper = nil
	vmc.PhysicalPMemory = nil
	vmc.FreeVirtualPages = nil
	vmc.UsedVirtualPages = nil
	vmc.FreePhysicalMemory = nil
	vmc.UsedPhysicalMemory = nil
	vmc.LRUCache = nil
	return nil
}

func (vmc *VMContainer) GetNumberFreePages() uint32 {
	return uint32(vmc.FreeVirtualPages.Len())
}

func (vmc *VMContainer) GetNumberUsedPages() uint32 {
	return uint32(vmc.UsedVirtualPages.Len())
}

func (vmc *VMContainer) AllocateNVirtualPages(num uint32) (*list.List, error) {
	// Make sure we enough free pages
	if vmc.GetNumberFreePages() < num {
		RemoteLogging.LogEvent("INFO", "AllocateNVirtualPages", "Not enough free pages, swapping out pages")
		vmc.SwapOldPages()
		// Try again
		if vmc.GetNumberFreePages() < num {
			RemoteLogging.LogEvent("ERROR", "AllocateNVirtualPages", "Failed to swap out pages")
			return nil, errors.New("Failed to allocate virtual pages")
		}
	}
	// We have enough free pages now, let's allocate some
	lst := list.New()
	for i := uint32(0); i < num; i++ {
		// Allocate a physical page
		if vmc.FreePhysicalMemory.Len() == 0 {
			RemoteLogging.LogEvent("INFO", "AllocateVirtualPages", "No free physical pages, swapping out pages")
			vmc.SwapOldPages()
		}
		if vmc.FreePhysicalMemory.Len() == 0 {
			RemoteLogging.LogEvent("ERROR", "AllocateNVirtualPages", "Failed to swap out pages")
			return nil, errors.New("Failed to allocate virtual pages")
		}
		if vmc.FreeVirtualPages.Len() == 0 {
			RemoteLogging.LogEvent("ERROR", "AllocateNVirtualPages", "Failed to swap out pages")
			return nil, errors.New("Failed to allocate virtual pages")
		}
		// OK, we have a page available
		newPPage := vmc.FreePhysicalMemory.Back().Value.(uint32)
		MoveFreeToUsed(vmc.FreePhysicalMemory, vmc.UsedPhysicalMemory, newPPage)
		// Get a virtual page
		newVPage := vmc.FreeVirtualPages.Back().Value.(uint32)
		MoveFreeToUsed(vmc.FreeVirtualPages, vmc.UsedVirtualPages, newVPage)
		// Set up the virtual page
		vmc.MemoryPages[newVPage].PhysicalPage = newPPage
		vmc.MemoryPages[newVPage].Status = PageStatus_Active | PageStatus_OnDisk
		lst.PushBack(newVPage)
	}
	DebugList("Free Virtual Pages", vmc.FreeVirtualPages)
	DebugList("Used virtual pages", vmc.UsedVirtualPages)
	RemoteLogging.LogEvent("INFO", "AllocateNVirtualPages", "Allocated "+strconv.Itoa(int(num))+" virtual pages")
	return lst, nil
}

func (vmc *VMContainer) ReturnNVirtualPages(pages *list.List) error {
	DebugList("ReturnPages Used virtual pages", vmc.UsedVirtualPages)
	for page := pages.Front(); page != nil; page = page.Next() {
		newpage := page.Value.(uint32)
		ppage := vmc.MemoryPages[newpage].PhysicalPage
		MoveUsedToFree(vmc.UsedVirtualPages, vmc.FreeVirtualPages, newpage)
		MoveUsedToFree(vmc.UsedPhysicalMemory, vmc.FreePhysicalMemory, ppage)
		vmc.PageIsNotActive(newpage)
	}
	RemoteLogging.LogEvent("INFO", "ReturnNVirtualPages", "Returned "+strconv.Itoa(int(pages.Len()))+" virtual pages")
	return nil
}

func (vmc *VMContainer) SwapOutPage(page uint32) error {
	RemoteLogging.LogEvent("INFO", "SwapOutPage", "Swapping out page "+strconv.Itoa(int(page)))
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is not active")
		return errors.New("Page is not active")
	}
	if vmc.IsPageOnDisk(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is already on disk")
		return errors.New("Page is already on disk")
	}
	vmc.PageIsOnDisk(page)
	pp := vmc.MemoryPages[page].PhysicalPage
	bp, err := vmc.PhysicalPMemory.ReadPhysicalPage(pp)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Failed to read physical page")
		return err
	}
	err = vmc.Swapper.SwapOutPage(page, bp)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Failed to swap out page")
		return err
	}
	MoveUsedToFree(vmc.UsedVirtualPages, vmc.FreeVirtualPages, page)
	return nil
}

func (vmc *VMContainer) SwapOldPages() error {
	RemoteLogging.LogEvent("INFO", "SwapOldPages", "Swapping out old pages")
	if vmc.UsedVirtualPages.Len() > MaxSwapPages &&
		vmc.LRUCache.Len() > MaxSwapPages {
		// We enough to swap out MaxSwapPages to make extra room
		// Swap out MaxSwapPages oldest pages
		for i := 0; i < MaxSwapPages; i++ {
			newPage := vmc.LRUCache.Back().Value.(uint32)
			err := vmc.SwapOutPage(newPage)
			if err != nil {
				return err
			}
			vmc.LRUCache.Remove(vmc.LRUCache.Back())
		}
	} else {
		// Not enough room, just swap out the oldest page
		newPage := vmc.LRUCache.Back().Value.(uint32)
		err := vmc.SwapOutPage(newPage)
		if err != nil {
			return err
		}
		vmc.LRUCache.Remove(vmc.LRUCache.Back())
		RemoteLogging.LogEvent("INFO", "SwapOldPages", "Swapped out oldest page")
		return nil
	}
	return nil
}

func (vmc *VMContainer) SwapInPage(page uint32) error {
	RemoteLogging.LogEvent("INFO", "SwapInPage", "Swapping in page "+strconv.Itoa(int(page)))
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapInPage", "Page is not active")
		return errors.New("Page is not active")
	}
	if !vmc.IsPageOnDisk(page) {
		return errors.New("Page is not on disk")
	}
	if vmc.FreePhysicalMemory.Len() == 0 {
		vmc.SwapOldPages()
		if vmc.FreePhysicalMemory.Len() == 0 {
			RemoteLogging.LogEvent("ERROR", "SwapInPage", "Failed to swap in page")
			return errors.New("Failed to swap in page")
		}
	}
	// Get a free page
	newPPage := vmc.FreePhysicalMemory.Back().Value.(uint32)
	MoveFreeToUsed(vmc.FreePhysicalMemory, vmc.UsedPhysicalMemory, newPPage)
	// Copy the buffer via the swapper
	bp := make([]byte, PhysicalMemory.PageSize)
	err := vmc.Swapper.SwapInPage(newPPage, bp)
	if err != nil {
		return err
	}
	vmc.MemoryPages[page].PhysicalPage = newPPage
	vmc.PageIsNotOnDisk(page)
	RemoteLogging.LogEvent("INFO", "SwapInPage", "Swapped in page "+strconv.Itoa(int(page)))
	return nil
}

func (vmc *VMContainer) ReadPage(page uint32) ([]byte, error) {
	RemoteLogging.LogEvent("INFO", "ReadPage", "Reading page "+strconv.Itoa(int(page)))
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "ReadPage", "Page is not active")
		return nil, errors.New("Page is not active")
	}
	if !vmc.IsPageOnDisk(page) {
		err := vmc.SwapInPage(page)
		if err != nil {
			return nil, err
		}
	}
	pp := vmc.MemoryPages[page].PhysicalPage
	bp, err := vmc.PhysicalPMemory.ReadPhysicalPage(pp)
	if err != nil {
		return nil, err
	}
	vmc.LRUCache.PushBack(page)
	RemoteLogging.LogEvent("INFO", "ReadPage", "Read page "+strconv.Itoa(int(page)))
	return bp, nil
}

func (vmc *VMContainer) WritePage(page uint32, buffer []byte) error {
	RemoteLogging.LogEvent("INFO", "WritePage", "Writing page "+strconv.Itoa(int(page)))
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "WritePage", "Page is not active")
		errors.New("Page is not active")
	}
	if !vmc.IsPageOnDisk(page) {
		err := vmc.SwapInPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "WritePage", "Failed to swap in page")
			return err
		}
	}
	pp := vmc.MemoryPages[page].PhysicalPage
	err := vmc.PhysicalPMemory.WritePage(pp, buffer)
	if err != nil {
		return err
	}
	return nil
}

func (vmc *VMContainer) ReadAddress(addr uint64) (byte, error) {
	RemoteLogging.LogEvent("INFO", "ReadAddress", "Reading address "+strconv.Itoa(int(addr)))
	page := addr / PhysicalMemory.PageSize
	offset := addr % PhysicalMemory.PageSize
	buffer, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return 0, err
	}
	return buffer[offset], nil

}

func (vmc *VMContainer) WriteAddress(addr uint64, value byte) error {
	RemoteLogging.LogEvent("INFO", "WriteAddress", "Writing address "+strconv.Itoa(int(addr)))
	page := addr / PhysicalMemory.PageSize
	offset := addr % PhysicalMemory.PageSize
	bp, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return err
	}
	bp[offset] = value
	return vmc.WritePage(uint32(page), bp)
}
