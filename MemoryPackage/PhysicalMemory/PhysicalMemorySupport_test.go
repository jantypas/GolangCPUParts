package PhysicalMemory

import (
	"GolangCPUParts/MemoryPackage/MemoryMap"
	"GolangCPUParts/RemoteLogging"
	"testing"
)

func TestPhysicalMemoryContainer_GetRegionByAddress(t *testing.T) {
	RemoteLogging.LogInit("test")
	TestMap := []MemoryMap.MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  MemoryMap.ReplaceableRange,
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
	}
	pmc, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
	reg := pmc.GetRegionByAddress(0x0000_0000_0001_0001)
	if reg.Key != 1 {
		t.Error("Failed to find region")
	}
	reg = pmc.GetRegionByAddress(0x0000_9999_9999_9999)
	if reg != nil {
		t.Error("Failed to report no region")
	}
}

func TestPhysicalMemoryContainer_GetRegionByKey(t *testing.T) {

}

func TestPhysicalMemoryContainer_ReadAddress(t *testing.T) {

}

func TestPhysicalMemoryContainer_ReadPage(t *testing.T) {

}

func TestPhysicalMemoryContainer_Terminate(t *testing.T) {

}

func TestPhysicalMemoryContainer_WriteAddress(t *testing.T) {

}

func TestPhysicalMemoryContainer_WritePage(t *testing.T) {

}

func TestPhysicalMemoryInitialize(t *testing.T) {
	RemoteLogging.LogInit("test")
	TestMap := []MemoryMap.MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  MemoryMap.ReplaceableRange,
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
	}
	_, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
}
