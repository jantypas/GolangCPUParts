package PhysicalMemory

import (
	"GolangCPUParts/MemorySuport/MemoryMap"
	"GolangCPUParts/RemoteLogging"
)

const PageSize = 4096

// PhysicalBlock
// The physical memory block contains a block of bytes for a physical memory region
type PhysicalBlock struct {
	Buffer      []byte // The buffer of bytes that contain data
	NumPages    uint32 // Number of pages in the buffer
	StartPage   uint32 // Where does our material start (on a page)
	Protections uint64 // Any protection rules
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
		totalBytes += mmap[i].EndAddress - mmap[i].StartAddress
		totalPages += int(mmap[i].EndAddress-mmap[i].StartAddress) / PageSize
	}
	msg := "Initialized physical memory with " + string(totalBytes) + " bytes and " + string(totalPages) + " pages"
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
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "Getting block for key "+string(key))
	// Walk the map looking for a key
	for i := range pmc.MyMap {
		if pmc.MyMap[i].Key == key {
			// Return it
			return &pmc.PhysicalBlocks[i]
		}
	}
	// No key found
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "No block found for key "+string(key))
	return nil
}

// GetRegionByAddress
// Retrieves a PhysicalBlock corresponding to the provided address in the memory map,
// or nil if not found.
func (pmc *PhysicalMemoryContainer) GetRegionByAddress(addr uint64) *PhysicalBlock {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryGetBlock", "Getting block for address "+string(addr))
	bl := MemoryMap.FindSegment(pmc.MyMap, addr)
	if bl != nil {
		return pmc.GetRegionByKey(bl.Key)
	}
	return nil
}

// ReadAddress retrieves a byte from the physical memory at the specified address.
// Returns the byte and an error if the address is invalid or not mapped.
func (pmc *PhysicalMemoryContainer) ReadAddress(addr uint64) (byte, error) {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryReadAddress", "Reading address "+string(addr))
	// Compute the block by address
	bl := pmc.GetRegionByAddress(addr)
	if bl != nil {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryReadAddress", "Invalid address "+string(addr))
		return 0, nil
	}
	// Compute the address in the buffer
	newAddr := (uint32(addr/PageSize) - bl.StartPage) * PageSize
	return bl.Buffer[newAddr], nil
}

// WriteAddress writes a byte of data to the specified physical memory address.
// Returns an error if the address is invalid or not writable.
func (pmc *PhysicalMemoryContainer) WriteAddress(addr uint64, data byte) error {
	RemoteLogging.LogEvent("INFO", "PhysicalMemoryWriteAddress", "Writing address "+string(addr))
	// Compute the block by address
	bl := pmc.GetRegionByAddress(addr)
	if bl != nil {
		RemoteLogging.LogEvent("ERROR", "PhysicalMemoryWriteAddress", "Invalid address "+string(addr))
		return nil
	}
	// Update buffer
	newAddr := (uint32(addr/PageSize) - bl.StartPage) * PageSize
	bl.Buffer[newAddr] = data
	return nil
}
