package main

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/ei-sugimoto/soybeans/internal"
)

func main() {

	// Mapping構造体を生成

	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = Attribute().SysProcAttr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("コマンドの実行に失敗しました: ", err)
		os.Exit(1)
	}

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

func Attribute() *internal.Attribute {
	configMap := Mapping()

	attr := internal.NewAttribute()
	attr.SetUID(configMap.UID(), os.Getuid(), 1)
	attr.SetGID(configMap.UID(), os.Getgid(), 1)
	err := attr.SetHostName(configMap.HostName())
	if err != nil {
		log.Fatalf("ホスト名の設定に失敗しました: %v", err)
	}

	return attr
}
