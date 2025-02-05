package VirtualMemory

import (
	"GolangCPUParts/Configuration"
	"testing"
)

func TestListFindUint32(t *testing.T) {
	cfg := Configuration.ConfigObject{}
	cfg.Initialize("Test")
	vmc, err := VirtualMemoryInitialize(cfg, "TEST-MAP")
	if err != nil {
		t.Error(err)
	}
	vmc.Terminate()
}

func TestListFindUint64(t *testing.T) {
}
