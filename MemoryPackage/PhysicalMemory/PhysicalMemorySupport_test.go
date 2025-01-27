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
			Permissions:  0x0,
			SegmentType:  0x0,
		},
		{
			Key:          1,
			StartAddress: 0x0000_0000_0001_0000,
			EndAddress:   0x0000_0000_0001_FFFF,
			Permissions:  0x1,
			SegmentType:  0x1,
		},
	}
	pmc, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
	reg := pmc.GetRegionByAddress(0x0000_0000_0001_0055)
	if reg.Key != 1 {
		t.Error("Failed to find correct region")
	}
	reg = pmc.GetRegionByAddress(0x0000_0000_0005_0000)
	if reg != nil {
		t.Error("Failed to report no region")
	}
}

func TestPhysicalMemoryContainer_GetRegionByKey(t *testing.T) {
	RemoteLogging.LogInit("test")
	TestMap := []MemoryMap.MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  0x0,
			SegmentType:  0x0,
		},
		{
			Key:          1,
			StartAddress: 0x0000_0000_0001_0000,
			EndAddress:   0x0000_0000_0001_FFFF,
			Permissions:  0x1,
			SegmentType:  0x1,
		},
	}
	pmc, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
	reg := pmc.GetRegionByKey(1)
	if reg.Key != 1 {
		t.Error("Failed to find correct region")
	}
	reg = pmc.GetRegionByAddress(52)
	if reg == nil {
		t.Error("Failed to report no region")
	}
}

func TestPhysicalMemoryContainer_ReadWriteTest(t *testing.T) {
	RemoteLogging.LogInit("test")
	TestMap := []MemoryMap.MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  0x0,
			SegmentType:  0x0,
		},
		{
			Key:          1,
			StartAddress: 0x0000_0000_0001_0000,
			EndAddress:   0x0000_0000_0001_FFFF,
			Permissions:  0x1,
			SegmentType:  0x1,
		},
	}
	pmc, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
	err = pmc.WriteAddress(0x0000_0000_0001_3000, 0x12)
	if err != nil {
		t.Error(err)
	}
	val, err := pmc.ReadAddress(0x0000_0000_0001_3000)
	if err != nil {
		t.Error(err)
	}
	if val != 0x12 {
		t.Error("Failed to read correct value")
	}
}

func TestPhysicalMemoryContainer_ReadWritePageTest(t *testing.T) {

}

func TestPhysicalMemoryInitialize(t *testing.T) {
	RemoteLogging.LogInit("test")
	TestMap := []MemoryMap.MemoryMapRegion{
		{
			Key:          0,
			StartAddress: 0x0000_0000_0000_0000,
			EndAddress:   0x0000_0000_0000_FFFF,
			Permissions:  0x0,
			SegmentType:  0x0,
		},
		{
			Key:          1,
			StartAddress: 0x0000_0000_0001_0000,
			EndAddress:   0x0000_0000_0001_FFFF,
			Permissions:  0x1,
			SegmentType:  0x1,
		},
	}
	_, err := PhysicalMemoryInitialize(TestMap)
	if err != nil {
		t.Error(err)
	}
}
