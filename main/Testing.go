package main

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/RemoteLogging"
	"fmt"
)

func main() {
	RemoteLogging.LogInit("test")
	RemoteLogging.SetLogginActive(false)
	mconf := MMUSupport.MMUConfig{
		NumVirtualPages:  1024,
		NumPhysicalPages: 64,
		Swapper:          MMUSupport.SwapperInterface{Filename: "/tmp/swap.swp"},
	}
	pt, err := MMUSupport.ProcessTableInitialize(&mconf)
	if err != nil {
		panic(err)
	}
	pt.CreateNewProcess("testproc", []string{}, 0, 1, 0, 1, 20)
	pt.WriteAddress(0, 1, 0, 1, 1, 500, 13)
	v, err := pt.ReadAddress(0, 1, 0, 1, 1, 500)
	fmt.Println(v)
	pt.DestroyProcess(1)
}
