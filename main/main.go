package main

import (
	"GolangCPUParts/Configuration"
	"GolangCPUParts/MemoryPackage/VirtualMemory"
)

func main() {
	cfg := Configuration.ConfigObject{}
	cfg.Initialize("")
	vmc, err := VirtualMemory.VirtualMemoryInitialize(cfg, "TEST-MAP")
	if err != nil {
		panic(err)
	} else {
		vmc.Terminate()
	}
}
