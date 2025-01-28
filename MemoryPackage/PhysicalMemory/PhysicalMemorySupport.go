package PhysicalMemory

import (
	"GolangCPUParts/MemoryPackage/MemoryMap"
	"GolangCPUParts/RemoteLogging"
	"errors"
	"strconv"
)

const PageSize = 4096

// PhysicalBlock
// The physical memory block contains a block of bytes for a physical memory region
type PhysicalBlock struct {
	Buffer      []byte // The buffer of bytes that contain data
	NumPages    uint32 // Number of pages in the buffer
	StartPage   uint32 // Where does our material start (on a page)
	Protections uint64 // Any protection rules
	Key         uint16
}

// PhysicalMemoryContainer
// THe PhysicalMemoryContainer contains all PhysicalMemoryBlocks
type PhysicalMemoryContainer struct {
	MyMap          []MemoryMap.MemoryMapRegion
	PhysicalBlocks []PhysicalBlock
}

// PhysicalMemoryInitialize
// Given a memory map, build the physical memory blocks as a container and returns it
func PhysicalMemoryInitialize(mmap []MemoryMap.MemoryMapRegion) (*PhysicalMemoryContainer, error) {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryInitialize", "Initializing physical memory")
	pmc := &PhysicalMemoryContainer{}
	pmc.PhysicalBlocks = make([]PhysicalBlock, len(mmap))
	totalBytes := uint64(0)
	totalPages := 0
	// For each item in the map, build an object
	for i := range mmap {
		pmc.PhysicalBlocks[i].Buffer = make([]byte, mmap[i].EndAddress-mmap[i].StartAddress)
		pmc.PhysicalBlocks[i].NumPages = uint32(mmap[i].EndAddress-mmap[i].StartAddress) / 4096
		pmc.PhysicalBlocks[i].Protections = mmap[i].Permissions
		pmc.PhysicalBlocks[i].StartPage = uint32(mmap[i].StartAddress) / PageSize
		pmc.PhysicalBlocks[i].Key = mmap[i].Key
		totalBytes += mmap[i].EndAddress - mmap[i].StartAddress
		totalPages += int(mmap[i].EndAddress-mmap[i].StartAddress) / PageSize
	}
	msg := "Initialized physical memory with " + strconv.Itoa(int(totalBytes)) +
		" bytes and " + strconv.Itoa(totalPages) + " pages"
	pmc.MyMap = mmap
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryInitialize", msg)
	return pmc, nil
}

// Terminate
// Given a memory container, terminate everything
func (pmc *PhysicalMemoryContainer) Terminate() {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryTerminate", "Terminating physical memory")
	for i := range pmc.PhysicalBlocks {
		pmc.PhysicalBlocks[i].Buffer = nil
	}
}

// GetRegionByKey
// Retrieves a PhysicalBlock from the PhysicalBlocks slice that corresponds to the given key in the MyMap
// slice.
// It returns a pointer to the matched PhysicalBlock, or nil if no match is found.
func (pmc *PhysicalMemoryContainer) GetRegionByKey(key uint16) *PhysicalBlock {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "Getting block for key "+
		strconv.Itoa(int(key)))
	// Walk the map looking for a key
	for i := range pmc.MyMap {
		if pmc.MyMap[i].Key == key {
			// Return it
			return &pmc.PhysicalBlocks[i]
		}
	}
	// No key found
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "No block found for key "+
		strconv.Itoa(int(key)))
	return nil
}

// GetRegionByAddress
// Retrieves a PhysicalBlock corresponding to the provided address in the memory map,
// or nil if not found.
func (pmc *PhysicalMemoryContainer) GetRegionByAddress(addr uint64) *PhysicalBlock {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "Getting block for address "+
		strconv.Itoa(int(addr)))
	bl := MemoryMap.FindSegment(pmc.MyMap, addr)
	if bl != nil {
		return pmc.GetRegionByKey(bl.Key)
	}
	return nil
}

// ReadAddress retrieves a byte from the physical memory at the specified address.
// Returns the byte and an error if the address is invalid or not mapped.
func (pmc *PhysicalMemoryContainer) ReadPhysicalAddress(addr uint64) (byte, error) {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryReadAddress", "Reading address "+
		strconv.Itoa(int(addr)))
	// Compute the block by address
	bl := pmc.GetRegionByAddress(addr)
	if bl == nil {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryReadAddress", "Invalid address "+
			strconv.Itoa(int(addr)))
		return 0, nil
	}
	// Compute the address in the buffer
	page := uint32(addr / PageSize)
	if page > bl.NumPages {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryReadAddress", "Invalid address "+
			strconv.Itoa(int(addr)))
	}
	newAddr := (page - bl.StartPage) * PageSize
	return bl.Buffer[newAddr], nil
}

// WriteAddress writes a byte of data to the specified physical memory address.
// Returns an error if the address is invalid or not writable.
func (pmc *PhysicalMemoryContainer) WritePhysicalAddress(addr uint64, data byte) error {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryWriteAddress", "Writing address "+
		strconv.Itoa(int(addr)))
	// Compute the block by address
	bl := pmc.GetRegionByAddress(addr)
	if bl == nil {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryWriteAddress", "Invalid address "+
			strconv.Itoa(int(addr)))
		return nil
	}
	// Update buffer
	page := uint32(addr / PageSize)
	if page-bl.StartPage > bl.NumPages {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryWriteAddress", "Invalid address "+
			strconv.Itoa(int(addr)))
		return nil
	}
	newAddr := (page - bl.StartPage) * PageSize
	bl.Buffer[newAddr] = data
	return nil
}

// ReadPage retrieves the data of a physical memory page specified by its page number.
// Returns the page as a byte slice and an error if the page is invalid or not mapped.
func (pmc *PhysicalMemoryContainer) ReadPhysicalPage(page uint32) ([]byte, error) {
	bl := pmc.GetRegionByAddress(uint64(page) * PageSize)
	if bl == nil {
		return nil, errors.New("Invalid page")
	}
	newPage := page - bl.StartPage
	return bl.Buffer[newPage*PageSize : (newPage+1)*PageSize], nil
}

// WritePage
// Write an entire page of memory to the buffer
func (pmc *PhysicalMemoryContainer) WritePhysicalPage(page uint32, data []byte) error {
	if len(data) != PageSize {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryWritePage", "Invalid data length ")
		return errors.New(
			"WritePage: Data length must be equal to page size (" +
				strconv.Itoa(PageSize) + ")")
	}
	bl := pmc.GetRegionByAddress(uint64(page) * PageSize)
	if bl == nil {
		return errors.New("Invalid page")
	}
	newPage := page - bl.StartPage
	copy(bl.Buffer[newPage*PageSize:newPage*PageSize+PageSize], data)
	return nil
}

func (pmc *PhysicalMemoryContainer) GetRegionByTag(tag string) *PhysicalBlock {
	for i := range pmc.MyMap {
		if pmc.MyMap[i].Tag == tag {
			return &pmc.PhysicalBlocks[i]
		}
	}
	return nil
}

func (pmc *PhysicalMemoryContainer) ReadAddress(addr uint64) (byte, error) {
	bl := MemoryMap.FindSegment(pmc.MyMap, addr)
	switch bl.SegmentType {
	case MemoryMap.SegmentTypeEmpty:
		return 0, errors.New("Invalid address")
	case MemoryMap.SegmentTypeVirtualRAM:
		return pmc.ReadPhysicalAddress(addr)
	case MemoryMap.SegmentTypePhysicalRAM:
		return pmc.ReadPhysicalAddress(addr)
	case MemoryMap.SegmentTypePhysicalIO:
		return 0, errors.New("IO not implemented yet")
	case MemoryMap.SegmentTypeBuffer:
		return 0, errors.New("Buffer not implemented yet")
	}
	return 0, errors.New("Invalid address")
}

func (pmc *PhysicalMemoryContainer) WriteAddress(addr uint64, data byte) error {
	bl := MemoryMap.FindSegment(pmc.MyMap, addr)
	switch bl.SegmentType {
	case MemoryMap.SegmentTypeEmpty:
		return errors.New("Invalid address")
	case MemoryMap.SegmentTypeVirtualRAM:
		return pmc.WritePhysicalAddress(addr, data)
	case MemoryMap.SegmentTypePhysicalRAM:
		return pmc.WritePhysicalAddress(addr, data)
	case MemoryMap.SegmentTypePhysicalIO:
		return errors.New("IO not implemented yet")
	case MemoryMap.SegmentTypeBuffer:
		return errors.New("Buffer not implemented yet")
	}
	return errors.New("Invalid address")
}

func (pmc *PhysicalMemoryContainer) ReadPPage(page uint32) ([]byte, error) {
	bl := MemoryMap.FindSegment(pmc.MyMap, uint64(page*PageSize))
	switch bl.SegmentType {
	case MemoryMap.SegmentTypeEmpty:
		return nil, errors.New("Invalid address")
	case MemoryMap.SegmentTypeVirtualRAM:
		return pmc.ReadPhysicalPage(page)
	case MemoryMap.SegmentTypePhysicalRAM:
		return pmc.ReadPhysicalPage(page)
	case MemoryMap.SegmentTypePhysicalIO:
		return nil, errors.New("IO not implemented yet")
	case MemoryMap.SegmentTypeBuffer:
		return nil, errors.New("Buffer not implemented yet")
	}
	return nil, errors.New("Invalid address")
}

func (pmc *PhysicalMemoryContainer) WritePage(page uint32, data []byte) error {
	bl := MemoryMap.FindSegment(pmc.MyMap, uint64(page*PageSize))
	switch bl.SegmentType {
	case MemoryMap.SegmentTypeEmpty:
		return errors.New("Invalid address")
	case MemoryMap.SegmentTypeVirtualRAM:
		return pmc.WritePhysicalPage(page, data)
	case MemoryMap.SegmentTypePhysicalRAM:
		return pmc.WritePhysicalPage(page, data)
	case MemoryMap.SegmentTypePhysicalIO:
		return errors.New("IO not implemented yet")
	case MemoryMap.SegmentTypeBuffer:
		return errors.New("Buffer not implemented yet")
	}
	return errors.New("Invalid address")
}
