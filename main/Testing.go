package main

import (
	"GolangCPUParts/MMUSupport"
	"GolangCPUParts/RemoteLogging"
)

func main() {
	RemoteLogging.LogInit("test")
	RemoteLogging.SetLogginActive(true)
	mconf := MMUSupport.MMUConfig{
		NumVirtualPages:  1024,
		NumPhysicalPages: 64,
		Swapper:          MMUSupport.SwapperInterface{Filename: "/tmp/swap.swp"},
	}
	_, err := MMUSupport.ProcessTableInitialize(&mconf)
	if err != nil {
		panic(err)
	}
}
