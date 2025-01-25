package LinearMemory

import (
	"GolangCPUParts/MemorySupport/Segment"
	"errors"
)

type LinearMemoryBlock struct {
	Buffer   []byte
	NumPages uint32
}

func NewLinearMemory(seg *Segment.SegmentDescriptor) (*LinearMemoryBlock, error) {
	if seg.Size%Segment.PageSize != 0 {
		return nil, errors.New("Size must be a multiple of page size")
	}
	lmb := LinearMemoryBlock{
		NumPages: uint32(seg.Size / Segment.PageSize),
		Buffer:   make([]byte, seg.Size),
	}
	return &lmb, nil
}

func (lmb *LinearMemoryBlock) Dispose() error {
	lmb.Buffer = nil
	return nil
}

func (lmb *LinearMemoryBlock) ReadAddress(page uint32, offset uint32) (byte, error) {
	if page > lmb.NumPages {
		return 0, errors.New("Page out of bounds")
	}
	ptr := page*Segment.PageSize + offset
	return lmb.Buffer[ptr], nil
}

func (lmb *LinearMemoryBlock) WriteAddress(page uint32, offset uint32, value byte) error {
	if page > lmb.NumPages {
		return errors.New("Page out of bounds")
	}
	ptr := page*Segment.PageSize + offset
	lmb.Buffer[ptr] = value
	return nil
}

func (lmb *LinearMemoryBlock) ReadPage(page uint32) ([]byte, error) {
	if page > lmb.NumPages {
		return nil, errors.New("Page out of bounds")
	}
	return lmb.Buffer[page*Segment.PageSize : (page+1)*Segment.PageSize], nil
}

func (lmb *LinearMemoryBlock) WritePage(page uint32, data []byte) error {
	if page > lmb.NumPages {
		return errors.New("Page out of bounds")
	}
	copy(lmb.Buffer[page*Segment.PageSize:(page+1)*Segment.PageSize], data)
	return nil
}
