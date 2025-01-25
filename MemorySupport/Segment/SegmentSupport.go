package Segment

import (
	"GolangCPUParts/MemorySupport/LinearMemory"
	"errors"
)

type SegmentDescriptor struct {
	Name        string
	Readable    bool
	Writable    bool
	Executable  bool
	Protect     bool
	Size        uint64
	BaseAddress uint64
	SegmentType int
	SegmentData interface{}
}

const (
	PageSize               = 4096
	EmptySegment           = 0
	LinearSegment          = 1
	IOSegment              = 2
	BufferSegment          = 3
	VirtualSegment         = 4
	ProtectedLinearSegment = 5
)

func NewSegment(name string,
	read bool, write bool, execute bool, protect bool,
	size uint64, baseAddress uint64) (*SegmentDescriptor, error) {
	seg := SegmentDescriptor{
		Name:        name,
		Readable:    read,
		Writable:    write,
		Executable:  execute,
		Size:        size,
		BaseAddress: baseAddress,
	}
	switch seg.SegmentType {
	case LinearSegment:
		obj, err := LinearMemory.NewLinearMemory(&seg)
		if err != nil {
			return nil, err
		}
		seg.SegmentData = obj
		return &seg, nil
	default:
		return nil, errors.New("Invalid segment type")
	}
}

func (seg *SegmentDescriptor) Dispose() error {
	switch seg.SegmentType {
	case LinearSegment:
		lmb := seg.SegmentData.(LinearMemory.LinearMemoryBlock)
		lmb.Dispose()
		return nil
	default:
		return errors.New("Invalid segment type")
	}
}

func GetPage(addr uint64) uint32 {
	return uint32(addr&0x0_0000_0_FFF_FFFF_000) >> 12
}

func GetOffset(addr uint64) uint32 {
	return uint32(addr & 0x0_0000_0_000_0000_FFF)
}

func (seg *SegmentDescriptor) ReadAddress(addr uint64) (byte, error) {
	switch seg.SegmentType {
	case LinearSegment:
		lmb := seg.SegmentData.(LinearMemory.LinearMemoryBlock)
		page := GetPage(addr)
		offset := GetOffset(addr)
		return lmb.ReadAddress(page, offset)
	default:
		return 0, errors.New("Invalid segment type")
	}
}

func (seg *SegmentDescriptor) WriteAddress(addr uint64, value byte) error {
	switch seg.SegmentType {
	case LinearSegment:
		lmb := seg.SegmentData.(LinearMemory.LinearMemoryBlock)
		page := GetPage(addr)
		offset := GetOffset(addr)
		return lmb.WriteAddress(page, offset, value)
	default:
		return errors.New("Invalid segment type")
	}
}

func (seg *SegmentDescriptor) ReadPage(addr uint64) ([]byte, error) {
	switch seg.SegmentType {
	case LinearSegment:
		lmb := seg.SegmentData.(LinearMemory.LinearMemoryBlock)
		page := GetPage(addr)
		return lmb.ReadPage(page)
	default:
		return nil, errors.New("Invalid segment type")
	}
}

func (seg *SegmentDescriptor) WritePage(addr uint64, data []byte) error {
	switch seg.SegmentType {
	case LinearSegment:
		lmb := seg.SegmentData.(LinearMemory.LinearMemoryBlock)
		page := GetPage(addr)
		return lmb.WritePage(page, data)
	default:
		return errors.New("Invalid segment type")
	}
}
