package internal

import "syscall"

type Attribute struct {
	SysProcAttr *syscall.SysProcAttr
}

func NewAttribute() *Attribute {
	return &Attribute{
		SysProcAttr: &syscall.SysProcAttr{},
	}
}

func (a *Attribute) SetFlag(str []string) {
	for _, s := range str {
		switch s {
		case "mount":
			a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
		case "pid":
			a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
		case "uts":
			a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUTS
		case "ipc":
			a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWIPC
		case "network":
			a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
		}
	}
}

func (a *Attribute) SetUID(ContainerID, HostID, Size int) {

	if a.SysProcAttr.Cloneflags&syscall.CLONE_NEWUSER == 0 {
		a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
	}
	newUidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: ContainerID,
			HostID:      HostID,
			Size:        Size,
		},
	}

	a.SysProcAttr.UidMappings = newUidMappings
}

func (a *Attribute) SetGID(ContainerID, HostID, Size int) {
	if a.SysProcAttr.Cloneflags&syscall.CLONE_NEWUSER == 0 {
		a.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
	}

	newGidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: ContainerID,
			HostID:      HostID,
			Size:        Size,
		},
	}

	a.SysProcAttr.GidMappings = newGidMappings
}

func (a *Attribute) SetHostName(hostname string) error {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return err
	}

	return nil
}
