/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ei-sugimoto/soybeans/internal/Err"
	"github.com/ei-sugimoto/soybeans/internal/attribute"
	"github.com/ei-sugimoto/soybeans/internal/config"
	"github.com/ei-sugimoto/soybeans/internal/mount"
	"github.com/ei-sugimoto/soybeans/internal/rootfs"
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
		if os.Getenv("REEXEC") == "true" {
			config, err := config.Load("config.json")
			if err != nil {
				util.Must(err)
			}

			util.Must(rootfs.PivotRoot(config.Root.Path))
			util.Must(mount.Mount(*config))

			log.Println("finished pivot_root and mount")
			return nil
		}
		var containerID = args[0]
		if len(args) < 1 {
			return Err.ManyArgs
		}
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

		util.Must(saveState(stateFilePath, state))
		// linuxに関する設定がない場合には、Unsahreを呼び出す
		if len(config.Linux.Namespaces) == 0 {
			log.Println("exec unshare")
			util.Must(unix.Unshare(unix.CLONE_NEWNS))
			util.Must(rootfs.PivotRoot(config.Root.Path))
			util.Must(mount.Mount(*config))

		} else {
			newCmd := exec.Command("/proc/self/exe", os.Args[1:]...)
			newCmd.Stdout = os.Stdout
			newCmd.Stderr = os.Stderr
			newCmd.Env = append(os.Environ(), "REEXEC=true")

			attribute.Attribute(newCmd, config)
			if err := newCmd.Start(); err != nil {
				log.Fatalf("Failed to start command: %v", err)
			}

			log.Printf("Running %v with pid %d\n", newCmd.Args, newCmd.Process.Pid)
			if err := newCmd.Wait(); err != nil {
				log.Fatalf("Command did not complete successfully: %v", err)
			}

		}

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
