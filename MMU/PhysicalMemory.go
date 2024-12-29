package MMU

import "errors"

// MakePhysicalPageTable -- Create the table for physical pages
// Used when we need to create the physical page table the first time.
func (mmu *MMUStruct) MakePhysicalPageTable() error {
	// For every possible page, put it on the free list and make a table entry
	for i := 0; i < mmu.MMUConfig.NumPhysicalPages; i++ {
		mmu.FreePhysicalPages = append(mmu.FreePhysicalPages, i)
	}
	mmu.UsedPhysicalPages = make([]int, mmu.MMUConfig.NumPhysicalPages)
	mmu.PhysicalMem = make([]byte, mmu.MMUConfig.NumPhysicalPages*PageSize)
	return nil
}

// PercentFreePages -- How much of the physical page table is free
func (mmu *MMUStruct) PercentFreePages() (float64, error) {
	return float64(len(mmu.FreePhysicalPages)) / float64(mmu.MMUConfig.NumPhysicalPages), nil
}

// AllocateNewPhysicalPage -- allocate a new physical page from the free pages list
// Returns the page number or error
func (mmu *MMUStruct) AllocateNewPhysicalPage() (int, error) {
	if len(mmu.FreePhysicalPages) == 0 {
		return 0, errors.New("no physical pages")
	}
	pageID := mmu.FreePhysicalPages[0]
	mmu.FreePhysicalPages = mmu.FreePhysicalPages[1:]
	mmu.UsedPhysicalPages = append(mmu.UsedPhysicalPages, pageID)
	return pageID, nil
}

// ReturnPhysicalPage -- Return a physical page to the free pool
// Takes a page number and returns an error
func (mmu *MMUStruct) ReturnPhysicalPage(page int) error {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return errors.New("invalid physical page")
	}
	for ix, _ := range mmu.PhysicalMem {
		if ix == page {
			mmu.FreePhysicalPages = append(mmu.FreePhysicalPages, page)
			mmu.UsedPhysicalPages = mmu.UsedPhysicalPages[:ix]
			mmu.UsedPhysicalPages = mmu.UsedPhysicalPages[1:]
			return nil
		}
	}
	return errors.New("invalid physical page")
}

// ReadPhysicalPage -- Reads a physical page to a buffer
// Takes a page and returns a buffer and an error
func (mmu *MMUStruct) ReadPhysicalPage(page int) ([]byte, error) {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return nil, errors.New("invalid physical page")
	}
	return mmu.PhysicalMem[page*PageSize : (page+1)*PageSize], nil
}

// WritePhysicalPage - Write a buffer to a physical page
// Takes a page number, a byte buffer and returns an error
func (mmu *MMUStruct) WritePhysicalPage(page int, buffer []byte) error {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return errors.New("invalid physical page")
	}
	copy(mmu.PhysicalMem[page*PageSize:(page+1)*PageSize], buffer)
	return nil
}
