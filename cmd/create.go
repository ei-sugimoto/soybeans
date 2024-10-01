/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/ei-sugimoto/soybeans/internal"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {

			return internal.ErrManyArgs
		}

		MapConfig := Mapping()

		ContainerCMD := exec.Command(MapConfig.Args()[0], MapConfig.Args()[1:]...)

		if MapConfig.Terminal() {
			ContainerCMD.Stdin = os.Stdin
			ContainerCMD.Stdout = os.Stdout
			ContainerCMD.Stderr = os.Stderr
		}

		originalHostname, err := os.Hostname()
		if err != nil {
			return err
		}

		err = syscall.Sethostname([]byte(MapConfig.HostName()))
		if err != nil {
			return err
		}

		ContainerCMD.SysProcAttr = Attribute(MapConfig).SysProcAttr

		if err := PivotRoot(MapConfig); err != nil {
			return err
		}

		ContainerCMD.Env = append(os.Environ(), MapConfig.Env()...)

		if err := internal.SetupCwd(MapConfig); err != nil {
			return err
		}

		if err := syscall.Sethostname([]byte(originalHostname)); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

}

func Mapping() internal.Mapping {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	// ファイルの内容を読み取る
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("ファイルを読み取れませんでした: %v", err)
	}

	configMap := internal.NewMapping()
	if err := configMap.Unmarshal(byteValue); err != nil {
		log.Fatalf("JSONのパースに失敗しました: %v", err)
	}

	return configMap
}

func Attribute(config internal.Mapping) *internal.Attribute {

	attr := internal.NewAttribute()

	attr.SetFlag(config.NameSpaces())

	attr.SetUID(config.UID(), os.Getuid(), 1)
	attr.SetGID(config.UID(), os.Getgid(), 1)
	err := attr.SetHostName(config.HostName())
	if err != nil {
		log.Fatalf("ホスト名の設定に失敗しました: %v", err)
	}

	return attr
}

func PivotRoot(config internal.Mapping) error {
	if err := internal.PivotRoot(config.RootPath()); err != nil {
		return err
	}

	return nil
}
