package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
	"strconv"
)

// PhysicalPage
// For every physical page we manage, we keep one of these structures
type PhysicalPage struct {
	Buffer     []byte
	MemoryType int
	IsInUse    bool
}

// PhysicalMemoryContainer
// For all of our physical memory regions, we keep the data in this container
type PhysicalMemoryContainer struct {
	Regions          []PhysicalMemoryRegion  // Our memory map
	MemoryPages      map[uint32]PhysicalPage // THe pages themselves
	FreeVirtualPages *list.List              // Vitual page lists
	UsedVirtualPages *list.List
}

// PhysicalMemory_Initialize
// initializes a physical memory container for a specified name.
// It retrieves the memory map, validates it, and populates physical pages based on their memory types.
// Returns a pointer to a PhysicalMemoryContainer and an error if initialization fails.
func PhysicalMemory_Initialize(name string) (*PhysicalMemoryContainer, error) {
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Initialize",
		"Setting yp physical memory")
	// See if the memory map exists for this name
	pr, ok := MemoryMapTable[name]
	if !ok {
		RemoteLogging.LogEvent("ERROR",
			"PhysicalMemory_Initialize",
			"Physical Memory map name not found")
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
				RemoteLogging.LogEvent("ERROR",
					"PhysicalMemory_Initialize",
					"Invalid memory type")
				return nil, errors.New("Invalid memory type")
			}
			currPage++
		}
	}
	// Done - return the container
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Initialize",
		"Physical memory initialized")
	return &pmc, nil
}

// Terminate
// Releases all memory pages in the container by clearing their buffers and removing them from the map.
// It also logs the termination event using the RemoteLogging system.
func (pmc *PhysicalMemoryContainer) Terminate() {
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Terminate",
		"Terminating physical memory")
	for page, val := range pmc.MemoryPages {
		val.Buffer = nil
		delete(pmc.MemoryPages, page)
	}
}

// ReturnListOfPageType
// Given a page type, iterates over memory pages and returns a list of pages matching the specified memory type.
func (pmc *PhysicalMemoryContainer) ReturnListOfPageType(ptype int) *list.List {
	RemoteLogging.LogEvent("INFO",
		"Physical_ReturnListOfPages",
		"Return a page list of type "+strconv.Itoa(ptype)+"")
	l := list.New()
	for page, val := range pmc.MemoryPages {
		if val.MemoryType == ptype {
			l.PushBack(page)
		}
	}
	RemoteLogging.LogEvent(
		"INFO",
		"Physical_ReturnListOfPages",
		"Returning list of "+strconv.Itoa(l.Len())+" pages of type "+strconv.Itoa(ptype)+"")
	return l
}

// ReturnTotalNumberOfPages
// Returns the total number of memory pages available in the PhysicalMemoryContainer.
func (pmc *PhysicalMemoryContainer) ReturnTotalNumberOfPages() uint32 {
	num := uint32(len(pmc.MemoryPages))
	RemoteLogging.LogEvent("INFO", "Physical_ReturnTotalNumberOfPages",
		"Returning total "+strconv.Itoa(int(num))+" pages")
	return num
}

// AllocateVirtualPage
// Reserves a free virtual memory page and marks it as in use, returning the page number or an error.
func (pmc *PhysicalMemoryContainer) AllocateVirtualPage() (uint32, error) {
	RemoteLogging.LogEvent("INFO", "Physical_AllocateVirtualPage", "Allocating virtual page")
	if pmc.FreeVirtualPages.Len() == 0 {
		RemoteLogging.LogEvent("ERROR",
			"Physical_AllocateVirtualPage", "No free pages")
		return 0, errors.New("No free virtual pages")
	}
	page := pmc.FreeVirtualPages.Remove(pmc.FreeVirtualPages.Front()).(uint32)
	val := pmc.MemoryPages[page]
	val.IsInUse = true
	pmc.MemoryPages[page] = val
	pmc.UsedVirtualPages.PushBack(page)
	RemoteLogging.LogEvent("INFO",
		"Physical_AllocateVirtualPage",
		"Allocated virtual page "+strconv.Itoa(int(page))+"")
	return page, nil
}

// ReturnVirtualPage
// Releases a virtual memory page back to the free pool and validates the page type and usage state.
// Returns an error if the page is not found, not in use, or not of the correct type (MemoryTypeVirtualRAM).
func (pmc *PhysicalMemoryContainer) ReturnVirtualPage(page uint32) error {
	RemoteLogging.LogEvent("INFO",
		"Physical_ReturnVirtualPage",
		"Returning virtual page "+strconv.Itoa(int(page)))
	val, ok := pmc.MemoryPages[page]
	if !ok {
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReturnVirtualPage",
			"Page not found")
		return errors.New("Page not found")
	}
	if val.IsInUse == false {

		RemoteLogging.LogEvent("ERROR",
			"Physical_ReturnVirtualPage",
			"Page is not in use")
		return errors.New("Page is not in use")
	}
	if pmc.MemoryPages[page].MemoryType == MemoryTypeVirtualRAM {
		pmc.FreeVirtualPages.PushBack(page)
		pmc.UsedVirtualPages.Remove(pmc.UsedVirtualPages.Front())
		val.IsInUse = false
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO", "Physical_ReturnVirtualPage",
			"Returned virtual page "+strconv.Itoa(int(page))+"")
		return nil
	}
	RemoteLogging.LogEvent("ERROR",
		"Physical_ReturnVirtualPage",
		"Page is not of type "+strconv.Itoa(int(MemoryTypeVirtualRAM))+"")
	return errors.New("Page is not of type " + strconv.Itoa(int(MemoryTypeVirtualRAM)))
}

func (pmc *PhysicalMemoryContainer) VirtualPercentFree() int {
	RemoteLogging.LogEvent("INFO",
		"Physical_VirtualPercentFree",
		"Calculating virtual percent free")
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

// ReturnMemoryMap
// Returns the memory map consisting of all physical memory regions contained in the PhysicalMemoryContainer.
func (pmc *PhysicalMemoryContainer) ReturnMemoryMap() []PhysicalMemoryRegion {
	RemoteLogging.LogEvent("INFO", "Physical_ReturnMemoryMap", "Returning memory map")
	return pmc.Regions
}
