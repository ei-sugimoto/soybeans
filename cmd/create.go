/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ei-sugimoto/soybeans/internal/v2_pkg/Err"
	"github.com/ei-sugimoto/soybeans/internal/v2_pkg/config"
	"github.com/spf13/cobra"
)

type ContainerState struct {
	Id        string `json:"id"`
	Pid       int    `json:"pid,omitempty"`
	Status    string `json:"status"`
	Bundle    string `json:"bundle"`
	CreatedAt string `json:"createdAt"`
	Owner     string `json:"owner"`
}

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
		var containerID = args[0]
		if len(args) < 1 {
			return Err.ManyArgs
		}
		containerDir := filepath.Join("/var/lib/soybeans", containerID)
		if err := os.MkdirAll(containerDir, 0755); err != nil {
			return fmt.Errorf("failed to create container directory: %v", err)
		}

		config, err := config.Load("config.json")
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		stateFilePath := filepath.Join(containerDir, "state.json")
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}

		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %v", err)
		}

		state := &ContainerState{
			Id:        containerID,
			Pid:       0,
			Status:    "created",
			Bundle:    cwd,
			CreatedAt: time.Now().Format(time.RFC3339),
			Owner:     hostname,
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

}
