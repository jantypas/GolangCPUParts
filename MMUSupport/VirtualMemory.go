package MMUSupport

import (
	"GolangCPUParts/RemoteLogging"
	"errors"
)

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
func (mmu *MMUStruct) pruneLRUCache(value int) {
	// Create a new slice to store the result
	result := make([]int, 0, 64)

	for _, v := range mmu.LRUCache {
		if v != value {
			result = append(result, v) // Add only items not equal to the value
		}
	}
	mmu.LRUCache = make([]int, 0, 64)
	mmu.LRUCache = append(mmu.LRUCache, value)
	mmu.LRUCache = append(mmu.LRUCache, result...)
}
func (mmu *MMUStruct) CheckPermissionsOk(mode int, mask int, prot int) bool {
	finalMask := 0

	switch mask {
	case PageProtectionMaskUser:
		finalMask = prot & PageProtectionMaskUser
		break
	case PageProtectionMaskGroup:
		finalMask = (prot & PageProtectionMaskGroup) >> 4
		break
	case PageProtectionMaskWorld:
		finalMask = (prot & PageProtectionMaskWorld) >> 8
		break
	}
	return mode&finalMask == mode
}

// VirtualMemoryInitialize -- Initialize the virtual memory system
func VirtualMemoryInitialize(mmu *MMUConfig) (MMUStruct, error) {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Initializing virtual memory system")
	m := MMUStruct{
		MMUConfig: *mmu,
	}
	err := m.MakeVirtualMemoryTable()
	if err != nil {
		return m, err
	}
	RemoteLogging.LogEvent("INFO", "VirtualMemoryInitialize", "Virtual memory system initialized")
	return m, nil
}

// VirtualMemoryTerminate -- Terminate the virtual memoryh system
func (mmu *MMUStruct) VirtualMemoryTerminate() error {
	RemoteLogging.LogEvent("INFO", "VirtualMemoryTerminate", "Terminating virtual memory system")
	err := mmu.MMUConfig.Swapper.Terminate()
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "VirtualMemoryTerminate", "Error terminating virtual memory system")
		return err
	}
	RemoteLogging.LogEvent("INFO", "VirtualMemoryTerminate", "Virtual memory system terminated")
	return nil
}

// MakeVirtualMemroyTable - Create the virtual emmroy table
// Build the virtual memory table strcutures.  Returns an error
func (mmu *MMUStruct) MakeVirtualMemoryTable() error {
	RemoteLogging.LogEvent("INFO", "MakeVirtualMemoryTable", "Creating virtual memory table")
	mmu.VirtualMemory = make([]VirtualPage, mmu.MMUConfig.NumVirtualPages)
	mmu.FreeVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	mmu.UsedVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	for i := 0; i < mmu.MMUConfig.NumVirtualPages; i++ {
		mmu.FreeVirtualPages = append(mmu.FreeVirtualPages, i)
	}
	mmu.LRUCache = make([]int, 64)
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
func (mmu *MMUStruct) AllocateNewVirtualPageNoSwap(
	owner int,
	group int,
	protect int,
	seg int) (int, int, error) {
	RemoteLogging.LogEvent("INFO", "AllocateNewVirtualPageNoSwap", "Allocating new virtual page")
	if len(mmu.FreeVirtualPages) == 0 {
		RemoteLogging.LogEvent("ERROR", "AllocateNewVirtualPageNoSwap", "No free virtual pages")
		return 0, VirtualErrorNoPages, errors.New("no virtual pages")
	}
	pageID := mmu.FreeVirtualPages[0]
	mmu.FreeVirtualPages = mmu.FreeVirtualPages[1:]
	mmu.UsedVirtualPages = append(mmu.UsedVirtualPages, pageID)
	mmu.VirtualMemory[pageID].Protection = protect
	mmu.VirtualMemory[pageID].ProcessID = owner
	mmu.VirtualMemory[pageID].SegmentID = seg
	mmu.VirtualMemory[pageID].GroupID = group
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
	err := mmu.MMUConfig.Swapper.SwapOut(physicalPage, mmu.PhysicalMem[physicalPage*PageSize:])
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
	err = mmu.MMUConfig.Swapper.SwapIn(page, mmu.PhysicalMem[page*PageSize:])
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "SwapInPhysicalPage", "Error swapping in physical page")
		return err
	}
	RemoteLogging.LogEvent("INFO", "SwapInPhysicalPage", "Physical page swapped in")
	return nil
}

func (mmu *MMUStruct) SwapOutOldPages() error {
	RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Swapping out old pages")
	swapOutList := make([]int, 0)
	if len(mmu.LRUCache) > 3 {
		RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Swapping out at least 3 old pages")
		// We can get at least three pages to swpa out
		for i := 0; i < 3; i++ {
			swapOutList = append(swapOutList, mmu.LRUCache[len(mmu.LRUCache)-1])
			mmu.LRUCache = mmu.LRUCache[:len(mmu.LRUCache)-1]
		}
	} else {
		RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Swapping out at most 1 old pages")
		if len(mmu.LRUCache) > 0 {
			// We can't get three pages, do one
			swapOutList = append(swapOutList, mmu.LRUCache[len(mmu.LRUCache)-1])
			mmu.LRUCache = mmu.LRUCache[:len(mmu.LRUCache)-1]
		} else {
			RemoteLogging.LogEvent("ERROR", "SwapOutOldPages", "No old pages to swap out")
			// There are no pages to swap out -- This is an error
			return errors.New("no pages to swap out")
		}
	}
	for _, page := range swapOutList {
		err := mmu.SwapOutPhysicalPage(page)
		if err != nil {
			RemoteLogging.LogEvent("ERROR", "SwapOutOldPages", "Error swapping out old page")
			return err
		}
	}
	RemoteLogging.LogEvent("INFO", "SwapOutOldPages", "Old pages swapped out")
	return nil
}

func (mmu *MMUStruct) TryPageSwap(page int) error {
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
	return errors.New("page swap failed")
}

func (mmu *MMUStruct) WriteVirtualPage(
	owner int, group int,
	mode int, seg int, page int,
	buffer []byte) error {
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		return errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		return errors.New("page is not active")
	}
	// Check our permissions
	vpage := mmu.VirtualMemory[page]
	mask := 0
	if owner == vpage.ProcessID {
		mask = PageProtectionMaskUser
	}
	if owner != vpage.ProcessID && group != vpage.GroupID {
		mask = PageProtectionMaskWorld
	}
	if !mmu.CheckPermissionsOk(mode, mask, vpage.Protection) {
		return errors.New("permission denied")
	}
	if mmu.VirtualMemory[page].SegmentID != seg {
		return errors.New("invalid segment")
	}
	// If page isn't in memory, bring it in
	if mmu.PageIsOnDisk(page) {
		err := mmu.TryPageSwap(page)
		if err != nil {
			return err
		}
	}
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID

	copy(mmu.PhysicalMem[physicalPage*PageSize:physicalPage*PageSize+PageSize], buffer)
	mmu.SetPageIsDirty(page)
	mmu.pruneLRUCache(page)
	return nil
}

func (mmu *MMUStruct) ReadVirtualPage(
	owner int, group int,
	mode int, seg int, page int) ([]byte, error) {
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		return nil, errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		return nil, errors.New("page is not active")
	}
	// Check our permissions
	vpage := mmu.VirtualMemory[page]
	mask := 0
	if owner == vpage.ProcessID {
		mask = PageProtectionMaskUser
	}
	if owner != vpage.ProcessID && group != vpage.GroupID {
		mask = PageProtectionMaskWorld
	}
	if !mmu.CheckPermissionsOk(mode, mask, vpage.Protection) {
		return nil, errors.New("permission denied")
	}
	if mmu.VirtualMemory[page].SegmentID != seg {
		return nil, errors.New("invalid segment")
	}
	// If page isn't in memory, bring it in
	if mmu.PageIsOnDisk(page) {
		err := mmu.TryPageSwap(page)
		if err != nil {
			return nil, err
		}
	}
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID
	mmu.SetPageIsDirty(page)
	mmu.pruneLRUCache(page)
	return mmu.PhysicalMem[physicalPage*PageSize : physicalPage*PageSize+PageSize], nil
}

func (mmu *MMUStruct) FreeBulkPages(pages []int) error {
	for _, page := range pages {
		err := mmu.FreeVirtualPage(page)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mmu *MMUStruct) AllocateBulkPages(uid int, gid int, prot int, seg int, desiredPages int) ([]int, error) {
	lst := make([]int, desiredPages)
	for i := 0; i < desiredPages; i++ {
		page, _, err := mmu.AllocateNewVirtualPageNoSwap(uid, gid, prot, seg)
		if err != nil {
			return nil, err
		}
		lst = append(lst, page)
	}
	return lst, nil
}
