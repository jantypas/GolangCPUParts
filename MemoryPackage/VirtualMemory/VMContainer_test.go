package VirtualMemory

import (
	"GolangCPUParts/Configuration"
	"GolangCPUParts/RemoteLogging"
	"fmt"
	"strconv"
	"testing"
)

func TestListFindUint32(t *testing.T) {
	RemoteLogging.LogInit("test")
	cfg := Configuration.ConfigObject{}
	cfg.Initialize("Test")
	vmc, err := VirtualMemoryInitialize(cfg, "TEST-MAP")
	if err != nil {
		t.Error(err)
	}
	vmc.Terminate()
}

func TestVirtualMemoryInitializeNil(t *testing.T) {
	RemoteLogging.LogInit("test")
	cfg := Configuration.ConfigObject{}
	cfg.Initialize("Test")
	vmc, err := VirtualMemoryInitialize(cfg, "FTEST-MAP")
	if err == nil {
		t.Error(err)
	} else {
		e := vmc.Terminate()
		if e == nil {
			t.Error(e)
		}
	}
}

func TestVirtualMemoryAllocate(t *testing.T) {
	RemoteLogging.LogInit("test")
	cfg := Configuration.ConfigObject{}
	cfg.Initialize("Test")
	vmc, err := VirtualMemoryInitialize(cfg, "TEST-MAP")
	if err != nil {
		t.Error(err)
	}
	lst, err := vmc.AllocateNVirtualPages(5)
	if err != nil {
		t.Error(err)
	}
	if lst.Len() != 5 {
		t.Error("Failed to allocate correct number of pages")
	}
	for l := lst.Front(); l != nil; l = l.Next() {
		fmt.Println("Page " + strconv.Itoa(int(l.Value.(uint32))))
	}
	vmc.Terminate()
}
