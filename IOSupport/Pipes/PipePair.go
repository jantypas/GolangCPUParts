package Pipes

import (
	"container/list"
	"errors"
	"time"
)

const MaxPipePairs = 64

type PipePair struct {
	Master	*CList
	Slave *CList
	MasterIsOpen	bool
	SlaveIsOpen		bool
}

type PipePairTable struct {
	Table []PipePairTable
	Acitve uint64
}

func NewPipePair() *PipePair {
	return &PipePair{
		Master: NewCList(),
		Slave: NewCList(),
		MasterIsOpen: false,
		SlaveIsOpen: false,
	}
}

func (p *PipePair) Destory() {
	p.Master.Flush()
	p.Slave.Flush()
	p.Master.Dispose()
	p.Slave.Dispose()
}

func (p. *PipePair) OpenMaster() {
	p.MasterIsOpen = true
}

func (p *PipePair) CloseMaster() {
	p.MasterIsOpen = false
}

func (p *PipePair) IsMasterOpen() bool {
	return p.MasterIsOpen
}

func (p *PipePair) MasterSize() int {
	return p.Master.Size()
}

func (p *PipePair) OpenSlave() {
	p.SlaveIsOpen = true
}

func (p *PipePair) CloseSlave() {
	p.SlaveIsOpen = false
}

func (p *PipePair) IsSlaveOpen() bool {
	return p.SlaveIsOpen
}

func (p *PipePair) SlaveSize() int  {
	return p.Slave.Size()
}

func (p *PipePair) ReadFromMaster(n int) (int, []byte, error) {
	if !p.MasterIsOpen {
		return 0, nil, errors.New("Master not open")
	}
	sz, bp := p.Master.ReadNBytes(n)
	return sz, bp, nil
}

func (p *PipePair) ReadFromSlave(n int) (int, []byte, error) {
	if !p.SlaveIsOpen {
		return 0, nil, errors.New("Slave not open")
	}
	sz, bp := p.Slave.ReadNBytes(n)
	return sz, bp, nil
}

func (p *PipePair) WriteToMaster(data []byte) error {
	if !p.SlaveIsOpen {
		return errors.New("Master not open")
	}
	p.Master.WriteNBytes(data)
	return nil
}

func (p *PipePair) WriteToSlave(data []byte) error {
	if !p.MasterIsOpen {
		return errors.New("Slave not open")
	}
	s.Slave.WriteNBytes(data)
	return nil
}

func NewPipePairTable() *PipePairTable {
	return &PipePairTable{
		Table:  make([]PipePairTable, MaxPipePairs),
		Acitve: 0,
	}
}

func (pt *PipePairTable) AllocatePipePair() (int, *PipePair) {
	x : = 0
	for i := 0; i < MaxPipePairs; i++ {
		if pt.Acitve&(1<<i) == 1 {
			x = i
			break
		}
	}
	pp := NewPipePair()
	return x, pp
}

func (pt *PipePairTable) FreePipePair(x int) {
	pt.Acitve = pt.Acitve & ^(1<<x)
	pp := pt.Table[x]
	pp.FreePipePair()
}

