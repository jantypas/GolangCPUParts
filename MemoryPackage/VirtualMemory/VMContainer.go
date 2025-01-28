package VirtualMemory

import (
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

func VirtualMemoryInitialize(name string, vpages uint32) (*VMContainer, error) {

}

func (vmc *VMContainer) Terminate() {

}
