package MMU

import (
	"errors"
	"time"
)

func ProcessTableInitialize(mconf *MMUConfig) (*ProcessTable, error) {
	// Create a virtual memory object
	mmu, err := VirtualMemoryInitialize(mconf)
	if err != nil {
		return nil, err
	}
	// Create the process table
	pt := ProcessTable{
		MMU:         mmu,
		ProcessList: make(map[int]ProcessObject),
		NextPID:     1,
	}
	return &pt, nil
}

func ProcessTableTerminate(pt *ProcessTable) error {
	err := pt.MMU.VirtualMemoryTerminate()
	if err != nil {
		return err
	}
	return nil
}

func (pt *ProcessTable) CreateNewProcess(
	name string, args []string,
	ppid int, gid int) error {
	po := ProcessObject{
		Name:         name,
		Args:         args,
		PID:          pt.NextPID,
		PPID:         ppid,
		GID:          gid,
		SegmentTable: make([]Segment, 0),
		State:        ProcessStateWaitingToRun,
		System:       false,
		CreatedOn:    time.Now(),
	}
	pt.ProcessList[pt.NextPID] = po
	pt.NextPID++
	return nil
}

func (pt *ProcessTable) DestroyProcess(pid int) error {
	po, ok := pt.ProcessList[pid]
	if !ok {
		return errors.New("invalid process")
	}
	for i, s := range po.SegmentTable {
		err := pt.MMU.FreeBulkPages(s.Memory)
		if err != nil {
			return err
		}
	}
	delete(pt.ProcessList, pid)
	return nil
}

func (pt *ProcessTable) GrowSegment(pid int, gid int, prot int, newPages int) error {
	po, ok := pt.ProcessList[pid]
	if !ok {
		return errors.New("invalid process")
	}
	pagelist, err := pt.MMU.AllocateBulkPages(pid, gid, prot, newPages)
	if err != nil {
		return err
	}
	po.SegmentTable[0].Memory = append(po.SegmentTable[0].Memory, pagelist...)
	return nil
}

func (pt *ProcessTable) GetProcessInfo(pid int) (ProcessObject, error) {
	po, ok := pt.ProcessList[pid]
	if !ok {
		return ProcessObject{}, errors.New("invalid process")
	}
	return po, nil
}

func (pt *ProcessTable) GetProcessList() map[int]ProcessObject {
	return pt.ProcessList
}

func (pt *ProcessTable) ReadAddress(pid int, addr int) (byte, error) {

}

func (pt *ProcessTable) SetProcessState(pid int, state int) error {

}

func (pt *ProcessTable) WriteAddress(pid int, addr int) error {

}

func (pt *ProcessTable) LoadPagesFromBuffer(pid int, addr int, buffer []byte) error {
}

func (pt *ProcessTable) SavePagesToBuffer(pid int, addr int, buffer []byte) error {

}
