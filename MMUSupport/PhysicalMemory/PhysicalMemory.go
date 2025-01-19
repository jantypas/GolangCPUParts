package PhysicalMemory

import (
	"container/list"
	"errors"
)

const (
	MemoryTypeEmpty       = 0
	MemoryTypeVirtualRAM  = 1
	MemoryTypePhysicalRAM = 2
	MemoryTypePhysicalROM = 3
	MemoryTypeKernelRAM   = 4
	MemoryTypeKernelROM   = 5
	MemoryTypeIORAM       = 6
	MemoryTypeIOROM       = 7
	MemoryTypeBufferRAM   = 8
	PageSize              = 4096
)

// PhysicalPage
// For every physical page we manage, we keep one of these structures
type PhysicalPage struct {
	Buffer     []byte
	MemoryType int
	IsInUse    bool
}

type PhysicalMemoryRegion struct {
	Comment    string
	NumPages   uint32
	MemoryType int
}

type PhysicalMemoryContainer struct {
	Regions          []PhysicalMemoryRegion
	MemoryPages      map[uint32]PhysicalPage
	FreeVirtualPages *list.List
	UsedVirtualPages *list.List
}

func PhysicalMemory_Initialize(name string) (*PhysicalMemoryContainer, error) {
	// See if the memory map exists for this name
	pr, ok := MemoryMapTable[name]
	if !ok {
		return nil, errors.New("Physical Memory Region not found")
	}
	// Build teh base of the memory container
	pmc := PhysicalMemoryContainer{
		Regions: pr,
	}
	// For each valid memory page, put it in the page map along with its type
	var currPage uint32 = 0
	for reg := range pr {
		for i := uint32(0); i < pr[reg].NumPages; i++ {
			switch pr[currPage].MemoryType {
			case MemoryTypeVirtualRAM:
				pmc.MemoryPages[currPage] = PhysicalPage{
					Buffer:     make([]byte, PageSize),
					MemoryType: pr[currPage].MemoryType,
					IsInUse:    false,
				}
				pmc.FreeVirtualPages.PushBack(currPage)
				break
			case MemoryTypeKernelRAM:
			case MemoryTypePhysicalRAM:
			case MemoryTypeBufferRAM:
			case MemoryTypeIORAM:
				pmc.MemoryPages[currPage] = PhysicalPage{
					Buffer:     make([]byte, PageSize),
					MemoryType: pr[currPage].MemoryType,
				}
				break
			case MemoryTypeEmpty:
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: pr[currPage].MemoryType,
				}
				break
			case MemoryTypePhysicalROM:
			case MemoryTypeKernelROM:
			case MemoryTypeIOROM:
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: pr[currPage].MemoryType,
				}
				break
			default:
				return nil, errors.New("Invalid memory type")
			}
			currPage++
		}
	}
	// Doen return the container
	return &pmc, nil
}

func (pmc *PhysicalMemoryContainer) Terminate() {
	for page, val := range pmc.MemoryPages {
		val.Buffer = nil
		delete(pmc.MemoryPages, page)
	}
}

func (pmc *PhysicalMemoryContainer) ReturnListOfPageType(ptype int) *list.List {
	l := list.New()
	for page, val := range pmc.MemoryPages {
		if val.MemoryType == ptype {
			l.PushBack(page)
		}
	}
	return l
}

func (pmc *PhysicalMemoryContainer) ReturnTotalNumberOfPages() uint32 {
	return uint32(len(pmc.MemoryPages))
}

func (pmc *PhysicalMemoryContainer) AllocateVirtualPage() (uint32, error) {
	if pmc.FreeVirtualPages.Len() == 0 {
		return 0, errors.New("No free virtual pages")
	}
	page := uint32(pmc.FreeVirtualPages.Remove(pmc.FreeVirtualPages.Front()).(uint32))
	val := pmc.MemoryPages[page]
	val.IsInUse = true
	pmc.MemoryPages[page] = val
	pmc.UsedVirtualPages.PushBack(page)
	return page, nil
}

func (pmc *PhysicalMemoryContainer) ReturnVirtualPage(page uint32) error {
	val, ok := pmc.MemoryPages[page]
	if !ok {
		return errors.New("Page not found")
	}
	if val.IsInUse == false {
		return errors.New("Page is not in use")
	}
	if pmc.MemoryPages[page].MemoryType == MemoryTypeVirtualRAM {
		pmc.FreeVirtualPages.PushBack(page)
		pmc.UsedVirtualPages.Remove(pmc.UsedVirtualPages.Front())
		val.IsInUse = false
		pmc.MemoryPages[page] = val
		return nil
	}
	return errors.New("Page is wrong type")
}

func (pmc *PhysicalMemoryContainer) VritualPercentFree() int {
	return int(float64(pmc.FreeVirtualPages.Len()) / float64(pmc.ReturnTotalNumberOfPages()) * 100)
}

func (pmc *PhysicalMemoryContainer) ReadAddress(addr uint64) (byte, error) {
	page := uint32(addr / PageSize)
	offset := addr % PageSize
	val, ok := pmc.MemoryPages[page]
	if !ok {
		return 0, errors.New("Page not found")
	}
	if val.IsInUse == false {
		return 0, errors.New("Page is not in use")
	}
	return val.Buffer[offset], nil
}

func (pmc *PhysicalMemoryContainer) WriteAddress(addr uint64, data byte) error {
	page := uint32(addr / PageSize)
	offset := addr % PageSize
	val, ok := pmc.MemoryPages[page]
	if !ok {
		return errors.New("Page not found")
	}
	if val.IsInUse == false {
		return errors.New("Page is not in use")
	}
	buffer := pmc.MemoryPages[page].Buffer
	buffer[offset] = data
	val.Buffer = buffer
	pmc.MemoryPages[page] = val
	return nil
}
