package VirtualMemory

type VMPage struct {
	VirutalPage  uint32
	PhysicalPage uint32
}

type VMContainer struct {
	MemoryPages []VMPage
}
