package MMU

import (
	"GolangCPUParts"
	"errors"
)

// Initialize -- Initialize the MMMU
func (mmu *MMUStruct) Initialize(conf *MMUConfig) error {
	mmu.MMUConfig = *conf
	mmu.TLB = make([]MMUTLB, mmu.MMUConfig.TLBSize)
	mmu.FreeVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	mmu.UsedVirtualPages = make([]int, mmu.MMUConfig.NumVirtualPages)
	err := mmu.MMUConfig.Swapper.Initialize()
	if err != nil {
		return err
	}
	return nil
}

// Terminate -- shutdown the MMU
func (mmu *MMUStruct) Terminate() error {
	err := mmu.MMUConfig.Swapper.Terminate()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// makePhysicalPage -- Create a blank physical page
func (mmu *MMUStruct) makePhysicalPage(page int) (PhysicalPage, error) {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return nil, errors.New("invalid physical page")
	}
	pp := PhysicalPage{
		PhysicalPage: page,
	}
	return pp, nil
}

// makePhysicalPageTable -- Create the table for physical pages
func (mmu *MMUStruct) makePhysicalPageTable() error {
	for i := 0; i < mmu.MMUConfig.NumPhysicalPages; i++ {
		mmu.FreePhysicalPages = append(mmu.FreePhysicalPages, i)
		mmu.PhysicalMem[i] = PhysicalPage{
			PhysicalPage: i,
		}
	}
	return nil
}

// allocateNewPhysicalPage -- Allocate a physical page in the apge table
func (mmu *MMUStruct) allocateNewPhysicalPage() (int, error) {
	if len(mmu.FreePhysicalPages) == 0 {
		return 0, errors.New("no physical pages")
	}
	pageID := mmu.FreePhysicalPages[0]
	mmu.FreePhysicalPages = mmu.FreePhysicalPages[1:]
	mmu.UsedPhysicalPages = append(mmu.UsedPhysicalPages, pageID)
	return pageID, nil
}

// freePhysicalPage -- free a physical page in the physical page table for re-use
func (mmu *MMUStruct) freePhysicalPage(page int) error {
	if page > mmu.MMUConfig.NumPhysicalPages {
		return errors.New("invalid physical page")
	}
	for ix, vx := range mmu.PhysicalMem {
		if ix == page {
			mmu.FreePhysicalPages = append(mmu.FreePhysicalPages, page)
			mmu.UsedPhysicalPages = mmu.UsedPhysicalPages[:ix]
			mmu.UsedPhysicalPages = mmu.UsedPhysicalPages[1:]
			return nil
		}
	}

}
