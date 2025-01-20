package ProcessMemory

import (
	"time"
)

type ProcessObject struct {
	Name      string
	Args      []string
	PID       uint16
	PPID      uint16
	UID       uint32
	GID       uint32
	CreatedOn time.Time
	System    bool
}

type ProcessTable struct {
	Processes map[uint16]ProcessObject
}
