/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"
	"syscall"

	"github.com/ei-sugimoto/soybeans/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]

		containerDir := filepath.Join("/var/lib/soybeans", containerID)
		state, err := loadState(containerDir)
		if err != nil {
			util.Must(err)
		}

		pid := state.Pid
		if err := syscall.Kill(pid, unix.SIGCONT); err != nil {
			util.Must(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
