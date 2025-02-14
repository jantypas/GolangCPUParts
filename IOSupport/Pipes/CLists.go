package Pipes

import (
	"sync"
)

const CListSliceSize = 4096

type CList struct {
	Buffer   []byte
	ReadPtr  int
	WritePtr int
	Lock     sync.Mutex
}

func NewCList() *CList {
	return &CList{
		Buffer:   make([]byte, CListSliceSize),
		ReadPtr:  0,
		WritePtr: 0,
	}
}

func (c *CList) IsEmpty() bool {
	return c.ReadPtr == c.WritePtr
}

func (c *CList) Extend() {
	newBuffer := make([]byte, len(c.Buffer)+CListSliceSize)
	copy(newBuffer, c.Buffer)
	c.Buffer = newBuffer
}

func (c *CList) Trim() {
	newBuffer := make([]byte, len(c.Buffer)-CListSliceSize)
	copy(newBuffer, c.Buffer[CListSliceSize:])
	c.Buffer = newBuffer
	c.ReadPtr -= CListSliceSize
	c.WritePtr -= CListSliceSize
}

func (c *CList) WriteNBytes(data []byte) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	for c.WritePtr+len(data) >= len(c.Buffer) {
		c.Extend()
	}
	for i := 0; i < len(data); i++ {
		c.Buffer[c.WritePtr] = data[i]
		c.WritePtr++
	}
}

func (c *CList) ReadNBytes(n int) (int, []byte) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.IsEmpty() {
		return 0, nil
	}
	if n > (c.WritePtr - c.ReadPtr) {
		n = c.WritePtr - c.ReadPtr
	}
	output := make([]byte, n)
	for c.ReadPtr > CListSliceSize {
		c.Trim()
	}
	for i := 0; i < n; i++ {
		output[i] = c.Buffer[c.ReadPtr]
		c.ReadPtr++
	}
	return n, output
}

func (c *CList) Flush() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.ReadPtr = 0
	c.WritePtr = 0
	c.Buffer = make([]byte, CListSliceSize)
}
