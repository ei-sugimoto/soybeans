package internal

import "syscall"

func Mount(config Mapping) error {
	for _, m := range config.Mounts() {
		options := ""

		if len(m.Options) > 0 {
			options = m.Options[0]
			for _, opt := range m.Options[1:] {
				options += "," + opt
			}
		}

		if err := syscall.Mount(m.Source, m.Destination, m.Type, syscall.MS_BIND, options); err != nil {
			return err
		}
	}

	return nil
}