package ProcessMemory

import (
	"GolangCPUParts/IOSupport/Pipes"
	"GolangCPUParts/MemoryPackage/VirtualMemory"
	"math/rand"
)

/*
   To understand how pages are handled, remember that an address can be broken down into fields.
          -- 16 bit process ID -- -- 16 big segment ID -- 20 bit page ID -- -- 12-bit offset --
   Just by looking at the address, we can determine which maps and table to cosnult.
*/

// MemoryPage represents a logical segment and its associated page in memory.
// The segment defines the protections that this page uses.
type MemoryPage struct {
	Segment int
	Page    uint32
}

// SegmentObject represents a memory segment with base page, size, and protection attributes.
// For that segment, we can see the protecitons and how large the segment is (in pages)
type SegmentObject struct {
	BasePage   int
	Size       int
	Protection uint64
}

// ProcessMemoryObject represents a collection of memory pages and the total size of memory
// It also points to the segment that manages these pages
type ProcessMemoryObject struct {
	Pages []MemoryPage
	Size  int
}

const (
	Protection_CanWrite    = 0x1 << iota
	Protection_CanExecute  = 0x2
	Protection_IsLocked    = 0x4
	Protection_NeedSystem  = 0x8
	Protection_ShadowStack = 0x10
	Protection_IsVirtual   = 0x20

	Protection_Code   = Protection_CanExecute
	Protection_Data   = Protection_CanWrite
	Protection_Stack  = Protection_CanWrite | Protection_ShadowStack
	Protection_Heap   = Protection_CanWrite
	Protection_Kernel = Protection_CanExecute | Protection_CanWrite |
							Protection_IsLocked | Protection_NeedSystem
	Protection_ROM = Protection_CanExecute
)

// ProcessTable represents a collection of segment objects and process memory objects.
// When we create a process memory segment, we're defining or using a segment, and allocating pages .
type ProcessTable struct {
	Segments      map[int]SegmentObject
	MemoryObjects map[uint32]ProcessMemoryObject
	VMC           *VirtualMemory.VMContainer
}

// Initialize the process table
func ProcessTable_Initialize(vmp *VirtualMemory.VMContainer) *ProcessTable {
	pt := ProcessTable{}
	pt.Segments = make(map[int]SegmentObject)
	pt.MemoryObjects = make(map[uint32]ProcessMemoryObject)
	pt.VMC = vmp
	return &pt
}

type PipeTable struct {
	Pipes 	[]Pipes.PipePair
}
// Terminate releases all resources associated with the process table and cleans up memory allocations.
func (pt *ProcessTable) Terminate() {
	for _, v := range pt.MemoryObjects {
		pl := make([]uint32, v.Size)
		for i := 0; i < v.Size; i++ {
			pl[i] = v[i].
		}
	}
	pt.MemoryObjects = nil
	for i, v := range pt.Segments {
		delete(pt.Segments, i)
	}
}
