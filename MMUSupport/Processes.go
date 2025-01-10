package MMUSupport

import "time"

const (
	SegmentTypeLinear  = 1
	SegmentTypeVirtual = 2
)

type ProcessObject struct {
	UID       uint
	GID       uint
	CreatedOn time.Time
	State     uint64
	Name      string
	Args      []string
	Segments  []SegmentObject
}

type SegmentObject struct {
	Start         uint32
	End           uint32
	Permissions   uint64
	SegmentType   int
	NumPages      uint32
	VirtualPages  []uint32
	PhysicalPages []byte
}

type ProcessTable struct {
	ProcessList map[uint16]ProcessObject
	NextPID     uint16
	MMU         MMUStruct
}

func ProcessTableInitialize(mconf *MMUConfig) (*ProcessTable, error) {
	mmu, err := VirtualMemoryInitialize(mconf)
	if err != nil {
		return nil, err
	}
	pt := &ProcessTable{
		ProcessList: make(map[uint16]ProcessObject),
		NextPID:     0,
		MMU:         mmu,
	}
	return pt, nil
}

func (pt *ProcessTable) Terminate() {
	for i, _ := range pt.ProcessList {
		pt.KillProcess(i)
	}
}

func (pt *ProcessTable) KillProcess(pid uint16) {
	for i, _ := range pt.ProcessList[pid].Segments {
		pt.ReleaseSegment(pid, i)
	}
	delete(pt.ProcessList, pid)
}

func (pt *ProcessTable) ReleaseSegment(pid uint16, index int) error {
	if pt.ProcessList[pid].Segments[index].SegmentType == SegmentTypeVirtual {
		err := pt.MMU.FreeBulkPages(pt.ProcessList[pid].Segments[index].VirtualPages)
		if err != nil {
			return err
		}
	} else {
		pt.ProcessList[pid].Segments[index].PhysicalPages = nil
	}
	return nil
}

func (pt *ProcessTable) CreateProcess(ps *ProcessObject) error {
	pso := ps
	pso.CreatedOn = time.Now()
	pso.State = ProcessStateWaitingToRun
	for _, s := range pso.Segments {
		switch s.SegmentType {
		case SegmentTypeLinear:
			s.PhysicalPages = make([]byte, s.NumPages*PageSize)
			break
		case SegmentTypeVirtual:
			pages, err := pt.MMU.AllocateBulkPages(s.NumPages)
			if err != nil {
				return err
			}
			s.VirtualPages = pages
		}
	}
	pt.ProcessList[pt.NextPID] = *pso
	pt.NextPID++
	return nil
}
