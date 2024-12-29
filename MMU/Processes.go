package MMU

func ProcessInitialize() {
}

func ProcessCreate(
	name string, args []string,
	pid int, gid int,
	sys bool, prot int, desiredPages int ) error {
	po := ProcessObject{
		PID: pid,
		GID: gid,
		Name: name,
		Args: args,
		State: ProcessStateWaitingToRun,
		System: sys}
		err, list := ProcessAllocatePages(pid, gid, prot, desiredPages)
		if err != nil {
			return err
		}
		po.VirtualMemory = list
		return nil
	}
}

func ProcessAllocatePages(pid int, gid int, prot int, desiredPages int) (error, []int) {

}