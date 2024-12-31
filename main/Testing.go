package main

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/RemoteLogging"
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
	pt.CreateNewProcess("testproc", []string{}, 0, 1, 0, 1, 50)
}
