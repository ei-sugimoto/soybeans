package cgroup

import (
	"os/exec"

	"github.com/ei-sugimoto/soybeans/internal/config"
	"github.com/ei-sugimoto/soybeans/internal/util"
)

func SetCGroup(config config.TConfig) {
	cmd := exec.Command("cgcreate", "-g", "memory,cpu,cpuacct:")
	util.Must(cmd.Run())

	cmd = exec.Command("cgset", "-r", "memory.limit_in_bytes="+config.Linux.CgroupsPath)

}
