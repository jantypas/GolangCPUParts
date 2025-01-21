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
		t.Errorf("PhysicalMemory_Initialize() pmc = %v", pmc)
	}
	page, err := pmc.AllocateVirtualPage()
	if err != nil {
		t.Errorf("AllocateVirtualPage() error = %v", err)
		return
	}
	for i := 0; i < 10; i++ {
		page, err = pmc.AllocateVirtualPage()
		if err != nil {
			if err.Error() != "No free virtual pages" {
				t.Errorf("AllocateVirtualPage() error = %v", err)
				return
			}
		} else {
			fmt.Println("Page " + strconv.Itoa(int(page)) + " allocated")
		}
	}
	fmt.Println("We should have four virtual pages")
	num := pmc.ReturnListOfPageType(MemoryTypeVirtualRAM)
	if num.Len() != 4 {
		t.Errorf("ReturnListOfPageType() = %v, want %v", num, 4)
	}
}

func TestPhysicalMemoryContainer_FreeVirtualPage(t *testing.T) {
	RemoteLogging.LogInit("TEST")
	pmc, err := PhysicalMemory_Initialize("TEST")
	if err != nil {
		t.Errorf("PhysicalMemory_Initialize() error = %v", err)
		return
	}
	if pmc == nil {
		t.Errorf("PhysicalMemory_Initialize() pmc = %v", pmc)
	}
	page, err := pmc.AllocateVirtualPage()
	if err != nil {
		t.Errorf("AllocateVirtualPage() error = %v", err)
		return
	}
	err = pmc.ReturnVirtualPage(uint32(page))
	if err != nil {
		t.Errorf("ReturnVirtualPage() error = %v", err)
	}
	err = pmc.ReturnVirtualPage(uint32(12))
	if err.Error() != "Page wrong type" {
		t.Errorf("ReturnVirtualPage() error = %v", err)
	}
}

func TestPhysicalMemoryContainer_ReadWriteTest(t *testing.T) {
	RemoteLogging.LogInit("TEST")
	pmc, err := PhysicalMemory_Initialize("TEST")
	if err != nil {
		t.Errorf("PhysicalMemory_Initialize() error = %v", err)
		return
	}
	if pmc == nil {
		t.Errorf("PhysicalMemory_Initialize() pmc = %v", pmc)
	}
	page, err := pmc.AllocateVirtualPage()
	fmt.Println("Page " + strconv.Itoa(int(page)) + " allocated")
	err = pmc.WriteAddress(512, 12)
	if err != nil {
		t.Errorf("WriteAddress() error = %v", err)
	}
	v, err := pmc.ReadAddress(512)
	if err != nil {
		t.Errorf("ReadAddress() error = %v", err)
	}
	if v != 12 {
		t.Errorf("ReadAddress() = %v, want %v", v, 12)
	}
}
