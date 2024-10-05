/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ei-sugimoto/soybeans/internal/config"
	"github.com/ei-sugimoto/soybeans/internal/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

type ContainerState struct {
	Id        string `json:"id"`
	Pid       int    `json:"pid,omitempty"`
	Status    string `json:"status"`
	Bundle    string `json:"bundle"`
	CreatedAt string `json:"createdAt"`
	Owner     string `json:"owner"`
}

const firstProcessPid = 1

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <container-id>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Println("container id is required")
			return nil
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
			Status:    "created",
			Bundle:    cwd,
			CreatedAt: time.Now().Format(time.RFC3339),
			Owner:     hostname,
		}

		argv := append([]string{"/proc/self/exe", "init"}, config.Process.Args...)
		pid, err := syscall.ForkExec("/proc/self/exe", argv, &syscall.ProcAttr{
			Env: config.Process.Env,
			Dir: cwd,
			Sys: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
				Setsid:     true,
				UidMappings: []syscall.SysProcIDMap{
					{
						ContainerID: config.Process.User.UID,
						HostID:      os.Geteuid(),
						Size:        1,
					},
				},
				GidMappings: []syscall.SysProcIDMap{
					{
						ContainerID: config.Process.User.GID,
						HostID:      os.Getegid(),
						Size:        1,
					},
				},
				Credential: &syscall.Credential{
					Uid: 0,
					Gid: 0,
				},
				AmbientCaps: []uintptr{unix.CAP_SYS_ADMIN},
			},
			Files: []uintptr{0, 1, 2},
		})
		if err != nil {
			util.Must(err)
		}

		fmt.Printf("Started process with PID %d\n", pid)
		log.Println("argv:", argv)
		if err := syscall.Kill(pid, syscall.SIGSTOP); err != nil {
			panic(fmt.Sprintf("Failed to stop process: %v", err))
		}

		state.Pid = pid
		if err := syscall.Kill(pid, 0); err != nil {
			log.Fatalf("Process with PID %d does not exist: %v", pid, err)
		}

		util.Must(saveState(stateFilePath, state))

		// ルートファイルシステムのピボット
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

}

func saveState(path string, state *ContainerState) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create state file: %v", err)
	}
	defer file.Close()

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	byteValue, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if _, err := file.Write(byteValue); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}
