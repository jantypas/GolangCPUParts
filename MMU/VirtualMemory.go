package MMU

import "errors"

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

// MakeVirtualMemroyTable - Create the virtual emmroy table
// Build the virtual memory table strcutures.  Returns an error
func (mmu *MMUStruct) MakeVirtualMemoryTable() error {
	mmu.VirtualMemory = make([]VirtualPage, mmu.MMUConfig.NumVirtualPages)
	mmu.FreeVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	mmu.UsedVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	for i := 0; i < mmu.MMUConfig.NumVirtualPages; i++ {
		mmu.FreeVirtualPages = append(mmu.FreeVirtualPages, i)
	}
	mmu.LRUCache = make([]int, mmu.MMUConfig.NumVirtualPages)
	return nil
}

// AllocateNewVirtualPageNoSwap -- Allocate a new virtual page from the virtual page table
// Return the page table and an error
func (mmu *MMUStruct) AllocateNewVirtualPageNoSwap(
	owner int,
	group int,
	protect int) (int, int, error) {
	if len(mmu.FreeVirtualPages) == 0 {
		return 0, VirtualErrorNoPages, errors.New("no virtual pages")
	}
	pageID := mmu.FreeVirtualPages[0]
	mmu.FreeVirtualPages = mmu.FreeVirtualPages[1:]
	mmu.UsedVirtualPages = append(mmu.UsedVirtualPages, pageID)
	mmu.VirtualMemory[pageID].Protection = protect
	mmu.VirtualMemory[pageID].ProcessID = owner
	mmu.VirtualMemory[pageID].GroupID = group
	mmu.SetPageIsActive(pageID)
	mmu.SetPageIsOnDisk(pageID)
	mmu.VirtualMemory[pageID].PhysicalPageID = -1
	return pageID, 0, nil
}

// FreeVirtualPage -- Free a virtual and its associated physical page
// Takes a page ID and returns an error
func (mmu *MMUStruct) FreeVirtualPage(page int) error {
	// Make sure page is a valid page
	if page > mmu.MMUConfig.NumVirtualPages {
		return errors.New("invalid virtual page")
	}
	// If the page is not already active
	if !mmu.PageIsActive(page) {
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
			return err
		}
		mmu.ClearPageIsActive(page)
		return nil
	}
}

// SwapOutPhysicalPage -- Swap a physical page out to disk
// Takes a virtual page ID and returns an error
func (mmu *MMUStruct) SwapOutPhysicalPage(page int) error {
	// Make sure page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		return errors.New("invalid virtual page")
	}
	// Make sure page is active
	if !mmu.PageIsActive(page) {
		return errors.New("page is already free")
	}
	// Make sure page isn't swapped already
	if mmu.PageIsOnDisk(page) {
		return errors.New("page is already on disk")
	}
	// Find the physical page
	physicalPage := mmu.VirtualMemory[page].PhysicalPageID
	// Swap it out
	err := mmu.MMUConfig.Swapper.SwapOut(physicalPage, mmu.PhysicalMem[physicalPage*PageSize:])
	if err != nil {
		return err
	}
	// Make the page as swapped out
	mmu.SetPageIsOnDisk(page)
	// Return the physical page
	err = mmu.ReturnPhysicalPage(physicalPage)
	if err != nil {
		return err
	}
	mmu.VirtualMemory[page].PhysicalPageID = -1
	return nil
}

// SwapInPhysicalPage -- Swap a physical page in from disk
// Takes a virtual page ID and returns an error
func (mmu *MMUStruct) SwapInPhysicalPage(page int) error {
	// Make sure the page is valid
	if page > mmu.MMUConfig.NumVirtualPages {
		return errors.New("invalid virtual page")
	}
	// Make sure the page is already active
	if !mmu.PageIsActive(page) {
		return errors.New("page is already free")
	}
	// Make sure the page is on disk to swap in
	if !mmu.PageIsOnDisk(page) {
		return errors.New("page is not on disk")
	}
	// Try to get a physical page to swap into
	page, err := mmu.AllocateNewPhysicalPage()
	if err != nil {
		return err
	}
	// Set up the physical page
	mmu.VirtualMemory[page].PhysicalPageID = page
	mmu.ClearPageIsOnDisk(page)
	err := mmu.MMUConfig.Swapper.SwapIn(page, mmu.PhysicalMem[page*PageSize:])
	if err != nil {
		return err
	}
	return nil
}

func (mmu *MMUStruct) SwapOutOldPages() error {
	for ix, vx :=
}