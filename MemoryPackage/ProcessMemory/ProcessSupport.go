package ProcessMemory

type ProcessSegment struct {
	MemoryPages []uint32
	Protection  uint64
	SegmentBase uint64
	SegmentSize uint64
}

type SegmentTable map[uint32]ProcessSegment

func AllocateSegment(
	numVPages uint32,
	base uint64,
	size uint64,
	prot uint64,
) (*ProcessSegment, error) {
	ps := ProcessSegment{
		Protection:  prot,
		SegmentBase: base,
		SegmentSize: size,
	}
	return &ps, nil
}
