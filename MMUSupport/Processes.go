package MMUSupport

import (
	"errors"
	"time"
)

// Proc3essObject -- for every system procecss, one of these structures must exist
type ProcessObject struct {
	InUse           bool
	Name            string    // The application name
	Args            []string  // Application arguments
	PID             uint      // Process ID
	PPID            uint      // Parent process ID
	UID             uint      // User ID
	GID             uint      // Group ID
	System          bool      // System privileges active
	State           uint      // Process state flags
	CreatedOn       time.Time // Creation date/time
	ProcessContext  interface{}
	SegmentIdx      uint
	SegmentRequests []SegmentRequest
	Segments        []Segment // Segments that make up this process
}

// ProcessTable -- all process structures above are kept in the process table
type ProcessTable struct {
	ProcessList []ProcessObject // The table of processes
	NextPID     uint            // Next process ID
	MMU         MMUStruct       // The MMU for all processes
}

type SegmentRequest struct {
	Flags        uint
	AddressBase  uint
	AddressLimit uint
	Perm         uint
	NumPages     uint
}
type Segment struct {
	Flags         uint
	AddressBase   uint
	AddressLimit  uint
	NumPages      uint
	Perm          uint
	VirtualMemory []uint
}

func ProcessTableInitialize(mmu *MMUStruct) (*ProcessTable, error) {
	pt := ProcessTable{
		ProcessList: make([]ProcessObject, 0),
		NextPID:     0,
		MMU:         *mmu,
	}
	return &pt, nil
}

func (pt *ProcessTable) ProcessTableTerminate() error {
	for k := 0; k < len(pt.ProcessList); k++ {
		if pt.ProcessList[k].InUse {
			err := pt.ProcessDestroyProcess(uint(k))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (pt *ProcessTable) ProcessTableList() []ProcessObject {
	return pt.ProcessList
}

func (pt *ProcessTable) ProcessCreate(p ProcessObject) error {
	p.PID = pt.NextPID
	p.State = ProcessStateWaitingToRun
	p.CreatedOn = time.Now()
	p.Segments = make([]Segment, len(p.SegmentRequests))
	for _, sr := range p.SegmentRequests {
		seg := Segment{
			Flags:         sr.Flags,
			AddressBase:   sr.AddressBase,
			AddressLimit:  sr.AddressLimit,
			NumPages:      sr.NumPages,
			Perm:          sr.Perm,
			VirtualMemory: make([]uint, sr.NumPages),
		}
		pages, err := pt.MMU.AllocateBulkPages(sr.NumPages)
		if err != nil {
			return err
		}
		seg.VirtualMemory = pages
		p.Segments[p.SegmentIdx] = seg
		p.SegmentIdx++
	}
	pt.ProcessList[pt.NextPID].InUse = true
	pt.ProcessList[pt.NextPID] = p
	pt.NextPID++
	return nil
}

func (pt *ProcessTable) ProcessFindProcess(pid uint) (uint, error) {
	for idx := 0; idx < len(pt.ProcessList); idx++ {
		if pt.ProcessList[idx].InUse && pt.ProcessList[idx].PID == pid {
			return uint(idx), nil
		}
	}
	return 0, errors.New("invalid process")

}

func (pt *ProcessTable) ProcessDestroyProcess(pid uint) error {
	idx, err := pt.ProcessFindProcess(pid)
	if err != nil {
		return errors.New("invalid process")
	}
	for _, seg := range pt.ProcessList[idx].Segments {
		err := pt.MMU.FreeBulkPages(seg.VirtualMemory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pt *ProcessTable) ProcessGrowSegment(pid uint, seg uint, numPages uint) error {
	idx, err := pt.ProcessFindProcess(pid)
	if err != nil {
		return errors.New("invalid process")
	}
	po := pt.ProcessList[idx]
	if seg > uint(len(po.Segments)) {
		return errors.New("invalid segment")
	}
	newPages, err := pt.MMU.AllocateBulkPages(numPages)
	pt.ProcessList[idx].Segments[seg].VirtualMemory = append(pt.ProcessList[idx].Segments[seg].VirtualMemory, newPages...)
	return nil
}

func (pt *ProcessTable) ProcessReadPage(pid uint, page uint) ([]byte, error) {

}
func (pt *ProcessTable) ProcessWritePage(pid uint, page uint, value []byte) error {

}
func (pt *ProcessTable) ProcessReadByte(pid uint, addr uint) (byte, error) {

}
func (pt *ProcessTable) ProcessWriteByte(pid uint, addr uint, value byte) error {

}
func (pt *ProcessTable) ProcessReadWord(pid uint, addr uint) (uint16, error) {

}
func (pt *ProcessTable) ProcessWriteWord(pd uint, addr uint, value uint16) error {

}
func (pt *ProcessTable) ProcessReadDouble(pid uint, addr uint) (uint32, error) {

}
func (pt *ProcessTable) ProcessWriteDouble(pid uint, addr uint, value uint32) error {

}
func (pt *ProcessTable) ProcessReadQuad(pid uint, addr uint) (uint64, error) {

}
func (pt *ProcessTable) ProcessWriteQuad(pid uint, addr uint, value uint64) error {

}
