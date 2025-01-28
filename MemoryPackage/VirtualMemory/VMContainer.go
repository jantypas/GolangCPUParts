package VirtualMemory

import (
	"GolangCPUParts/MemoryPackage/PhysicalMemory"
	"GolangCPUParts/MemoryPackage/Swapper"
	"container/list"
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

func VirtualMemoryInitialize(
	vpages uint32,
	pmc *PhysicalMemory.PhysicalMemoryContainer,
	swap *Swapper.SwapperContainer) *VMContainer {
	vmc := VMContainer{
		MemoryPages: make(map[uint32]VMPage),
	}
	vmc.MemoryPages = make(map[uint32]VMPage, vpages)
	vmc.PMemory = pmc
	vmc.Swapper = swap
	UsedVirtualPages := list.New()
	LRUCache := list.New()
	FreeVirtualPages := list.New()
	vmc.PMemory.Get
	return &vmc
}

func (vmc *VMContainer) Terminate() {
	vmc.Swapper.Terminate()
	vmc.PMemory.Terminate()
	vmc.MemoryPages = nil
}
