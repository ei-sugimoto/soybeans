/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ei-sugimoto/soybeans/internal/config"
	"github.com/ei-sugimoto/soybeans/internal/mount"
	"github.com/ei-sugimoto/soybeans/internal/rootfs"
	"github.com/ei-sugimoto/soybeans/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// stateCmd represents the state command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run <container-id>",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("container id is required")
			return
		}
		containerID := args[0]

		containerDir := filepath.Join("/var/lib/soybeans", containerID)

		util.Must(os.MkdirAll(containerDir, 0755))

		config, err := config.Load("config.json")
		if err != nil {
			util.Must(err)
		}

		stateFilePath := filepath.Join(containerDir, "state.json")
		cwd, err := os.Getwd()
		if err != nil {
			util.Must(err)
		}

		hostname, err := os.Hostname()
		if err != nil {
			util.Must(err)
		}

		state := &ContainerState{
			Id:        containerID,
			Pid:       os.Getpid(),
			Status:    "running",
			Bundle:    cwd,
			CreatedAt: time.Now().Format(time.RFC3339),
			Owner:     hostname,
		}

		util.Must(saveState(stateFilePath, state))

		util.Must(unix.Unshare(unix.CLONE_NEWPID | unix.CLONE_NEWNS | unix.CLONE_NEWNET | unix.CLONE_NEWUTS | unix.CLONE_NEWIPC))

		// ルートファイルシステムのピボット
		util.Must(rootfs.PivotRoot(config.Root.Path))

		// 必要なファイルシステムのマウント
		util.Must(mount.Mount(*config))

		if err := unix.Sethostname([]byte(config.Hostname)); err != nil {
			panic(err)
		}
		binary, err := exec.LookPath(config.Process.Args[0])
		if err != nil {
			panic(fmt.Sprintf("バイナリが見つかりません: %v", err))
		}

		if err := unix.Exec(binary, config.Process.Args, config.Process.Env); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

}
