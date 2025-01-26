package MemorySuport

import (
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
)

type MemoryMapRegion struct {
	StartAddress uint64
	EndAddress   uint64
	Permissions  uint64
	SegmentType  uint64
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
