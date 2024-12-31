package MMUSupport

import (
	"GolangCPUParts/RemoteLogging"
	"errors"
	"time"
)

func ProcessTableInitialize(mconf *MMUConfig) (*ProcessTable, error) {
	// Create a virtual memory object
	RemoteLogging.LogEvent("INFO", "ProcessTableInitialize", "Initialization started")
	mmu, err := VirtualMemoryInitialize(mconf)
	if err != nil {
		return nil, err
	}
	// Create the process table
	RemoteLogging.LogEvent("INFO", "ProcessTableInitialize", "Creating process table")
	pt := ProcessTable{
		MMU:         mmu,
		ProcessList: make(map[int]ProcessObject),
		NextPID:     1,
	}
	RemoteLogging.LogEvent("INFO", "ProcessTableInitialize", "Initialization complete")
	return &pt, nil
}

func ProcessTableTerminate(pt *ProcessTable) error {
	RemoteLogging.LogEvent("INFO", "ProcessTableTerminate", "Termination started")
	err := pt.MMU.VirtualMemoryTerminate()
	if err != nil {
		return err
	}
	RemoteLogging.LogEvent("INFO", "ProcessTabbleTerminate", "Termination complete")
	return nil
}

func (pt *ProcessTable) CreateNewProcess(
	name string, args []string,
	ppid int, gid int) error {
	RemoteLogging.LogEvent("INFO", Process CreateProcess", "CreateProcess start")
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
	RemoteLogging.LogEvent("INFO", "Process CreateProcess", "Create process completed")
	return nil
}

func (pt *ProcessTable) DestroyProcess(pid int) error {
	RemoteLogging.LogEvent("INFO", "Process DestroyProcess", "DestroyProcess started")
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Process DestroyProcess", "Process not found")
		return errors.New("invalid process")
	}
	err := pt.MMU.FreeBulkPages(po.Memory)
	if err != nil {
		return err
	}
	delete(pt.ProcessList, pid)
	RemoteLogging.LogEvent("INFO", "Process DestroyProcess", "Destroy process completed"")
	return nil
}

func (pt *ProcessTable) GrowSegment(pid int, gid int, prot int, seg int, newPages int) error {
	RemoteLogging.LogEvent("INFO", "Process Growpages", "Growpages started")
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Process Growpages", "Process not found")
		return errors.New("invalid process")
	}
	pagelist, err := pt.MMU.AllocateBulkPages(pid, gid, prot, seg, newPages)
	if err != nil {
		return err
	}
	po.Memory = append(po.Memory, pagelist...)
	RemoteLogging.LogEvent("INFO", "Process Growpages", "Growpages completed")
	return nil
}

func (pt *ProcessTable) GetProcessInfo(pid int) (ProcessObject, error) {
	RemoteLogging.LogEvent("INFO", "Process GetProcessInfo", "GetProcessInfo started"")
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Process GetProcessInfo", "Process not found")
		return ProcessObject{}, errors.New("invalid process")
	}
	RemoteLogging.LogEvent("INFO", "Process GetProcessInfo", "GetProcessInfo completed")
	return po, nil
}

func (pt *ProcessTable) GetProcessList() map[int]ProcessObject {
	return pt.ProcessList
}

func (pt *ProcessTable) SetProcessState(pid int, state int) error {
	RemoteLogging.LogEvent("INFO", "Process SetProcessState", "SetProcessState started")
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Process SetProcessState", "Process not found")
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
	RemoteLogging.LogEvent("INFO", "Process ReadAddress", "ReadAddress started")
	_, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent("ERROR", "Process ReadAddress", "Process not found")
		return 0, errors.New("invalid process")
	}
	page := addr / PageSize
	if page > pt.MMU.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent("ERROR", "Process ReadAddress", "Invalid address")
		return 0, errors.New("invalid address")
	}
	virtualPage, err := pt.MMU.ReadVirtualPage(uid, gid, mode, seg, page)
	if err != nil {
		return 0, err
	}
	offset := addr % PageSize
	RemoteLogging.LogEvent("INFO", "Process ReadAddress", "ReadAddress completed")
	return virtualPage[offset], nil
}

func (pt *ProcessTable) WriteAddress(
	uid int, gid int,
	mode int, seg int,
	pid int, addr int, value byte) error {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process WriteAddress",
		EventMsg:    "Starting Write Address",
	})
	_, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
			EventTime:   time.Now().Format("2006-01-02 15:04:05"),
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Process WriteAddress",
			EventMsg:    "Process not found",
		})
		return errors.New("invalid process")
	}
	page := addr / PageSize
	if page > pt.MMU.MMUConfig.NumVirtualPages {
		RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
			EventTime:   time.Now().Format("2006-01-02 15:04:05"),
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Process WriteAddress",
			EventMsg:    "Invalid address",
		})
		return errors.New("invalid address")
	}
	virtualPage, err := pt.MMU.ReadVirtualPage(uid, gid, mode, seg, page)
	if err != nil {
		return err
	}
	offset := addr % PageSize
	virtualPage[offset] = value
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process WriteAddress",
		EventMsg:    "Write Address complete",
	})
	return pt.MMU.WriteVirtualPage(uid, gid, mode, seg, page, virtualPage)
}
