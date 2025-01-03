package MMUSupport

import "errors"

// AllocateNewPhysicalPage -- allocate a new physical page from the free pages list
// Returns the page number or error
func (mmu *MMUStruct) AllocateNewPhysicalPage() (int, error) {
	if mmu.FreePhysicalPages.Len() == 0 {
		return 0, errors.New("no physical pages")
	}
	pageID := mmu.FreePhysicalPages.Front().Value.(int)
	mmu.FreePhysicalPages.Remove(mmu.FreePhysicalPages.Front())
	mmu.UsedVirtualPages.PushBack(pageID)
	return pageID, nil
}

// ReturnPhysicalPage -- Return a physical page to the free pool
// Takes a page number and returns an error
func (mmu *MMUStruct) ReturnPhysicalPage(page int) error {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return errors.New("invalid physical page")
	}
	mmu.FreePhysicalPages.PushBack(page)
	mmu.UsedVirtualPages.Remove(mmu.UsedVirtualPages.Front())
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
