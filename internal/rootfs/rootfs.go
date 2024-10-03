package rootfs

import (
	"fmt"
	"os"
	"syscall"
)

func PivotRoot(newRoot string) error {
	// Bind mount newRoot to itself to ensure it's a mount point.
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to bind mount new root: %v", err)
	}

	// Create a directory for the old root.
	putOld := ".pivot_root"
	if err := os.MkdirAll(newRoot+"/"+putOld, 0700); err != nil {
		return fmt.Errorf("failed to create putold directory: %v", err)
	}

	// Change directory to new root.
	if err := syscall.Chdir(newRoot); err != nil {
		return fmt.Errorf("failed to change directory to new root: %v", err)
	}

	// Perform the pivot_root.
	if err := syscall.PivotRoot(".", putOld); err != nil {
		return fmt.Errorf("pivot_root failed: %v", err)
	}

	// Change the current working directory to the new root.
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change directory to new root: %v", err)
	}

	// Unmount the old root.
	if err := syscall.Unmount("/"+putOld, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("failed to unmount old root: %v", err)
	}

	// Remove the old root directory.
	if err := os.RemoveAll("/" + putOld); err != nil {
		return fmt.Errorf("failed to remove old root directory: %v", err)
	}

	return nil
}
