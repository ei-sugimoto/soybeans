package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/ei-sugimoto/soybeans/internal"
)

func main() {

	originalHostName, err := os.Hostname()
	if err != nil {
		log.Fatalf("ホスト名の取得に失敗しました: %v", err)
	}

	cmd := exec.Command("/bin/sh")

	config := Mapping()
	cmd.SysProcAttr = Attribute(config).SysProcAttr
	if err := PivotRoot(config); err != nil {
		log.Fatalf("PivotRootの実行に失敗しました: %v", err)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("コマンドの実行に失敗しました: ", err)
		syscall.Sethostname([]byte(originalHostName))
		os.Exit(1)
	}
	syscall.Sethostname([]byte(originalHostName))

	os.Exit(0)

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
