package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
	"container/list"
	"errors"
	"strconv"
)

const (
	MemoryTypeEmpty       = 0x0001
	MemoryTypeVirtualRAM  = 0x0002
	MemoryTypePhysicalRAM = 0x0004
	MemoryTypePhysicalROM = 0x0008
	MemoryTypeKernelRAM   = 0x0010
	MemoryTypeKernelROM   = 0x0020
	MemoryTypeIORAM       = 0x0040
	MemoryTypeIOROM       = 0x0080
	MemoryTypeBufferRAM   = 0x0100

	PageSize = 4096
)

var MemoryTypeNames = map[int]string{
	MemoryTypeEmpty:       "Empty",
	MemoryTypeVirtualRAM:  "Virtual RAM",
	MemoryTypePhysicalRAM: "Physical RAM",
	MemoryTypePhysicalROM: "Physical ROM",
	MemoryTypeKernelRAM:   "Kernel RAM",
	MemoryTypeKernelROM:   "Kernel ROM",
	MemoryTypeIORAM:       "I/O RAM",
	MemoryTypeIOROM:       "I/O ROM",
	MemoryTypeBufferRAM:   "Buffer RAM",
}

type PhysicalMemoryRegion struct {
	Comment    string
	NumPages   uint32
	MemoryType int
}

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
		"Setting up physical memory")
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
		Regions:     pr,
		MemoryPages: make(map[uint32]PhysicalPage),
	}
	// For each valid memory page, put it in the page map along with its type
	var currPage uint32 = 0
	var i uint32 = 0
	lpr := len(pr)
	for idx := 0; idx < lpr; idx++ {
		xp := pr[idx]
		switch xp.MemoryType {
		case MemoryTypeVirtualRAM:
			pmc.FreeVirtualPages = list.New()
			pmc.UsedVirtualPages = list.New()
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
					IsInUse:    false,
					Buffer:     make([]byte, PageSize),
				}
				pmc.FreeVirtualPages.PushBack(currPage)
				currPage++
			}
			continue
		case MemoryTypeKernelRAM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
					Buffer:     make([]byte, PageSize),
				}
				currPage++
			}
			continue
		case MemoryTypePhysicalRAM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
					Buffer:     make([]byte, PageSize),
				}
				currPage++
			}
			continue
		case MemoryTypeBufferRAM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
					Buffer:     make([]byte, PageSize),
				}
				currPage++
			}
			continue
		case MemoryTypeIORAM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
				}
				currPage++
			}
			continue
		case MemoryTypeEmpty:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
				}
				currPage++
			}
			continue
		case MemoryTypePhysicalROM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
				}
				currPage++
			}
			continue
		case MemoryTypeKernelROM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
				}
				currPage++
			}
			continue
		case MemoryTypeIOROM:
			for i = 0; i < xp.NumPages; i++ {
				pmc.MemoryPages[currPage] = PhysicalPage{
					MemoryType: xp.MemoryType,
				}
				currPage++
			}
			continue
		default:
		}
	}
	// Done - return the container
	RemoteLogging.LogEvent("INFO",
		"PhysicalMemory_Initialize",
		"Physical memory initialized with "+strconv.Itoa(int(currPage))+" pages")
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
		// Can't find the page
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReturnVirtualPage",
			"Page not found")
		return errors.New("Page not found")
	}
	if val.MemoryType != MemoryTypeVirtualRAM {
		// Not a virtual page
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReturnVirtualPage",
			"Page wrong type")
		return errors.New("Page wrong type")
	}
	if val.IsInUse == false {
		// Page is not in use-- we can't do this
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReturnVirtualPage",
			"Page is not in use")
		return errors.New("Page is not in use")
	}
	// Free the old page
	if pmc.MemoryPages[page].MemoryType == MemoryTypeVirtualRAM {
		pmc.FreeVirtualPages.PushBack(page)
		pmc.UsedVirtualPages.Remove(pmc.UsedVirtualPages.Front())
		val.IsInUse = false
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO", "Physical_ReturnVirtualPage",
			"Returned virtual page "+strconv.Itoa(int(page))+"")
		return nil
	}
	// Done
	RemoteLogging.LogEvent("ERROR",
		"Physical_ReturnVirtualPage",
		"Page is not of type "+strconv.Itoa(int(MemoryTypeVirtualRAM))+"")
	return errors.New("Page is not of type " + strconv.Itoa(int(MemoryTypeVirtualRAM)))
}

// VirtualPercentFree
// Calculates the percentage of free virtual memory pages in the memory container.
func (pmc *PhysicalMemoryContainer) VirtualPercentFree() int {
	RemoteLogging.LogEvent("INFO",
		"Physical_VirtualPercentFree",
		"Calculating virtual percent free")
	return int(float64(pmc.FreeVirtualPages.Len()) / float64(pmc.ReturnTotalNumberOfPages()) * 100)
}

// ReadAddress
// Retrieves the byte stored at the given virtual address and returns an error if access fails.
func (pmc *PhysicalMemoryContainer) ReadAddress(addr uint64) (byte, error) {
	page := uint32(addr / PageSize)
	offset := addr % PageSize
	val, ok := pmc.MemoryPages[page]
	if !ok {
		// Page not found
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReadAddress",
			"Page not found")
		return 0, errors.New("Page not found")
	}
	if val.IsInUse == false {
		RemoteLogging.LogEvent("ERROR", "Physical_ReadAddress",
			"Page is not in use")
		return 0, errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
		if !val.IsInUse {
			RemoteLogging.LogEvent("ERROR",
				"Physical_ReadAddress",
				"Page is not in use")
			return 0, errors.New("Page is not in use")
		} else {
			RemoteLogging.LogEvent("INFO",
				"Physical_ReadAddress",
				"Reading address "+strconv.Itoa(int(addr))+"")
			return val.Buffer[offset], nil
		}
	case MemoryTypePhysicalRAM:
		RemoteLogging.LogEvent("INFO",
			"Physical_ReadAddress",
			"Reading address "+strconv.Itoa(int(addr))+"")
		return val.Buffer[offset], nil
	case MemoryTypeBufferRAM:
		RemoteLogging.LogEvent("INFO", "Physical_ReadAddress",
			"Reading address "+strconv.Itoa(int(addr))+"")
		return val.Buffer[offset], nil
	case MemoryTypePhysicalROM:
		RemoteLogging.LogEvent("INFO", "Physical_ReadAddress",
			"Reading address "+strconv.Itoa(int(addr))+"")
		return val.Buffer[offset], nil
	case MemoryTypeKernelRAM:
		RemoteLogging.LogEvent("INFO", "Physical_ReadAddress",
			"Reading address "+strconv.Itoa(int(addr))+"")
		return val.Buffer[offset], nil
	case MemoryTypeKernelROM:
		RemoteLogging.LogEvent("INFO", "Physical_ReadAddress",
			"Read complete")
		return val.Buffer[offset], nil
	case MemoryTypeIORAM:
		RemoteLogging.LogEvent("INFO",
			"Physical_ReadAddress", "I/O not implemented")
		return 0, errors.New("I/O not implemented")
	case MemoryTypeIOROM:
		RemoteLogging.LogEvent("ERROR",
			"Physical_ReadAddress",
			"I/O not implemented")
		return 0, errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent(
			"ERROR",
			"Physical_ReadAddress",
			"Page is empty")
		return 0, errors.New("Page is empty")
	default:
		RemoteLogging.LogEvent(
			"ERROR",
			"Physical_ReadAddress",
			"Page is wrong type")
		return 0, errors.New("Page is wrong type")
	}
	RemoteLogging.LogEvent("ERROR",
		"Physical_ReadAddress", "Page is wrong type")
	return 0, errors.New("Page is wrong type")
}

// WriteAddress
// Writes a byte of data to the specified memory address within the physical memory container.
// It identifies the corresponding memory page, validates its state, and ensures it is writable.
// Returns an error if the page is not found, not in use, read-only, or of the wrong type.
func (pmc *PhysicalMemoryContainer) WriteAddress(addr uint64, data byte) error {
	RemoteLogging.LogEvent("INFO",
		"Physical_WriteAddress",
		"Writing address "+strconv.Itoa(int(addr))+" to "+strconv.Itoa(int(data))+"")
	page := uint32(addr / PageSize)
	offset := addr % PageSize
	val, ok := pmc.MemoryPages[page]
	if !ok {
		RemoteLogging.LogEvent("ERROR",
			"Physical_WriteAddress",
			"Page not found")
		return errors.New("Page not found")
	}
	if val.MemoryType == MemoryTypeVirtualRAM && val.IsInUse == false {
		RemoteLogging.LogEvent(
			"ERROR",
			"Physical_WriteAddress",
			"Page is not in use")
		return errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
		val := pmc.MemoryPages[page]
		val.Buffer[offset] = data
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO",
			"Physical_WriteAddress",
			"Write address completed")
		return nil
	case MemoryTypePhysicalRAM:
		val := pmc.MemoryPages[page]
		val.Buffer[offset] = data
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO",
			"Physical_WriteAddress",
			"Write address completed")
		return nil
	case MemoryTypeBufferRAM:
		val := pmc.MemoryPages[page]
		val.Buffer[offset] = data
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO",
			"Physical_WriteAddress",
			"Write address completed")
		return nil
	case MemoryTypeKernelRAM:
		val := pmc.MemoryPages[page]
		val.Buffer[offset] = data
		pmc.MemoryPages[page] = val
		RemoteLogging.LogEvent("INFO",
			"Physical_WriteAddress",
			"Write address completed")
		return nil
	case MemoryTypePhysicalROM:
		RemoteLogging.LogEvent("ERROR",
			"Physical_WriteAddress", "Page is read only")
		return errors.New("Page is read only")
	case MemoryTypeKernelROM:
		RemoteLogging.LogEvent("ERROR",
			"Physical_WriteAddress", "Page is read only")
		return errors.New("Page is read only")
	case MemoryTypeIORAM:
		RemoteLogging.LogEvent("ERROR",
			"Physical_WriteAddress", "I/O not implemented")
		return errors.New("I/O not implemented")
	case MemoryTypeIOROM:
		RemoteLogging.LogEvent("ERROR",
			"Physical_WriteAddress", "I/O not implemented")
		return errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent(
			"ERROR", "Physical_WriteAddress", "Page is empty")
		return errors.New("Page is empty")
	}
	RemoteLogging.LogEvent("ERROR", "Physical_WriteAddress", "Page is wrong type")
	return errors.New("Page is wrong type")
}

// LoadPage
// Loads a memory page specified by its page number into the provided buffer if it meets the required conditions.
// Returns an error if the page does not exist, is not in use, is of the wrong type, or is read-only.
func (pmc *PhysicalMemoryContainer) LoadPage(page uint32, buffer []byte) error {
	val, ok := pmc.MemoryPages[uint32(page)]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "Page not found")
		return errors.New("Page not found")
	}
	if val.IsInUse == false {
		RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "Page is not in use")
		return errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypeKernelRAM:
		val := pmc.MemoryPages[page]
		val.Buffer = buffer
		copy(pmc.MemoryPages[page].Buffer, buffer)
		RemoteLogging.LogEvent("INFO", "Physical_LoadPage", "Page loaded")
		return nil
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelROM:
		RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "Page is read only")
		return errors.New("Page is read only")
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "I/O not implemented")
		return errors.New("I/O not implemented")
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "Page is empty")
		return errors.New("Page is empty")
	}
	RemoteLogging.LogEvent("ERROR", "Physical_LoadPage", "Page is wrong type")
	return errors.New("Page is wrong type")
}

// SavePage
// Retrieves the buffer associated with a memory page if it is in use and of a valid type, or returns an error.
func (pmc *PhysicalMemoryContainer) SavePage(page uint32) ([]byte, error) {
	val, ok := pmc.MemoryPages[page]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page not found")
		return nil, errors.New("Page not found")
	}
	if val.IsInUse == false {
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page is not in use")
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
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page not found")
		return nil, errors.New("Page not found")
	}
	if val.IsInUse == false {
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page is not in use")
		return nil, errors.New("Page is not in use")
	}
	switch val.MemoryType {
	case MemoryTypeVirtualRAM:
	case MemoryTypePhysicalRAM:
	case MemoryTypeBufferRAM:
	case MemoryTypePhysicalROM:
	case MemoryTypeKernelRAM:
	case MemoryTypeKernelROM:

		RemoteLogging.LogEvent("INFO", "Physical_SavePage", "Page saved")
		return val.Buffer, nil
	case MemoryTypeIORAM:
	case MemoryTypeIOROM:
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page saved")
		return pmc.MemoryPages[page].Buffer, nil
	case MemoryTypeEmpty:
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page is empty")
		return nil, errors.New("Page is empty")
	default:
		RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page is wrong type")
		return nil, errors.New("Page is wrong type")
	}
	RemoteLogging.LogEvent("ERROR", "Physical_SavePage", "Page is wrong type")
	return nil, errors.New("Page is wrong type")
}

// NumberOfFreeVirtualPages
// Returns the number of free virtual memory pages in the PhysicalMemoryContainer.
func (pmc *PhysicalMemoryContainer) NumberOfFreeVirtualPages() uint32 {
	return uint32(pmc.FreeVirtualPages.Len())
}

// NumberOfUsedVirtualPages
// Returns the number of virtual memory pages currently marked as in use in the container.
func (pmc *PhysicalMemoryContainer) NumberOfUsedVirtualPages() uint32 {
	return uint32(pmc.UsedVirtualPages.Len())
}

// AllocateNVirtualPage
// Allocates a specified number of virtual memory pages if available, returning a list of allocated pages or an error.
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

// ReturnNVirtualPage
// Releases a list of virtual memory pages back to the free pool and returns any encountered error.
func (pmc *PhysicalMemoryContainer) ReturnNVirtualPage(l *list.List) error {
	RemoteLogging.LogEvent("INFO",
		"Physical_ReturnNVirtualPage",
		"Returning "+strconv.Itoa(l.Len())+" pages")
	for e := l.Front(); e != nil; e = e.Next() {
		page := e.Value.(uint32)
		err := pmc.ReturnVirtualPage(page)
		if err != nil {
			return err
		}
	}
	RemoteLogging.LogEvent("INFO", "ReturnNVirtulPages",
		"Returned "+strconv.Itoa(l.Len())+" pages")
	return nil
}

// ReturnMemoryMap
// Returns the memory map consisting of all physical memory regions contained in the PhysicalMemoryContainer.
func (pmc *PhysicalMemoryContainer) ReturnMemoryMap() []PhysicalMemoryRegion {
	RemoteLogging.LogEvent("INFO", "Physical_ReturnMemoryMap", "Returning memory map")
	return pmc.Regions
}

// GetMemoryType
// Retrieves the memory type of a specified memory page within the PhysicalMemoryContainer.
func (pmc *PhysicalMemoryContainer) GetMemoryType(page uint32) int {
	return pmc.MemoryPages[page].MemoryType
}
