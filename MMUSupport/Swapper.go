package MMUSupport

import (
	"GolangCPUParts/MMUSupport/PhysicalMemory"
	"GolangCPUParts/RemoteLogging"
	"os"
	"strconv"
)

const SwapPageSize = 4096

// SwapperInterface
// The SwapperInterface lets us swap pages in and out of memory
type SwapperInterface struct {
	FileHandle *os.File
	Filename   string
}

func (s *SwapperInterface) Initialize() error {
	RemoteLogging.LogEvent("INFO", "Swapper Initialize", "Starting initialization")
	file, err := os.OpenFile(s.Filename, os.O_RDWR|os.O_CREATE, 0666)
	s.FileHandle = file
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper Initialize", "Initialization failed")
		return err
	}
	RemoteLogging.LogEvent("INFO", "Swapper Initialize", "Initialization completed")
	return nil
}

func (s *SwapperInterface) Terminate() error {
	RemoteLogging.LogEvent("INFO", "Swapper Terminate", "Starting termination")
	err := s.FileHandle.Close()
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper Terminate", "Termination failed: "+err.Error())
		return err
	}
	err = os.Remove(s.Filename)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper Terminate", "Termination failed: "+err.Error())
		return err
	}
	RemoteLogging.LogEvent("INFO", "Swapper Terminate", "Termination completed")
	return nil
}

func (s *SwapperInterface) SwapOut(pm PhysicalMemory.PhysicalMemory, page uint32) error {
	RemoteLogging.LogEvent("INFO", "Swapper SwapOut", "Swapping out page "+strconv.Itoa(int(page)))
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper SwapOut", "Swap failed: "+err.Error())
		return err
	}
	// Copy data from the physical page
	_, err = s.FileHandle.Write(pm.PhysicalPages[page].Buffer[:])
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper SwapOut", "Swap failed: "+err.Error())
		return err
	}
	RemoteLogging.LogEvent("INFO", "Swapper SwapOut", "Swap completed")
	return nil
}

func (s *SwapperInterface) SwapIn(pm PhysicalMemory.PhysicalMemory, page int) error {
	RemoteLogging.LogEvent("INFO", "Swapper SwapIn", "Swapping in page "+strconv.Itoa(page))
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper SwapIn", "Swap failed: "+err.Error())
		return err
	}
	_, err = s.FileHandle.Read(pm.PhysicalPages[page].Buffer[:])
	if err != nil {
		RemoteLogging.LogEvent("ERROR", "Swapper SwapIn", "Swap failed: "+err.Error())
		return err
	}
	RemoteLogging.LogEvent("INFO", "Swapper SwapIn", "Swap completed")
	return nil
}
