package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
	"fmt"
	"strconv"
	"testing"
)

func TestPhysicalMemory_Initialize(t *testing.T) {
	RemoteLogging.LogInit("TEST")
	pmc, err := PhysicalMemory_Initialize("TEST")
	if err != nil {
		t.Errorf("PhysicalMemory_Initialize() error = %v", err)
		return
	}
	if pmc == nil {
		t.Errorf("PhysicalMemory_Initialize() pmc = %v", pmc)
	}
	for i := 0; i < len(pmc.MemoryPages); i++ {
		t := pmc.GetMemoryType(uint32(i))
		fmt.Println("Page " + strconv.Itoa(i) + " type " + MemoryTypeNames[t])
	}
	pmc.Terminate()
}

func TestTerminate(t *testing.T) {
	RemoteLogging.LogInit("TEST")
	pmc, err := PhysicalMemory_Initialize("TEST")
	if err != nil {
		t.Errorf("PhysicalMemory_Initialize() error = %v", err)
		return
	}
	if pmc == nil {
		t.Errorf("PhysicalMemory_Initialize() pmc = %v", pmc)
	}
	for i := 0; i < len(pmc.MemoryPages); i++ {
		t := pmc.GetMemoryType(uint32(i))
		fmt.Println("Page " + strconv.Itoa(i) + " type " + MemoryTypeNames[t])
	}
}

func TestPhysicalMemoryContainer_AllocateVirtualPage(t *testing.T) {
	RemoteLogging.LogInit("TEST")
	pmc, err := PhysicalMemory_Initialize("TEST")
	if err != nil {
		t.Errorf("PhysicalMemory_Initialize() error = %v", err)
		return
	}
	if pmc == nil {
	}
}
