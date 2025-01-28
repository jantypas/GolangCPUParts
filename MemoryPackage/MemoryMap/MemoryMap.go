package MemoryMap

import (
	"fmt"
	"sort"
)

const (
	SegmentTypeEmpty       = 0x0000
	SegmentTypeVirtualRAM  = 0x0001
	SegmentTypePhysicalRAM = 0x0002
	SegmentTypePhysicalIO  = 0x0003
	SegmentTypeBuffer      = 0x0004

	ProtectionWritable   = 0x0001
	ProtectionExecutable = 0x0002
	ProtectionSystem     = 0x0004
	ReplaceableRange     = 0x0008
	SegmentLocked        = 0x0010
)

var SegmentTypeNames = []string{
	"Empty",
	"Virtual RAM",
	"Physical RAM",
	"Physical IO",
	"Buffer",
}

type MemoryMapRegion struct {
	Key          uint16
	Comment      string
	Tag          string
	StartAddress uint64
	EndAddress   uint64
	Permissions  uint64
	SegmentType  uint16
}

var MMUTable = map[uint16][]MemoryMapRegion{}

func FindSegment(mr []MemoryMapRegion, addr uint64) *MemoryMapRegion {
	// Target address we want to search for
	target := uint64(0x0000_0000_0150_0000)

	// Perform the binary search using sort.Search
	index := sort.Search(len(mr), func(i int) bool {
		// Returns true when the current region's StartAddress >= target.
		return mr[i].StartAddress > addr
	})

	// Determine if the target is within a range
	if index > 0 &&
		mr[index-1].StartAddress <= target &&
		mr[index-1].EndAddress >= addr {
		return &mr[index-1]
	} else {
		return nil
	}
}

func DumpMemoryMap(m []MemoryMapRegion) {
	fmt.Println("--- Dumping Memory Map ---")
	for _, v := range m {
		pages := int(v.EndAddress-v.StartAddress) / 4096
		fmt.Println("Start    End      Pages")
		fmt.Printf("%08X-%08X  %4X\n", v.StartAddress, v.EndAddress, pages)
		fmt.Printf("\tSegment Type: %s\n", SegmentTypeNames[v.SegmentType])
		fmt.Printf("\tPermissions : %04X\n", v.Permissions)
		fmt.Printf("\tKey         : %04X\n", v.Key)
		fmt.Printf("\tTag         : %s\n", v.Tag)
		fmt.Printf("\tComment     : %s\n", v.Comment)
	}
	fmt.Println("--- End of Memory Map ---")
}
