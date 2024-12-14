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
func (s *SwapperInterface) SwapOut(pos int64, buffer []byte) error {
	_, err := s.FileHandle.Seek(pos, 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
func (s *SwapperInterface) SwapIn(pos int64, buffer []byte) error {
	_, err := s.FileHandle.Seek(pos, 0)
	if err != nil {
		return err
	}
	_, err = s.FileHandle.Read(buffer)
	if err != nil {
	}
	return nil
}
