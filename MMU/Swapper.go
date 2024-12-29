package MMU

import "os"

// SwapperInterface
// The SwapperInterface lets us swap pages in and out of memory
type SwapperInterface struct {
	FileHandle *os.File
	Filename   string
}

func (s *SwapperInterface) Initialize() error {
	file, err := os.OpenFile(s.Filename, os.O_RDWR|os.O_CREATE, 0666)
	s.FileHandle = file
	if err != nil {
		return err
	}
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
	_, err := s.FileHandle.Seek(int64(page*PageSize), 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Write(buffer[:PageSize])
	if err != nil {
		return err
	}
	return nil
}

func (s *SwapperInterface) SwapIn(page int, buffer []byte) error {
	_, err := s.FileHandle.Seek(int64(page*PageSize), 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Read(buffer[:PageSize])
	if err != nil {
		return err
	}
	return nil
}
