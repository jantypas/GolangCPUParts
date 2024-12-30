package MMU

import (
	"GolangCPUParts/RemoteLogging"
	"errors"
	"time"
)

func ProcessTableInitialize(mconf *MMUConfig) (*ProcessTable, error) {
	// Create a virtual memory object
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "ProcessTableInitialize",
		EventMsg:    "Starting initialization",
	})
	mmu, err := VirtualMemoryInitialize(mconf)
	if err != nil {
		return nil, err
	}
	// Create the process table
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "ProcessTableInitialize",
		EventMsg:    "Creating procesws table",
	})
	pt := ProcessTable{
		MMU:         mmu,
		ProcessList: make(map[int]ProcessObject),
		NextPID:     1,
	}
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "ProcessTableInitialize",
		EventMsg:    "Initialization complete",
	})
	return &pt, nil
}

func ProcessTableTerminate(pt *ProcessTable) error {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "ProcessTableTerminate",
		EventMsg:    "Starting termination",
	})
	err := pt.MMU.VirtualMemoryTerminate()
	if err != nil {
		return err
	}
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "ProcessTableTerminate",
		EventMsg:    "Termination complete",
	})
	return nil
}

func (pt *ProcessTable) CreateNewProcess(
	name string, args []string,
	ppid int, gid int) error {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process CreateProcess",
		EventMsg:    "Creation started",
	})
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
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process CreateProcess",
		EventMsg:    "Creation complete Process ID"})
	pt.NextPID++
	return nil
}

func (pt *ProcessTable) DestroyProcess(pid int) error {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process DestroyProcess",
		EventMsg:    "Starting Destroy Process",
	})
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
			EventTime:   time.Now().Format("2006-01-02 15:04:05"),
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Process DestroyProcess",
			EventMsg:    "Process not found",
		})
		return errors.New("invalid process")
	}
	err := pt.MMU.FreeBulkPages(po.Memory)
	if err != nil {
		return err
	}
	delete(pt.ProcessList, pid)
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process DestroyProcess",
		EventMsg:    "Process destroyed",
	})
	return nil
}

func (pt *ProcessTable) GrowSegment(pid int, gid int, prot int, newPages int) error {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process GrowSegment",
		EventMsg:    "Startied growing segment",
	})
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
			EventTime:   time.Now().Format("2006-01-02 15:04:05"),
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Process GrowSegment",
			EventMsg:    "Process not found",
		})
		return errors.New("invalid process")
	}
	pagelist, err := pt.MMU.AllocateBulkPages(pid, gid, prot, newPages)
	if err != nil {
		return err
	}
	po.Memory = append(po.Memory, pagelist...)
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process GrowSegment",
		EventMsg:    "Grow segment complete",
	})
	return nil
}

func (pt *ProcessTable) GetProcessInfo(pid int) (ProcessObject, error) {
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process GetProcessInfo",
		EventMsg:    "Starting Get Process Info",
	})
	po, ok := pt.ProcessList[pid]
	if !ok {
		RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
			EventTime:   time.Now().Format("2006-01-02 15:04:05"),
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Process GetProcessInfo",
			EventMsg:    "Process not found",
		})
		return ProcessObject{}, errors.New("invalid process")
	}
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   time.Now().Format("2006-01-02 15:04:05"),
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Process GetProcessInfo",
		EventMsg:    "Get Process Info complete",
	})
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
