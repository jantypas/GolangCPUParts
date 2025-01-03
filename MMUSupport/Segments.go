package MMUSupport

import "syscall"

type Segment struct {
	Name         string
	Protection   uint64
	UID          int
	GID          int
	System       bool
	VirtualPages []int
}

type SegmentRequest struct {
	Name       string
	Protection uint64
	UID        int
	GID        int
	System     bool
	NumPages   int
}

type SegmentList struct {
	NumSegments int
	Segments    []Segment
}

func NewSegmentList() *SegmentList {
	return &SegmentList{
		NumSegments: 0,
		Segments:    make([]Segment, 0),
	}
}

func (mmu *MMUStruct) CreateNewSegment(seg SegmentRequest) (*Segment, error) {
	lst, err := mmu.AllocateBulkPages(seg.NumPages)
	if err != nil {
		return nil, err
	}
	segment := Segment{
		Name:         seg.Name,
		Protection:   seg.Protection,
		UID:          seg.UID,
		GID:          seg.GID,
		System:       seg.System,
		VirtualPages: lst,
	}
	return &segment, nil
}

func (mmu *MMUStruct) CerateSegment(seg SegmentRequest) error {
	sp := Segment{
		Name:         seg.Name,
		Protection:   seg.Protection,
		UID:          seg.UID,
		GID:          seg.GID,
		System:       seg.System,
		VirtualPages: make([]int, 0),
	}
	pages, err := mmu.AllocateBulkPages(seg.NumPages)
	if err != nil {
		return err
	}
	sp.VirtualPages = pages
	return nil
}

func (mmu *MMUStruct) DeleteSegment(seg Segment) error {
	return mmu.FreeBulkPages(seg.VirtualPages)
}
