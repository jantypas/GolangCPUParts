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

	MinFreePages    = 8
	MaxVirtualPages = 1024 * 1024
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
func (vmc *VMContainer) GetBuffer(page uint32) []byte {
	ppage := vmc.MemoryPages[page].PhysicalPage
	return vmc.PhysicalPMemory.Blocks[ppage].Buffer
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
	pmc, err := PhysicalMemory.PhysicalMemoryInitialize(&cfg, name)
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
	if numVPages < MinFreePages+1 {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "Not enough virtual memory pages found")
		return nil, errors.New("Not enough virtual memory pages found")
	}
	if numVPages > 1024*1024 {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryInitialize", "Too many virtual memory pages found")
		return nil, errors.New("Too many virtual memory pages found")
	}
	// Build the virtual paages into the lower 4GB (20-bits) of the map
	for i := 0; i < numVPages; i++ {
		vmc.MemoryPages[uint32(i)] =
			VMPage{uint32(i), uint32(i), 0}
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

func (vmc *VMContainer) AllocateVirtualPages(numPagse int) ([]uint32, error) {
	RemoteLogging.LogEvent("INFO", "AllocateVirtualPages", "Allocating "+strconv.Itoa(numPagse)+" virtual pages")
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPages", "VMContainer is nil")
		return nil, errors.New("VMContainer is nil")
	}
	if numPagse == 0 {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPages", "Invalid number of pages")
		return nil, errors.New("Invalid number of pages")
	}
	if numPagse > MaxVirtualPages {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPages", "Too many pages requested")
		return nil, errors.New("Too many pages requested")
	}
	if vmc.FreeVirtualPages.Len() < numPagse {
		RemoteLogging.LogEvent("ERROR", "AllocateVirtualPages", "Not enough free pages")
		return nil, errors.New("Not enough free pages")
	}
	lst := make([]uint32, numPagse)
	pageIdx := 0
	for vmc.FreeVirtualPages.Len() > 0 {
		elm := vmc.FreeVirtualPages.Front()
		vmc.FreeVirtualPages.Remove(elm)
		pgValue := elm.Value.(uint32)
		lst[pageIdx] = pgValue
		vmc.MemoryPages[pgValue] = VMPage{
			VirutalPage:  pgValue,
			PhysicalPage: 0,
			Status:       PageStatus_Active | PageStatus_OnDisk,
		}
		pageIdx++
	}
	return lst, nil
}

func (vmc *VMContainer) ReturnVirtualPages(pages []uint32) error {
	RemoteLogging.LogEvent("INFO", "ReturnVirtualPages", "Returning "+strconv.Itoa(len(pages))+" virtual pages")
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "ReturnVirtualPages", "VMContainer is nil")
		return errors.New("VMContainer is nil")
	}
	if len(pages) == 0 {
		RemoteLogging.LogEvent("ERROR", "ReturnVirtualPages", "Invalid number of pages")
		return errors.New("Invalid number of pages")
	}
	for _, pg := range pages {
		MoveUsedToFree(vmc.UsedVirtualPages, vmc.FreeVirtualPages, pg)
		MoveUsedToFree(vmc.UsedPhysicalMemory, vmc.FreePhysicalMemory, vmc.MemoryPages[pg].PhysicalPage)
		delete(vmc.MemoryPages, pg)
	}
	return nil
}

func (vmc *VMContainer) SwapOutPage(page uint32) error {
	RemoteLogging.LogEvent("INFO", "SwapOutPage", "Swapping out page "+strconv.Itoa(int(page)))
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "VMContainer is nil")
		return errors.New("VMContainer is nil")
	}
	if page > MaxVirtualPages {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Invalid page number")
		return errors.New("Invalid page number")
	}
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is not active")
		return errors.New("Page is not active")
	}
	if vmc.IsPageOnDisk(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is on disk")
		return errors.New("Page is on disk")
	}
	if vmc.IsPageLocked(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is locked")
		return errors.New("Page is locked")
	}
	bptr := vmc.GetBuffer(page)
	err := vmc.Swapper.SwapOutPage(page, bptr)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Failed to swap out page")
		return err
	}
	vmc.SetPageIsOnDisk(page)
	MoveUsedToFree(vmc.UsedPhysicalMemory, vmc.FreePhysicalMemory, vmc.MemoryPages[page].PhysicalPage)
	return nil
}

func (vmc *VMContainer) SwapInPage(page uint32) error {
	RemoteLogging.LogEvent("INFO", "SwapInPage", "Swapping in page "+strconv.Itoa(int(page)))
	// Make sure we don't have a null container
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "VMContainer is nil")
		return errors.New("VMContainer is nil")
	}
	// Make sure our page is in range
	if page > MaxVirtualPages {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Invalid page number")
		return errors.New("Invalid page number")
	}
	// If the page is not active, this is an error
	if !vmc.IsPageActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is not active")
		return errors.New("Page is not active")
	}
	// If page is not on disk, we can't swap it in
	if !vmc.IsPageOnDisk(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Page is not on disk")
		return errors.New("Page is not on disk")
	}
	// Do we have a free physical page to put this into?
	if vmc.FreePhysicalMemory.Len() == 0 {
		// No, try swapping out some old pages
		err := vmc.SwapOutOldPages()
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Failed to swap out old pages")
			return err
		}
		// We've got a page, so point the virtual page to it
		newPage := vmc.FreePhysicalMemory.Front().Value.(uint32)
		MoveFreeToUsed(vmc.FreePhysicalMemory, vmc.UsedPhysicalMemory, newPage)
		s := vmc.MemoryPages[page]
		s.PhysicalPage = newPage
		vmc.MemoryPages[page] = s
		// Page will now be in memory
		vmc.SetPageIsNotOnDisk(page)
		// Swap the page in
		err = vmc.Swapper.SwapInPage(page, vmc.GetBuffer(page))
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapOutPage", "Failed to swap out page")
			return err
		}
		return nil
	}
	return errors.New("Unknown swap in error")
}

func (vmc *VMContainer) SwapOutOldPages() error {
	RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Swapping out old pages")
	if vmc == nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutOldPages", "VMContainer is nil")
		return errors.New("VMContainer is nil")
	}
	// If the number of free pages > MinFreePages, nothing for us to do
	if vmc.FreePhysicalMemory.Len() > MinFreePages {
		return nil
	} else {
		for vmc.FreePhysicalMemory.Len() < MinFreePages {
			pg := vmc.FreePhysicalMemory.Front()
			val := pg.Value.(uint32)
			vmc.LRUCache.Remove(pg)
			err := vmc.Swapper.SwapOutPage(val, vmc.GetBuffer(val))
			if err != nil {
				RemoteLogging.LogEvent("ERROR", "SwapOutOldPages", "Failed to swap out page")
				return err
			}
		}
	}
	return nil
}

func (vmc *VMContainer) AttachRegion(reg int) ([]uint32, error) {
	RemoteLogging.LogEvent("INFO", "AttachRegion", "Attaching region "+strconv.Itoa(reg))
	blk, err := vmc.PhysicalPMemory.GetBlockByKey(reg)
	if err != nil {
		return nil, err
	}
	lst := make([]uint32, blk.NumPages)
	total := 0
	var i uint64
	for i = 0; i < uint64(blk.NumPages); i++ {
		nextAddr := blk.StartAddress + (i * PhysicalMemory.PhysicalPageSize)
		base := 1024 * 1024 * blk.Key
		lst[total] = uint32(base) + uint32(nextAddr/PhysicalMemory.PhysicalPageSize)
		total++
	}
	return lst, nil
}

func (vmc *VMContainer) ReadPage(page uint32) ([]byte, error) {
	if page > MaxVirtualPages {
		return nil, errors.New("Invalid page number")
	}
	_, ok := vmc.MemoryPages[page]
	if !ok {
		return nil, errors.New("Page not found")
	}
	if !vmc.IsPageActive(page) {
		return nil, errors.New("Page is not active")
	}
	if vmc.IsPageOnDisk(page) {
		err := vmc.SwapInPage(page)
		if err != nil {
			return nil, err
		}
	}
	vmc.LRUCache.PushFront(page)
	return vmc.GetBuffer(page), nil
}

func (vmc *VMContainer) WritePage(page uint32, buf []byte) error {
	if page > MaxVirtualPages {
		return errors.New("Invalid page number")
	}
	_, ok := vmc.MemoryPages[page]
	if !ok {
		return errors.New("Page not found")
	}
	if !vmc.IsPageActive(page) {
		return errors.New("Page is not active")
	}
	if vmc.IsPageOnDisk(page) {
		err := vmc.SwapInPage(page)
		if err != nil {
			return err
		}
	}
	blk := vmc.GetBuffer(page)
	copy(blk, buf)
	vmc.LRUCache.PushFront(page)
	return nil
}

func (vmc *VMContainer) ReadAddress(addr uint64) (byte, error) {
	page := addr / PhysicalMemory.PhysicalPageSize
	offset := addr % PhysicalMemory.PhysicalPageSize
	buf, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return 0, err
	}
	return buf[offset], nil
}

func (vmc *VMContainer) WriteAddress(addr uint64, value byte) error {
	page := addr / PhysicalMemory.PhysicalPageSize
	offset := addr % PhysicalMemory.PhysicalPageSize
	buf, err := vmc.ReadPage(uint32(page))
	if err != nil {
		return err
	}
	buf[offset] = value
	return vmc.WritePage(uint32(page), buf)
}
