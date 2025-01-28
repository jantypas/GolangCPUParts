package VirtualMemory

import (
	"GolangCPUParts/Configuration"
	"GolangCPUParts/MemoryPackage/MemoryMap"
	"GolangCPUParts/MemoryPackage/PhysicalMemory"
	"GolangCPUParts/MemoryPackage/Swapper"
	"container/list"
	"errors"
)

type VMPage struct {
	VirutalPage  uint32
	PhysicalPage uint32
}

type VMContainer struct {
	MemoryPages      map[uint32]VMPage
	Swapper          *Swapper.SwapperContainer
	PMemory          *PhysicalMemory.PhysicalMemoryContainer
	FreeVirtualPages *list.List
	UsedVirtualPages *list.List
	LRUCache         *list.List
}

func VirtualMemoryInitialize(cfg Configuration.ConfigObject, name string, vpages uint32) (*VMContainer, error) {
	// First try to find the memory map
	mm, ok := MemoryMap.ProductionMap[name]
	if !ok {
		return nil, errors.New("Failed to find memory map")
	}
	// With the map, build the physical memory container
	pmc, err := PhysicalMemory.PhysicalMemoryInitialize(mm)
	if err != nil {
		return nil, err
	}
	// Build the virtual memory container
	vp := pmc.GetRegionByTag("VIRTUAL-RAM")
	vmc, err := VirtualMemoryInitialize(cfg, name, vp.NumPages)
	// Build the swapper the VM Container will need
	vmc.Swapper, err = Swapper.Swapper_Initialize(cfg.SwapFileNames, vmc)
	return nil, nil
}

func (vmc *VMContainer) Terminate() {

}
