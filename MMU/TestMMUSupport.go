package main

import (
	"GolangMMU/MMU"
	"fmt"
)

func main() {
	conf := MMU.MMUConfig{
		VirtualMemoryPages:  1024 * 1024,
		PhysicalMemoryPages: 12,
		TLBSize:             256,
		MaxDiskPages:        1024 * 1024,
		MinEvictPages:       4,
	}
	mmu := MMU.NewMMU(conf)

	// Access some pages to trigger page faults and swapping
	for i := 0; i < 64; i++ { // Access virtual pages beyond physical memory
		virtualAddr := i * MMU.PageSize
		err := mmu.Write(virtualAddr, byte(i))
		if err != nil {
			return
		}
		value, err := mmu.Read(virtualAddr)
		if err != nil {
			fmt.Printf("Error reading virtual address %x: %v\n", virtualAddr, err)
		} else {
			fmt.Printf("Value at virtual address %x: %x\n", virtualAddr, value)
		}
	}

	// Show statistics
	fmt.Printf("TLB Hits: %d\n", mmu.TLBHitCount)
	fmt.Printf("TLB Misses: %d\n", mmu.TLBMissCount)
	fmt.Printf("Page Faults: %d\n", mmu.PageFaultCount)
	fmt.Printf("Swaps: %d\n", mmu.SwapCount)
}
