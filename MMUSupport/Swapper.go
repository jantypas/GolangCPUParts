package MMUSupport

import (
	"GolangCPUParts/RemoteLogging"
	"os"
)

const SwapPageSize = 4096

func (s *SwapperInterface) Initialize() error {
	RemoteLogging.LogEvent(
		RemoteLogging.LogEventStruct{
			EventTime:   "2006-01-02 15:04:05",
			EventApp:    "",
			EventLevel:  "INFO",
			EventSource: "Swapper Initialize",
			EventMsg:    "Starting initialization",
		})
	file, err := os.OpenFile(s.Filename, os.O_RDWR|os.O_CREATE, 0666)
	s.FileHandle = file
	if err != nil {
		RemoteLogging.LogEvent(
			RemoteLogging.LogEventStruct{
				EventTime:   "2006-01-02 15:04:05",
				EventApp:    "",
				EventLevel:  "INFO",
				EventSource: "Swapper Initialize",
				EventMsg:    "Error opening file",
			})
		return err
	}
	RemoteLogging.LogEvent(RemoteLogging.LogEventStruct{
		EventTime:   "2006-01-02 15:04:05",
		EventApp:    "",
		EventLevel:  "INFO",
		EventSource: "Swapper Initialize",
		EventMsg:    "Initialization complete",
	})
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

func (s *SwapperInterface) SwapOut(page int, buffer []byte) error {
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Write(buffer[:SwapPageSize])
	if err != nil {
		return err
	}
	return nil
}

func (s *SwapperInterface) SwapIn(page int, buffer []byte) error {
	_, err := s.FileHandle.Seek(int64(page*SwapPageSize), 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Read(buffer[:SwapPageSize])
	if err != nil {
		return err
	}
	return nil
}
