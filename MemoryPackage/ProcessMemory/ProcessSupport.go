package ProcessMemory

type ProcessSegment struct {
	MemoryPages []uint32
	Protection  uint64
	SegmentBase uint64
	SegmentSize uint64
	HumPages    uint32
}

type SegmentRequest struct {
	Key        uint
	Protection uint64
}

type SegmentTable map[uint32]ProcessSegment

func AllocateSegment(numVPages uint32, req []SegmentRequest) (*ProcessSegment, error) {

	return &ps, nil
}
