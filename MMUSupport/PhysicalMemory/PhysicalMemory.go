package PhysicalMemory

import (
	"container/list"
	"errors"
)

// PhysicalPage
// For every physical page we manage, we keep one of these structures
type PhysicalPage struct {
	Buffer     []byte
	MemoryType int
	IsInUse    bool
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
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelRAM:
	case MemoryTypeKernelROM:
		return val.Buffer[offset], nil
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		return 0, errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		return 0, errors.New("Page is empty")
	default:
		return 0, errors.New("Page is wrong type")
	}
	return 0, errors.New("Page is wrong type")
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
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypeKernelRAM:
		val := pmc.MemoryPages[page]
		val.Buffer[offset] = data
		pmc.MemoryPages[page] = val
		return nil
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelROM:
		return errors.New("Page is read only")
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		return errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		return errors.New("Page is empty")
	}
	return errors.New("Page is wrong type")
}

func (pmc *PhysicalMemoryContainer) LoadPage(page uint32, buffer []byte) error {
	val, ok := pmc.MemoryPages[uint32(page)]
	if !ok {
		return errors.New("Page not found")
	}
	if val.IsInUse == false {
		return errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypeKernelRAM:
		val := pmc.MemoryPages[page]
		val.Buffer = buffer
		pmc.MemoryPages[page] = val
		return nil
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelROM:
		return errors.New("Page is read only")
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		return errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		return errors.New("Page is empty")
	}
	return errors.New("Page is wrong type")
}

func (pmc *PhysicalMemoryContainer) SavePage(page uint32) ([]byte, error) {
	val, ok := pmc.MemoryPages[page]
	if !ok {
		return nil, errors.New("Page not found")
	}
	if val.IsInUse == false {
		return nil, errors.New("Page is not in use")
	}
	val, ok = pmc.MemoryPages[page]
	if !ok {
		return nil, errors.New("Page not found")
	}
	if val.IsInUse == false {
		return nil, errors.New("Page is not in use")
	}
	val, ok = pmc.MemoryPages[page]
	if !ok {
		return nil, errors.New("Page not found")
	}
	if val.IsInUse == false {
		return nil, errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelRAM:
	case MemoryTypeKernelROM:
		return val.Buffer, nil
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		return nil, errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		return nil, errors.New("Page is empty")
	default:
		return nil, errors.New("Page is wrong type")
	}
	return nil, errors.New("Page is wrong type")
}

func (pmc *PhysicalMemoryContainer) NumberOfFreeVirtualPages() uint32 {
	return uint32(pmc.FreeVirtualPages.Len())
}

func (pmc *PhysicalMemoryContainer) NumberOfUsedVirtualPages() uint32 {
	return uint32(pmc.UsedVirtualPages.Len())
}

func (pmc *PhysicalMemoryContainer) AllocateNVirtualPage(num uint32) (*list.List, error) {
	l := list.New()
	if pmc.NumberOfFreeVirtualPages() < num {
		return nil, errors.New("Not enough free pages")
	}
	for i := uint32(0); i < num; i++ {
		page, err := pmc.AllocateVirtualPage()
		if err != nil {
			return nil, err
		} else {
			l.PushBack(page)
		}
	}
	return l, nil
}

func (pmc *PhysicalMemoryContainer) ReturnNVirtualPage(l *list.List) error {
	for e := l.Front(); e != nil; e = e.Next() {
		page := e.Value.(uint32)
		err := pmc.ReturnVirtualPage(page)
		if err != nil {
			return err
		}
	}
	return nil
}
