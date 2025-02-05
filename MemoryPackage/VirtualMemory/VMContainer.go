package VirtualMemory

import (
	"GolangCPUParts/Configuration"
	"GolangCPUParts/MemoryPackage/MemoryMap"
	"GolangCPUParts/MemoryPackage/PhysicalMemory"
	"GolangCPUParts/MemoryPackage/Swapper"
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
)

const (
	PageStatus_Active = 0x0001
	PageStatus_OnDisk = 0x0002
	PageStatus_Locked = 0x0004

	MaxSwapPages = 4
)

type VMContainer struct {
	MemoryPages        []VMPage
	MemoryMap          []MemoryMap.MemoryMapRegion
	Swapper            *Swapper.SwapperContainer
	PhysicalPMemory    *PhysicalMemory.PhysicalMemoryContainer
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

func (vmc *VMContainer) PageIsActive(page uint32) {
	vmc.MemoryPages[page].Status |= PageStatus_Active
}
func (vmc *VMContainer) PageIsOnDisk(page uint32) {
	vmc.MemoryPages[page].Status |= PageStatus_OnDisk
}
func (vmc *VMContainer) PageIsLocked(page uint32) {
	vmc.MemoryPages[page].Status |= PageStatus_Locked
}
func (vmc *VMContainer) PageIsNotOnDisk(page uint32) {
	vmc.MemoryPages[page].Status &= ^PageStatus_OnDisk
}
func (vmc *VMContainer) PageIsNotActive(page uint32) {
	vmc.MemoryPages[page].Status &= ^PageStatus_Active
}
func (vmc *VMContainer) PageIsNotLocked(page uint32) {
	vmc.MemoryPages[page].Status &= ^PageStatus_Locked
}
func ListFindUint32(l *list.List, v uint32) *list.Element {
	for l := l.Front(); l != nil; l.Next() {
		if l.Value.(uint32) == v {
			return l
		}
	}
	return nil
}

func MoveFreeToUsed(freelst *list.List, usedlst *list.List, pg uint32) {
	elm := ListFindUint32(freelst, pg)
	freelst.Remove(elm)
	usedlst.PushBack(pg)
}

func MoveUsedToFree(freelst *list.List, usedlst *list.List, pg uint32) {
	elm := ListFindUint32(usedlst, pg)
	usedlst.Remove(elm)
	freelst.PushBack(pg)
}

func VirtualMemoryInitialize(
	cfg Configuration.ConfigObject,
	name string, vpages uint32) (*VMContainer, error) {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Initializing virtual memory")
	// See if the memory map is valid
	mr, ok := MemoryMap.ProductionMap[name]
	if !ok {
		return nil, errors.New("Invalid memory map")
	}
	// Create physical memory for our virtual memory
	pmc, err := PhysicalMemory.PhysicalMemoryInitialize(mr)
	if err != nil {
		return nil, err
	}
	vmc := VMContainer{}
	// Create the overall page map -- for each physical page, createa  virtual one
	totalPagesNeeded := pmc.GetTotalPages()
	vmc.MemoryPages = make([]VMPage, totalPagesNeeded)
	vmc.UsedVirtualPages = list.New()
	vmc.FreeVirtualPages = list.New()
	vmc.LRUCache = list.New()
	// Special handling for virtual region -- we need to know where these pages are
	// all other pages should be locked in memory since we can't swap them
	for i := uint32(0); i < totalPagesNeeded; i++ {
		if pmc.GetPageType(i) == MemoryMap.SegmentTypeVirtualRAM {
			vmc.FreeVirtualPages.PushBack(i)
			vmc.FreePhysicalMemory.PushBack(i)
			vmc.MemoryPages[i] = VMPage{
				VirutalPage:  i,
				PhysicalPage: i,
				Status:       0,
			}
		} else {
			vmc.MemoryPages[i] = VMPage{
				VirutalPage:  i,
				PhysicalPage: i,
				Status:       PageStatus_Locked | PageStatus_Active,
			}
		}
	}
	// Start the swapper
	vmc.Swapper, err = Swapper.Swapper_Initialize(cfg.SwapFileNames, &vmc)
	if err != nil {
		return nil, err
	}
	vmc.PhysicalPMemory = pmc
	return &vmc, nil
}

func (vmc *VMContainer) Terminate() {
	vmc.Swapper.Terminate()
	vmc.PhysicalPMemory.Terminate()
	vmc.MemoryPages = nil
	vmc.MemoryMap = nil
	vmc.Swapper = nil
	vmc.PhysicalPMemory = nil
	vmc.FreeVirtualPages = nil
	vmc.UsedVirtualPages = nil
	vmc.FreePhysicalMemory = nil
	vmc.UsedPhysicalMemory = nil
	vmc.LRUCache = nil
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
		vmc.SwapOldPages()
		// Try again
		if vmc.GetNumberFreePages() < num {
			return nil, errors.New("Failed to allocate virtual pages")
		}
	}
	// We have enough free pages now, let's allocate some
	lst := list.New()
	for i := uint32(0); i < num; i++ {
		// Allocate a physical page
		if vmc.FreePhysicalMemory.Len() == 0 {
			vmc.SwapOldPages()
		}
		if vmc.FreePhysicalMemory.Len() == 0 {
			return nil, errors.New("Failed to allocate virtual pages")
		}
		if vmc.FreeVirtualPages.Len() == 0 {
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
	return lst, nil
}

func (vmc *VMContainer) ReturnNVirtualPages(pages *list.List) error {
	for page := pages.Front(); page != nil; page.Next() {
		page := page.Value.(uint32)
		vmc.PageIsNotActive(page)
		MoveUsedToFree(vmc.UsedVirtualPages, vmc.FreeVirtualPages, page)
	}
	return nil
}

func (vmc *VMContainer) SwapOutPage(page uint32) error {
	if !vmc.IsPageActive(page) {
		return errors.New("Page is not active")
	}
	if vmc.IsPageOnDisk(page) {
		return errors.New("Page is already on disk")
	}
	vmc.PageIsOnDisk(page)
	pp := vmc.MemoryPages[page].PhysicalPage
	bp, err := vmc.PhysicalPMemory.ReadPhysicalPage(pp)
	if err != nil {
		return err
	}
	err = vmc.Swapper.SwapOutPage(page, bp)
	if err != nil {
		return err
	}
	MoveUsedToFree(vmc.UsedVirtualPages, vmc.FreeVirtualPages, page)
	return nil
}

func (vmc *VMContainer) SwapOldPages() {
	if vmc.UsedVirtualPages.Len() > MaxSwapPages &&
		vmc.LRUCache.Len() > MaxSwapPages {
		// We enough to swap out MaxSwapPages to make extra room
		// Swap out MaxSwapPages oldest pages
		for i := 0; i < MaxSwapPages; i++ {
			newPage := vmc.LRUCache.Back().Value.(uint32)
			err := vmc.SwapOutPage(newPage)
			if err != nil {
				panic(err)
			}
			vmc.LRUCache.Remove(vmc.LRUCache.Back())
		}
	} else {
		// Not enough room, just swap out the oldest page
		newPage := vmc.LRUCache.Back().Value.(uint32)
		err := vmc.SwapOutPage(newPage)
		if err != nil {
			panic(err)
		}
		vmc.LRUCache.Remove(vmc.LRUCache.Back())
	}
}

func (vmc *VMContainer) SwapInPage(page uint32) error {
	if !vmc.IsPageActive(page) {
		return errors.New("Page is not active")
	}
	if !vmc.IsPageOnDisk(page) {
		return errors.New("Page is not on disk")
	}
	if vmc.FreePhysicalMemory.Len() == 0 {
		vmc.SwapOldPages()
		if vmc.FreePhysicalMemory.Len() == 0 {
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
	return nil
}

func (vmc *VMContainer) ReadPage(page uint32) ([]byte, error) {
	if !vmc.IsPageActive(page) {
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
	return bp, nil
}

func (vmc *VMContainer) WritePage(page uint32, buffer []byte) error {
	if !vmc.IsPageActive(page) {
		return errors.New("Page is not active")
	}
	if !vmc.IsPageOnDisk(page) {
		err := vmc.SwapInPage(page)
		if err != nil {
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
	page := addr / PhysicalMemory.PageSize
	offset := addr % PhysicalMemory.PageSize
	buffer, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return 0, err
	}
	return buffer[offset], nil

}

func (vmc *VMContainer) WriteAddress(addr uint64, value byte) error {
	page := addr / PhysicalMemory.PageSize
	offset := addr % PhysicalMemory.PageSize
	bp, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return err
	}
	bp[offset] = value
	return vmc.WritePage(uint32(page), bp)
}
