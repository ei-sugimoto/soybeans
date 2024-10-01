package internal

import (
	"fmt"
	"os"
)

func SetupCwd(config Mapping) error {
	if err := os.Chdir(config.Cwd()); err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", config.Cwd(), err)
	}
	return nil
}
