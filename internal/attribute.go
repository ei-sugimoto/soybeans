package internal

import "syscall"

type Attribute struct {
	SysProcAttr *syscall.SysProcAttr
}

func NewAttribute() *Attribute {
	return &Attribute{
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS,
		},
	}
}

func (a *Attribute) SetUID(ContainerID, HostID, Size int) {
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
