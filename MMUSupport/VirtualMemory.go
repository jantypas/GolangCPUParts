package MMUSupport

import (
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
)

const (
	LRUCacheSize   = 32
	LRUMinSwapSize = 3
)

type VirtualPage struct {
	PhysicalPageID int
	Flags          int
}

// A bund of useful utility functions
func (mmu *MMUStruct) PageIsOnDisk(page int) bool {
	return mmu.VirtualMemory[page].Flags&PageIsOnDisk == PageIsOnDisk
}
func (mmu *MMUStruct) SetPageIsOnDisk(page int) {
	mmu.VirtualMemory[page].Flags |= PageIsOnDisk
}
func (mmu *MMUStruct) ClearPageIsOnDisk(page int) {
	mmu.VirtualMemory[page].Flags &= ^PageIsOnDisk
}
func (mmu *MMUStruct) PageIsActive(page int) bool {
	return mmu.VirtualMemory[page].Flags&PageIsActive == PageIsActive
}
func (mmu *MMUStruct) SetPageIsActive(page int) {
	mmu.VirtualMemory[page].Flags |= PageIsActive
}
func (mmu *MMUStruct) ClearPageIsActive(page int) {
	mmu.VirtualMemory[page].Flags &= ^PageIsActive
}
func (mmu *MMUStruct) PageIsDirty(page int) bool {
	return mmu.VirtualMemory[page].Flags&PageIsDirty == PageIsDirty
}
func (mmu *MMUStruct) SetPageIsDirty(page int) {
	mmu.VirtualMemory[page].Flags |= PageIsDirty
}
func (mmu *MMUStruct) ClearPageIsDirty(page int) {
	mmu.VirtualMemory[page].Flags &= ^PageIsOnDisk
}

// VirtualMemoryInitialize -- Initialize the virtual memory system
func VirtualMemoryInitialize(mmu *MMUConfig) (MMUStruct, error) {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Initializing virtual memory system")
	// Set up the MMU Struct
	m := MMUStruct{MMUConfig: *mmu}
	// Make the tables
	err := m.MakeVirtualMemoryTable()
	if err != nil {
		return m, err
	}
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Virtual memory system initialized")
	return m, nil
}

// VirtualMemoryTerminate -- Terminate the virtual memory system
func (mmu *MMUStruct) VirtualMemoryTerminate() error {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryTerminate", "Terminating virtual memory system")
	// Stop the swapper
	err := mmu.MMUConfig.Swapper.Terminate()
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryTerminate", "Error terminating virtual memory system")
		return err
	}
	// Free Virtual Pages
	for i := 0; i < mmu.MMUConfig.NumVirtualPages; i++ {
		err := mmu.FreeVirtualPage(i)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "VirtualMemoryTerminate", "Error freeing virtual page")
		}
	}
	RemoteLogging.LogEvent("INFO", "VirtualMemoryTerminate", "Virtual memory system terminated")
	return nil
}

// MakeVirtualMemroyTable - Create the virtual emmroy table
// Build the virtual memory table strcutures.  Returns an error
func (mmu *MMUStruct) MakeVirtualMemoryTable() error {
	RemoteLogging.LogEvent("INFO", "MakeVirtualMemoryTable", "Creating virtual memory table")
	mmu.VirtualMemory = make([]VirtualPage, mmu.MMUConfig.NumVirtualPages)
	mmu.FreeVirtualPages = list.New()
	mmu.UsedVirtualPages = list.New()
	mmu.FreePhysicalPages = list.New()
	mmu.UsedPhysicalPages = list.New()
	for i := 0; i < mmu.MMUConfig.NumVirtualPages; i++ {
		mmu.FreeVirtualPages.PushBack(i)
	}
	for i := 0; i < mmu.MMUConfig.NumPhysicalPages; i++ {
		mmu.FreePhysicalPages.PushBack(i)
	}
	err := mmu.MMUConfig.Swapper.Initialize()
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "MakeVirtualMemoryTable", "Error initializing virtual memory table")
		return err
	}
	RemoteLogging.LogEvent("INFO", "MakeVirtualMemoryTable", "Virtual memory table created")
	return nil
}

// AllocateNewVirtualPageNoSwap -- Allocate a new virtual page from the virtual page table
// Return the page table and an error
func (mmu *MMUStruct) AllocateNewVirtualPageNoSwap() (int, int, error) {
	RemoteLogging.LogEvent("INFO", "AllocateNewVirtualPageNoSwap", "Allocating new virtual page")
	// See if we're out of virtual pages
	if mmu.FreeVirtualPages.Len() == 0 {
		RemoteLogging.LogEvent("ERROR", "AllocateNewVirtualPageNoSwap", "No free virtual pages")
		return 0, VirtualErrorNoPages, errors.New("no virtual pages")
	}
	// Get a free virtual page
	pageID := mmu.FreeVirtualPages.Front().Value.(int)
	mmu.SetPageIsActive(pageID)
	mmu.SetPageIsOnDisk(pageID)
	mmu.VirtualMemory[pageID].PhysicalPageID = -1
	RemoteLogging.LogEvent("INFO", "AllocateNewVirtualPageNoSwap", "Virtual page allocated")
	return pageID, 0, nil
}

// FreeVirtualPage -- Free a virtual and its associated physical page
// Takes a page ID and returns an error
func (mmu *MMUStruct) FreeVirtualPage(page int) error {
	RemoteLogging.LogEvent("INFO", "FreeVirtualPage", "Freeing virtual page")
	// Make sure page is a valid page
	if page > mmu.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "FreeVirtualPage", "Invalid virtual page")
		return errors.New("invalid virtual page")
	}
	// If the page is not already active
	if !mmu.PageIsActive(page) {
		RemoteLogging.LogEvent("ERROR", "FreeVirtualPage", "Page is already free")
		return errors.New("page is already free")
	}
	// Is the page on disk
	if mmu.PageIsOnDisk(page) {
		// NO physical page to free, just free the virtual page
		mmu.ClearPageIsActive(page)
		return nil
	} else {
		// We have a physical page to free
		err := mmu.ReturnPhysicalPage(mmu.VirtualMemory[page].PhysicalPageID)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "FreeVirtualPage", "Error returning physical page")
			return err
		}
		mmu.ClearPageIsActive(page)
		RemoteLogging.LogEvent("INFO", "FreeVirtualPage", "Virtual page freed")
		return nil
	}
}

// SwapOutPhysicalPage -- Swap a physical page out to disk
// Takes a virtual page ID and returns an error
func (mmu *MMUStruct) SwapOutPhysicalPage(page int) error {
	RemoteLogging.LogEvent("INFO", "SwapOutPhysicalPage", "Swapping out physical page")
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "SwapOutPhysicalPage", "Invalid virtual page")
		return errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPhysicalPage", "Page is already free")
		return errors.New("page is already free")
	}
	// Make sure page isn't swapped already
	if mmu.PageIsOnDisk(page) {
		RemoteLogging.LogEvent("ERROR", "SwapOutPhysicalPage", "Page is already on disk")
		return errors.New("page is already on disk")
	}
	// Find the physical page
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID
	// Swap it out
	err := mmu.MMUConfig.Swapper.SwapOut(physicalPage)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPhysicalPage", "Error swapping out physical page")
		return err
	}
	// Make the page as swapped out
	mmu.SetPageIsOnDisk(page)
	// Return the physical page
	err = mmu.ReturnPhysicalPage(physicalPage)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapOutPhysicalPage", "Error returning physical page")
		return err
	}
	mmu.VirtualMemory[page].PhysicalPageID = -1
	RemoteLogging.LogEvent("INFO", "SwapOutPhysicalPage", "Physical page swapped out")
	return nil
}

// SwapInPhysicalPage -- Swap a physical page in from disk
// Takes a virtual page ID and returns an error
func (mmu *MMUStruct) SwapInPhysicalPage(page int) error {
	RemoteLogging.LogEvent("INFO", "SwapInPhysicalPage", "Swapping in physical page")
	// Make sure the page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Invalid virtual page")
		return errors.New("invalid virtual page")
	}
	// Make sure the page is already active
	if !mmu.PageIsActive(page) {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Page is already free")
		return errors.New("page is already free")
	}
	// Make sure the page is on disk to swap in
	if !mmu.PageIsOnDisk(page) {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Page is not on disk")
		return errors.New("page is not on disk")
	}
	// Try to get a physical page to swap into
	page, err := mmu.AllocateNewPhysicalPage()
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Error allocating physical page")
		return err
	}
	// Set up the physical page
	mmu.VirtualMemory[page].PhysicalPageID = page
	mmu.ClearPageIsOnDisk(page)
	err = mmu.MMUConfig.Swapper.SwapIn(page)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Error swapping in physical page")
		return err
	}
	RemoteLogging.LogEvent("INFO", "SwapInPhysicalPage", "Physical page swapped in")
	return nil
}

func (mmu *MMUStruct) UpdateLRU(page int) {
	mmu.LRUCache.PushBack(page)
	for mmu.LRUCache.Len() > LRUCacheSize {
		mmu.LRUCache.Remove(mmu.LRUCache.Front())
	}
}

func (mmu *MMUStruct) SwapOutOldPages() error {
	RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Swapping out old pages")
	for mmu.LRUCache.Len() > LRUMinSwapSize {
		page := mmu.LRUCache.Front().Value.(int)
		err := mmu.SwapOutPhysicalPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapOutOldPages", "Error swapping out old page")
			return err
		}
		mmu.LRUCache.Remove(mmu.LRUCache.Front())
	}
	RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Old pages swapped out")
	return nil
}

func (mmu *MMUStruct) TryPageSwap(page int) error {
	RemoteLogging.LogEvent("INFO", "TryPageSwap", "Trying to swap page")
	count := 2
	for count != 0 {
		result := mmu.SwapInPhysicalPage(page)
		if result == nil {
			return nil
		}
		result2 := mmu.SwapOutOldPages()
		if result2 != nil {
			return result2
		}
		count--
	}
	RemoteLogging.LogEvent("INFO", "TryPageSwap", "TryPageSwap complete")
	return errors.New("page swap failed")
}

func (mmu *MMUStruct) WriteVirtualPage(page int, buffer []byte) error {
	RemoteLogging.LogEvent("INFO", "WriteVirtualPage", "Writing virtual page")
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "WriteVirtualPage", "Invalid virtual page")
		return errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		RemoteLogging.LogEvent("ERROR", "WriteVirtualPage", "Page is not active")
		return errors.New("page is not active")
	}
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID
	copy(mmu.PhysicalMem[physicalPage*PageSize:physicalPage*PageSize+PageSize], buffer)
	mmu.SetPageIsDirty(page)
	mmu.UpdateLRU(page)
	RemoteLogging.LogEvent("INFO", "WriteVirtualPage", "Virtual page written")
	return nil
}

func (mmu *MMUStruct) ReadVirtualPage(
	owner int, group int,
	mode int, seg int, page int) ([]byte, error) {
	RemoteLogging.LogEvent("INFO", "ReadVirtualPage", "Reading virtual page")
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "ReadVirtualPage", "Invalid virtual page")
		return nil, errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		RemoteLogging.LogEvent("ERROR", "ReadVirtualPage", "Page is not active")
		return nil, errors.New("page is not active")
	}
	// If page isn't in memory, bring it in
	if mmu.PageIsOnDisk(page) {
		err := mmu.TryPageSwap(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "ReadVirtualPage", "Error trying to swap page")
			return nil, err
		}
	}
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID
	mmu.SetPageIsDirty(page)
	mmu.UpdateLRU(page)
	RemoteLogging.LogEvent("INFO", "ReadVirtualPage", "Virtual page read")
	return mmu.PhysicalMem[physicalPage*PageSize : physicalPage*PageSize+PageSize], nil
}

func (mmu *MMUStruct) FreeBulkPages(pages []uint) error {
	RemoteLogging.LogEvent("INFO", "FreeBulkPages", "Freeing bulk pages")
	for _, page := range pages {
		err := mmu.FreeVirtualPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "FreeBulkPages", "Error freeing page")
			return err
		}
	}
	RemoteLogging.LogEvent("INFO", "FreeBulkPages", "Bulk pages freed")
	return nil
}

func (mmu *MMUStruct) AllocateBulkPages(desiredPages uint) ([]uint, error) {
	RemoteLogging.LogEvent("INFO", "AllocateBulkPages", "Allocating bulk pages")
	lst := make([]int, 0)
	var i uint
	for i = 0; i < desiredPages; i++ {
		page, _, err := mmu.AllocateNewVirtualPageNoSwap()
		if err != nil {
			return nil, err
		}
		lst = append(lst, page)
	}
	RemoteLogging.LogEvent("INFO", "AllocateBulkPages", "Bulk pages allocated")
	return lst, nil
}
