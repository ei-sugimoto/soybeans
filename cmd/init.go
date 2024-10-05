/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/ei-sugimoto/soybeans/internal/config"
	"github.com/ei-sugimoto/soybeans/internal/mount"
	"github.com/ei-sugimoto/soybeans/internal/rootfs"
	"github.com/ei-sugimoto/soybeans/internal/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		config, err := config.Load("config.json")
		if err != nil {
			util.Must(err)
		}

		util.Must(unix.Unshare(unix.CLONE_NEWPID | unix.CLONE_NEWNS | unix.CLONE_NEWNET | unix.CLONE_NEWUTS | unix.CLONE_NEWIPC))

		// ルートファイルシステムのピボット
		util.Must(rootfs.PivotRoot(config.Root.Path))

		// 必要なファイルシステムのマウント
		util.Must(mount.Mount(*config))

		if err := unix.Sethostname([]byte(config.Hostname)); err != nil {
			panic(err)
		}

		if err := unix.Setgid(config.Process.User.GID); err != nil {
			panic(err)
		}
		if err := unix.Setuid(config.Process.User.UID); err != nil {
			panic(err)
		}

		if err := unix.Chdir(config.Process.Cwd); err != nil {
			panic(err)
		}

		if err := unix.Exec(args[0], args[1:], os.Environ()); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}

func loadState(containerDir string) (*ContainerState, error) {
	stateFilePath := filepath.Join(containerDir, "state.json")
	state := &ContainerState{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	file, err := os.Open(stateFilePath)
	if err != nil {
		return nil, errors.New("failed to open state file")
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("failed to read state file")
	}

	if err := json.Unmarshal(byteValue, &state); err != nil {
		return nil, errors.New("failed to unmarshal state file")
	}

	return state, nil
}
