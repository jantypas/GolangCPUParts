package MemoryMap

import (
	"testing"
)

func TestFindSegment(t *testing.T) {
	TestMap := []MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  ReplaceableRange,
			SegmentType:  0x0,
		},
		{
			Key:          1,
			StartAddress: 0x0000_0000_0001_0000,
			EndAddress:   0x0000_0000_0001_FFFF,
			Permissions:  0x1,
			SegmentType:  0x1,
		},
		{
			Key:          2,
			StartAddress: 0x0000_0000_0002_0000,
			EndAddress:   0x0000_0000_0002_FFFF,
			Permissions:  0x2,
			SegmentType:  0x2,
		},
		{
			Key:          3,
			StartAddress: 0x0000_0000_0003_0000,
			EndAddress:   0x0000_0000_0003_FFFF,
			Permissions:  0x3,
			SegmentType:  0x3,
		},
	}

	seg := FindSegment(TestMap, 0x0000_0000_0001_0050)
	if seg == nil {
		t.Error("Failed to find segment")
	}
	if seg.StartAddress != 0x0000_0000_0001_0000 {
		t.Error("Failed to find correct segment")
	}
	seg = FindSegment(TestMap, 0x0000_0000_0005_0000)
	if seg != nil {
		t.Error("Failed to report no segment")
	}
}
