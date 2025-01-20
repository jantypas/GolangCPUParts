package PhysicalMemory

import (
	"GolangCPUParts/RemoteLogging"
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
}
