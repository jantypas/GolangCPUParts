package PhysicalMemory

import (
	"GolangCPUParts/Configuration"
	"errors"
)

const (
	PhysicalPageSize       = 4096
	MemoryType_Empty       = 0
	MemoryType_VirtualRAM  = 1
	MemoryType_PhysicalRAM = 2
	MemoryType_BufferRAM   = 3
	MemoryType_KernelRAM   = 4
	MemoryType_IORAM       = 5
	MemoryType_ROM         = 6

	Protection_NoAccess   = 0
	Protection_CanRead    = 0x1
	Protection_CanWrite   = 0x2
	Protection_CanExecute = 0x4
	Protection_NeedSystem = 0x8
)

type PhysicalMemoryBlock struct {
	Buffer       []byte
	StartAddress uint64
	EndAddress   uint64
	StartPage    uint32
	EndPage      uint32
	Protection   uint64
	MemoryType   int
	NumPages     int
	Key          int
}

type PhysicalMemoryManager struct {
	Blocks    []PhysicalMemoryBlock
	NumBlocks int
}

func PhysicalMemoryInitialize(
	cfg *Configuration.ConfigObject,
	name string) (*PhysicalMemoryManager, error) {
	// Make sure config is not null -- if it is, we have an invalid config name
	if cfg == nil {
		return nil, errors.New("Config object is nil")
	}
	sd := cfg.GetConfigByName(name)
	if sd == nil {
		return nil, errors.New("Config not found")
	}
	// Get access to the configured memory regions
	memoryRegions := sd.Description.Memory
	tatalRegions := len(memoryRegions)
	// At least one must be defined
	if tatalRegions == 0 {
		return nil, errors.New("No memory regions found")
	}
	// Make the container for all the blocks
	pmc := PhysicalMemoryManager{}
	pmc.NumBlocks = tatalRegions
	pmc.Blocks = make([]PhysicalMemoryBlock, tatalRegions)

	for idx, memoryRegion := range memoryRegions {
		pmc.Blocks[idx].StartAddress = memoryRegion.StartAddress
		pmc.Blocks[idx].EndAddress = memoryRegion.EndAddress
		pmc.Blocks[idx].Buffer = make([]byte, memoryRegion.EndAddress-memoryRegion.StartAddress)
		pmc.Blocks[idx].NumPages = int(memoryRegion.EndAddress-memoryRegion.StartAddress) / PhysicalPageSize
		pmc.Blocks[idx].StartPage = uint32(pmc.Blocks[idx].StartAddress / PhysicalPageSize)
		pmc.Blocks[idx].EndPage = uint32(pmc.Blocks[idx].EndAddress / PhysicalPageSize)
		pmc.Blocks[idx].Key = idx
		switch memoryRegion.MemoryType {
		case "Empty":
			{
				pmc.Blocks[idx].MemoryType = 0
			}
		case "Swap":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_Empty
				pmc.Blocks[idx].Protection = Protection_NoAccess
			}
		case "Virtual-RAM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_VirtualRAM
				pmc.Blocks[idx].Protection = Protection_CanWrite | Protection_CanRead | Protection_CanExecute
			}
		case "Physical-RAM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_PhysicalRAM
				pmc.Blocks[idx].Protection = Protection_CanWrite | Protection_CanRead | Protection_CanExecute
			}
		case "Buffer-RAM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_BufferRAM
				pmc.Blocks[idx].Protection = Protection_CanWrite | Protection_CanRead | Protection_CanExecute
			}
		case "Kernel-RAM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_KernelRAM
				pmc.Blocks[idx].Protection =
					Protection_CanWrite | Protection_CanRead | Protection_CanExecute | Protection_NeedSystem
			}
		case "Physical-ROM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_ROM
				pmc.Blocks[idx].Protection = Protection_CanRead | Protection_CanExecute
			}
		case "Physical-IORAM":
			{
				pmc.Blocks[idx].MemoryType = MemoryType_IORAM
				pmc.Blocks[idx].Protection = Protection_CanWrite | Protection_CanRead | Protection_CanExecute
			}
		default:
			{
				return nil, errors.New("Unknown memory type")
			}
		}
	}
	return &pmc, nil
}

func (pmc *PhysicalMemoryManager) Terminate() {
	for _, block := range pmc.Blocks {
		block.Buffer = nil
	}
}

func (pmc *PhysicalMemoryManager) GetBlockByKey(idx int) (*PhysicalMemoryBlock, error) {
	for _, block := range pmc.Blocks {
		if idx == block.Key {
			return &block, nil
		}
	}
	return nil, errors.New("Block not found")
}

func (pmc *PhysicalMemoryManager) GetBlockByType(t int) (*PhysicalMemoryBlock, error) {
	for _, block := range pmc.Blocks {
		if t == block.MemoryType {
			return &block, nil
		}
	}
	return nil, errors.New("Block not found")
}

func (pmc *PhysicalMemoryManager) GetBlockByAddress(addr uint64) (*PhysicalMemoryBlock, error) {
	for _, block := range pmc.Blocks {
		if addr >= block.StartAddress && addr <= block.EndAddress {
			return &block, nil
		}
	}
	return nil, errors.New("Block not found")
}

func (pmc *PhysicalMemoryManager) GetBlockByPage(p uint32) (*PhysicalMemoryBlock, error) {
	for _, block := range pmc.Blocks {
		if p >= block.StartPage && p <= block.EndPage {
			return &block, nil
		}
	}
	return nil, errors.New("Block not found")
}

func (pmc *PhysicalMemoryManager) GetNumberPages() uint32 {
	total := 0
	for _, block := range pmc.Blocks {
		total += block.NumPages
	}
	return uint32(total)
}

func (pmc *PhysicalMemoryManager) ReadPage(p uint32) ([]byte, error) {
	block, err := pmc.GetBlockByPage(p)
	if err != nil {
		return nil, err
	}
	return block.Buffer[p*PhysicalPageSize : (p+1)*PhysicalPageSize], nil
}

func (pmc *PhysicalMemoryManager) WritePage(p uint32, data []byte) error {
	block, err := pmc.GetBlockByPage(p)
	if err != nil {
		return err
	}
	copy(block.Buffer[p*PhysicalPageSize:], data)
	return nil
}

func (pmc *PhysicalMemoryManager) ReadAddress(addr uint64) (uint8, error) {
	block, err := pmc.GetBlockByAddress(addr)
	if err != nil {
		return 0, err
	}
	return block.Buffer[addr-block.StartAddress], nil
}

func (pmc *PhysicalMemoryManager) WriteAddress(addr uint64, data uint8) error {
	block, err := pmc.GetBlockByAddress(addr)
	if err != nil {
		return err
	}
	block.Buffer[addr-block.StartAddress] = data
	return nil
}
