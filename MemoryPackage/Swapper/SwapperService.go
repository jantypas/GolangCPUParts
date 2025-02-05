package Swapper

import (
	"GolangCPUParts/MemoryPackage/PhysicalMemory"
	"errors"
	"os"
)

type SwapperContainer struct {
	Filename   string
	FileHandle *os.File
}

func Swapper_Initialize(name string) (*SwapperContainer, error) {
	sc := SwapperContainer{Filename: name}
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	sc.FileHandle = file
	return &sc, nil
}

func (sc *SwapperContainer) Terminate() {
	err := sc.FileHandle.Close()
	if err != nil {
		panic(err)
	}
}

func (sc *SwapperContainer) SwapOutPage(page uint32, buffer []byte) error {
	_, err := sc.FileHandle.Seek(int64(page*PhysicalMemory.PageSize), 0)
	if err != nil {
		return errors.New("Failed to seek to page")
	}
	_, err = sc.FileHandle.Write(buffer)
	if err != nil {
		return errors.New("Failed to write page")
	}
	return nil
}

func (sc *SwapperContainer) SwapInPage(page uint32, buffer []byte) error {
	_, err := sc.FileHandle.Seek(int64(page*PhysicalMemory.PageSize), 0)
	if err != nil {
		return errors.New("Failed to seek to page")
	}
	_, err = sc.FileHandle.Read(buffer)
	if err != nil {
		return errors.New("Failed to read page")
	}
	return nil
}
