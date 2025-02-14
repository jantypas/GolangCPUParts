package Pipes

type PipePair struct {
	Master *CList
	Slave  *CList
}

func NewPipePair() *PipePair {
	return &PipePair{
		Master: NewCList(),
		Slave:  NewCList(),
	}
}

func (p *PipePair) Flush() {
	p.Master.Flush()
	p.Slave.Flush()
}

func (p *PipePair) IsMasterEmpty() bool {
	return p.Master.IsEmpty()
}

func (p *PipePair) IsSlaveEmpty() bool {
	return p.Slave.IsEmpty()
}

func (p *PipePair) WriteMaster(b []byte) {
	p.Master.WriteNBytes(b)
}

func (p *PipePair) WriteSlave(b []byte) {
	p.Slave.WriteNBytes(b)
}

func (p *PipePair) ReadNMaster(n int) (int, []byte) {
	return p.Master.ReadNBytes(n)
}

func (p *PipePair) ReadSlave(n int) (int, []byte) {
	return p.Slave.ReadNBytes(n)
}
