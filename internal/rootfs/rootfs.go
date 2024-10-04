package rootfs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func PivotRoot(newRoot string) error {
	// 1. newRoot をマウントポイントにする
	if err := unix.Mount(newRoot, newRoot, "", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to bind mount new root: %v", err)
	}

	// 2. putOld ディレクトリを作成
	putOld := newRoot + "/.pivot_root"
	if err := os.MkdirAll(putOld, 0700); err != nil {
		return fmt.Errorf("failed to create putold directory: %v", err)
	}

	// 3. putOld をマウントポイントにする
	if err := unix.Mount(putOld, putOld, "", unix.MS_BIND, ""); err != nil {
		return fmt.Errorf("failed to bind mount putold directory: %v", err)
	}

	// 4. カレントディレクトリを newRoot に変更
	if err := unix.Chdir(newRoot); err != nil {
		return fmt.Errorf("failed to change directory to new root: %v", err)
	}

	// 5. pivot_root を実行
	if err := unix.PivotRoot(".", ".pivot_root"); err != nil {
		return fmt.Errorf("pivot_root failed: %v", err)
	}

	return nil
}
