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
		Name:      name,
		Args:      args,
		PID:       pt.NextPID,
		PPID:      ppid,
		GID:       gid,
		Memory:    make([]int, 0),
		State:     ProcessStateWaitingToRun,
		System:    false,
		CreatedOn: time.Now(),
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
	err := pt.MMU.FreeBulkPages(po.Memory)
	if err != nil {
		return err
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
	po.Memory = append(po.Memory, pagelist...)
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

func (pt *ProcessTable) SetProcessState(pid int, state int) error {
	po, ok := pt.ProcessList[pid]
	if !ok {
		return errors.New("invalid process")
	}
	po.State = state
	pt.ProcessList[pid] = po
	return nil
}

func (pt *ProcessTable) ReadAddress(
	uid int, gid int,
	mode int, seg int,
	pid int, addr int) (byte, error) {
	_, ok := pt.ProcessList[pid]
	if !ok {
		return 0, errors.New("invalid process")
	}
	page := addr / PageSize
	if page > pt.MMU.MMUConfig.NumVirtualPages {
		return 0, errors.New("invalid address")
	}
	virtualPage, err := pt.MMU.ReadVirtualPage(uid, gid, mode, seg, page)
	if err != nil {
		return 0, err
	}
	offset := addr % PageSize
	return virtualPage[offset], nil
}

func (pt *ProcessTable) WriteAddress(
	uid int, gid int,
	mode int, seg int,
	pid int, addr int, value byte) error {
	_, ok := pt.ProcessList[pid]
	if !ok {
		return errors.New("invalid process")
	}
	page := addr / PageSize
	if page > pt.MMU.MMUConfig.NumVirtualPages {
		return errors.New("invalid address")
	}
	virtualPage, err := pt.MMU.ReadVirtualPage(uid, gid, mode, seg, page)
	if err != nil {
		return err
	}
	offset := addr % PageSize
	virtualPage[offset] = value
	return pt.MMU.WriteVirtualPage(uid, gid, mode, seg, page, virtualPage)
}
