package PhysicalMemory

import "errors"

const BufferMultiple = 8

type BaseBuffer struct {
	Buffer []byte
}

func (bb *BaseBuffer) MakeBaseBuffer() *BaseBuffer {
	return &BaseBuffer{
		Buffer: make([]byte, 4096),
	}
}

type BufferMemoryObject struct {
	BufferSpace [BufferMultiple]BaseBuffer
	ReadPtr     int
	WritePtr    int
}

func NewBufferMemoryObject() *BufferMemoryObject {
	b := &BufferMemoryObject{
		BufferSpace: [BufferMultiple]BaseBuffer{},
		ReadPtr:     0,
		WritePtr:    0,
	}
	for i := 0; i < BufferMultiple; i++ {
		b.BufferSpace[i].MakeBaseBuffer()
	}
	return b
}

func (b *BufferMemoryObject) ReadPage() (*[]byte, error) {
	if b.ReadPtr == b.WritePtr {
		return nil, errors.New("Buffer is empty")
	}
	if b.ReadPtr == BufferMultiple {
		b.ReadPtr = 0
	}
	result := &b.BufferSpace[b.ReadPtr].Buffer
	b.ReadPtr++
	return result, nil
}

func (b *BufferMemoryObject) WritePage(v []byte) error {
	if b.ReadPtr == b.WritePtr {
		return errors.New("Buffer is empty")
	}
	if b.ReadPtr == BufferMultiple {
		b.WritePtr = 0
	}
	copy(b.BufferSpace[b.WritePtr].Buffer, v)
	b.WritePtr++
	return nil
}

func (b *BufferMemoryObject) BufferIsEmpty() bool {
	return b.ReadPtr == b.WritePtr
}
