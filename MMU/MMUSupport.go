package MMU

import (
	"errors"
	"fmt"
)

// MMUConfig
// The MMUConfig defines the parameters for our MMU
type MMUConfig struct {
	VirtualMemoryPages  int // Number of virtual memory pages we'll support
	PhysicalMemoryPages int // Number of physical pages we'll support
	TLBSize             int // How large is our TLB in pages
	MinEvictPages       int // The minimal page swap out
	SwapIO              SwapperInterface
}

// NewMMUConfig
// initializes and returns a default configuration for an MMU with predefined page and TLB settings.
func NewMMUConfig() *MMUConfig {
	swapper := SwapperInterface{Filename: "swap.bin"}
	return &MMUConfig{
		VirtualMemoryPages:  1024 * 1024,
		PhysicalMemoryPages: 16384,
		TLBSize:             256,
		MinEvictPages:       1024,
		SwapIO:              swapper,
	}
}

// Protection flags
// The bits that define the protection for a memory page
const (
	OwnerRead   = 1 << 0 // Read permission
	OwnerWrite  = 1 << 1 // Write permission
	OwnerExec   = 1 << 2 // Execute permission
	OwnerSystem = 1 << 3 // System privilege on the page
	GroupRead   = 1 << 4
	GroupWrite  = 1 << 5
	GroupExec   = 1 << 6
	GroupSystem = 1 << 7
	WorldRead   = 1 << 8
	WorldWrite  = 1 << 9
	WorldExec   = 1 << 10
	WorldSystem = 1 << 11

	PageSize = 4096
)

// PageTableEntry
// The PageTableEntry represents a single entry in the virtual page table
type PageTableEntry struct {
	Present      bool // Indicates if the page is in physical memory
	OnDisk       bool // Indicates if the page resides on the disk
	Dirty        bool // Page should be written to disk
	DiskPage     int  // Page index on the simulated disk (if on disk)
	Protection   int  // Page protection flags (read/write/execute)
	PhysicalPage int  // Physical page index (if in memory)
	ProcessID    int
	GroupID      int
}

// TLBEntry represents a single TLB entry
// The TLB is used as a cache to avoid digging into the page table
type TLBEntry struct {
	VirtualPage  int  // Virtual page index
	PhysicalPage int  // Physical page index
	Valid        bool // Indicates if the entry is valid
}

// MMU
// MMU represents a Memory Management Unit with TLB and disk swapping support
type MMU struct {
	PageTable      []PageTableEntry // Page table
	TLB            []TLBEntry       // Translation Lookaside Buffer
	PhysicalMem    []byte           // Physical memory representation
	Disk           [][]byte         // Simulated disk storage
	FreePages      []int            // List of free physical pages
	FreeDiskSlots  []int            // List of free disk slots
	TLBHitCount    int              // Count of TLB hits
	TLBMissCount   int              // Count of TLB misses
	PageFaultCount int              // Count of page faults
	SwapCount      int              // Count of pages swapped to/from memory
	MinEvictPages  int
	Swapper        SwapperInterface
}

// NewMMU
// Initializes a new MMU instance with disk support
func NewMMU(cnf MMUConfig) *MMU {
	numVirtualPages := cnf.VirtualMemoryPages
	numPhysicalPages := cnf.PhysicalMemoryPages

	// Initialize free physical pages
	freePages := make([]int, numPhysicalPages)
	for i := 0; i < numPhysicalPages; i++ {
		freePages[i] = i
	}

	// Initialize TLB
	tlb := make([]TLBEntry, cnf.TLBSize)

	// We've built the MMU
	return &MMU{
		PageTable:     make([]PageTableEntry, numVirtualPages),
		TLB:           tlb,
		PhysicalMem:   make([]byte, cnf.PhysicalMemoryPages*PageSize),
		FreePages:     freePages,
		Swapper:       cnf.SwapIO,
		MinEvictPages: cnf.MinEvictPages,
	}
}

// findInTLB checks if a virtual page exists in the TLB
func (mmu *MMU) findInTLB(virtualPage int) (int, bool) {
	// Look through the TLB for a page match
	for _, entry := range mmu.TLB {
		if entry.Valid && entry.VirtualPage == virtualPage {
			// Mark this page as used
			mmu.TLBHitCount++
			return entry.PhysicalPage, true
		}
	}
	// We didn't find a page
	mmu.TLBMissCount++
	return 0, false
}

// updateTLB
// Updates the TLB with a new virtual-to-physical page mapping
func (mmu *MMU) updateTLB(virtualPage, physicalPage int) {
	// For each entry in the TLB
	for i := range mmu.TLB {
		if !mmu.TLB[i].Valid {
			mmu.TLB[i] = TLBEntry{
				VirtualPage:  virtualPage,
				PhysicalPage: physicalPage,
				Valid:        true,
			}
			return
		}
	}
	// FIFO eviction policy
	mmu.TLB = append(mmu.TLB[1:], TLBEntry{
		VirtualPage:  virtualPage,
		PhysicalPage: physicalPage,
		Valid:        true,
	})
}

func (mmu *MMU) evictPage() (int, error) {
	for i := 0; i < mmu.MinEvictPages; i++ {
		if _, err := mmu.evictSinglePage(); err != nil {
			return 0, err
		}
	}
	return 0, nil
}

// evictPage handles evicting a page from memory to disk
func (mmu *MMU) evictSinglePage() (int, error) {
	// Evict the first in-memory page (simple FIFO eviction)
	for i, entry := range mmu.PageTable {
		if entry.Present {
			// Copy physical memory to disk
			physicalPage := entry.PhysicalPage
			err := mmu.Swapper.SwapOut(
				int64(physicalPage*PageSize),
				mmu.PhysicalMem[physicalPage*PageSize:(physicalPage+1)*PageSize])
			if err != nil {
				return 0, err
			}
			// Mark the page as swapped out
			entry.Present = false
			entry.OnDisk = true
			mmu.FreePages = append(mmu.FreePages, physicalPage)

			mmu.PageTable[i] = entry
			fmt.Printf("Page swapped out: Virtual page %d to disk page %d\n", evictedPage, diskSlot)
			mmu.SwapCount++
			return 0, nil
		}
	}
	return 0, errors.New("no pages available for eviction")
}

// handlePageFault handles a page fault
func (mmu *MMU) handlePageFault(virtualPage int) error {
	if len(mmu.FreePages) == 0 {
		// If physical memory is full, evict a page
		if _, err := mmu.evictPage(); err != nil {
			return err
		}
	}

	// Allocate a physical page
	physicalPage := mmu.FreePages[len(mmu.FreePages)-1]
	mmu.FreePages = mmu.FreePages[:len(mmu.FreePages)-1]

	entry := &mmu.PageTable[virtualPage]
	// If the page was on disk, load it back into memory
	if entry.OnDisk {
		tmpBuffer := make([]byte, PageSize)
		err := mmu.Swapper.SwapIn(int64(physicalPage*PageSize), tmpBuffer)
		copy(mmu.PhysicalMem[physicalPage*PageSize:physicalPage*PageSize+PageSize], tmpBuffer)
		entry.OnDisk = false
		entry.DiskPage = -1
		mmu.FreeDiskSlots = append(mmu.FreeDiskSlots, diskPage)
		fmt.Printf("Page loaded from disk: Virtual page %d from disk page %d to physical page %d\n", virtualPage, diskPage, physicalPage)
	} else {
		// If it's a new allocation, initialize it
		fmt.Printf("Page fault handled: Allocated virtual page %d to physical page %d\n", virtualPage, physicalPage)
	}

	entry.Present = true
	entry.PhysicalPage = physicalPage
	mmu.PageFaultCount++
	return nil
}

// Translate translates a virtual address to a physical address
func (mmu *MMU) Translate(virtualAddr int, accessType int) (int, error) {
	virtualPage := virtualAddr / PageSize
	offset := virtualAddr % PageSize

	if virtualPage >= len(mmu.PageTable) {
		return 0, errors.New("invalid virtual address")
	}

	// Check TLB for the mapping
	if physicalPage, found := mmu.findInTLB(virtualPage); found {
		return physicalPage*PageSize + offset, nil
	}

	// Page table lookup
	entry := &mmu.PageTable[virtualPage]
	if !entry.Present {
		// Handle page fault
		if err := mmu.handlePageFault(virtualPage); err != nil {
			return 0, err
		}
		entry = &mmu.PageTable[virtualPage]
	}

	// Update the TLB
	mmu.updateTLB(virtualPage, entry.PhysicalPage)

	return entry.PhysicalPage*PageSize + offset, nil
}

// Write writes a byte to a virtual address
func (mmu *MMU) Write(virtualAddr int, value byte) error {
	physicalAddr, err := mmu.Translate(virtualAddr, Write)
	if err != nil {
		return err
	}
	mmu.PhysicalMem[physicalAddr] = value
	return nil
}

// Read reads a byte from a virtual address
func (mmu *MMU) Read(virtualAddr int) (byte, error) {
	physicalAddr, err := mmu.Translate(virtualAddr, Read)
	if err != nil {
		return 0, err
	}
	return mmu.PhysicalMem[physicalAddr], nil
}
