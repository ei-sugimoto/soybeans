package attribute

import (
	"os/exec"

	"github.com/ei-sugimoto/soybeans/internal/config"
	"golang.org/x/sys/unix"
)

func Attribute(cmd *exec.Cmd, config *config.TConfig) {
	for _, t := range config.Linux.Namespaces {
		switch t.Type {
		case "pid":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWPID
		case "network":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWNET
		case "ipc":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWIPC
		case "uts":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWUTS
		case "user":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWUSER
		case "cgroup":
			cmd.SysProcAttr.Cloneflags |= unix.CLONE_NEWCGROUP
		}
	}
}
