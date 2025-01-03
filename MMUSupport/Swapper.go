package MMUSupport

import (
	"GolangCPUParts/RemoteLogging"
	"os"
)

const SwapPageSize = 4096

// SwapperInterface
// The SwapperInterface lets us swap pages in and out of memory
type SwapperInterface struct {
	ServingMMU *MMUStruct
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
	err := s.FileHandle.Close()
	if err != nil {
		return err
	}
	err = os.Remove(s.Filename)
	if err != nil {
		return err
	}
	return nil
}

func (s *SwapperInterface) SwapOut(page int) error {
	var buffer [SwapPageSize]byte
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		return err
	}
	// Copy data from the phyical page
	copy(buffer[:], s.ServingMMU.PhysicalMem[page*SwapPageSize:page*SwapPageSize+SwapPageSize])
	_, err = s.FileHandle.Write(buffer[:])
	if err != nil {
		return err
	}
	return nil
}

func (s *SwapperInterface) SwapIn(page int) error {
	var buffer [SwapPageSize]byte
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Read(buffer[:])
	if err != nil {
		return err
	}
	copy(s.ServingMMU.PhysicalMem[page*SwapPageSize:page*SwapPageSize+SwapPageSize], buffer[:])
	return nil
}
