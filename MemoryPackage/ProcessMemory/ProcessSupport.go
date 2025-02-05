package ProcessMemory

type ProcessSegment struct {
	MemoryPages []uint32
	Protection  uint64
	SegmentBase uint64
	SegmentSize uint64
}

type ProcessMemory struct {
	ProcessSegments []ProcessSegment
}
