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
)

type VMContainer struct {
	MemoryPages      []VMPage
	MemoryMap        []MemoryMap.MemoryMapRegion
	Swapper          *Swapper.SwapperContainer
	PhysicaklPMemory *PhysicalMemory.PhysicalMemoryContainer
	FreeVirtualPages *list.List
	UsedVirtualPages *list.List
	LRUCache         *list.List
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
	vmc.PhysicaklPMemory = pmc
	return &vmc, nil
}

func (vmc *VMContainer) Terminate() {
	vmc.Swapper.Terminate()
	vmc.PhysicaklPMemory.Terminate()
	vmc.MemoryPages = nil
	vmc.MemoryMap = nil
	vmc.Swapper = nil
	vmc.PhysicaklPMemory = nil
	vmc.FreeVirtualPages = nil
	vmc.UsedVirtualPages = nil
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
		ok := vmc.SwapOldPages()
		// Nope, try to make some byh swapping some out
		if !ok {
			// Not possible
			return nil, errors.New("Failed to allocate virtual pages")
		}
		// Try again
		if vmc.GetNumberFreePages() < num {
			return nil, errors.New("Failed to allocate virtual pages")
		}
	}
	// We have enoughh free pages now, let's allocate some
	lst := list.New()
	for i := uint32(0); i < num; i++ {
		elm := vmc.FreeVirtualPages.Front()
		page := elm.Value.(uint32)
		vmc.FreeVirtualPages.Remove(elm)
		vmc.UsedVirtualPages.PushBack(page)
		vmc.PageIsActive(page)
		vmc.PageIsOnDisk(page)
		lst.PushBack(page)
	}
	return lst, nil
}

func (vmc *VMContainer) ReturnNVirtualPages(pages *list.List) error {
	for page := pages.Front(); page != nil; page.Next() {
		page := page.Value.(uint32)
		vmc.PageIsNotActive(page)
		vmc.FreeVirtualPages.PushBack(page)
		for elm := vmc.UsedVirtualPages.Front(); elm != nil; elm = elm.Next() {
			if elm.Value.(uint32) == page {
				vmc.UsedVirtualPages.Remove(elm)
				break
			}
		}
	}
	return nil
}

func (vmc *VMContainer) SwapOutPage(page uint32) error {

}

func (vmc *VMContainer) SwapOldPages() bool {

}

func (vmc *VMContainer) SwapInPage(page uint32) error {

}

func (vmc *VMContainer) ReadPage(page uint32) ([]byte, error) {

}

func (vnc *VMContainer) WritePage(page uint32, buffer []byte) error {

}

func (vmc *VMContainer) ReadAddress(pid uint32, segment uint32, offset uint32) (uint64, error) {

}

func (vmc *VMContainer) WriteAddress(pid uint32, segment uint32, offset uint32, value uint64) error {

}
