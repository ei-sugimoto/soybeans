package config

type TConfig struct {
	Version string `json:"ociVersion"`
	Process struct {
		Terminal bool `json:"terminal"`
		User     struct {
			UID            int   `json:"uid"`
			GID            int   `json:"gid"`
			AdditionalGids []int `json:"additionalGids"`
		} `json:"user"`
		Args []string `json:"args"`
		Env  []string `json:"env"`
		Cwd  string   `json:"cwd"`
	} `json:"process"`
	Root struct {
		Path     string `json:"path"`
		Readonly bool   `json:"readonly"`
	} `json:"root"`
	Hostname string `json:"hostname"`
	Mounts   []struct {
		Destination string   `json:"destination"`
		Type        string   `json:"type"`
		Source      string   `json:"source"`
		Options     []string `json:"options"`
	} `json:"mounts"`
	Linux struct {
		Namespaces []struct {
			Type string `json:"type"`
		} `json:"namespaces"`
		Resources struct {
			Memory struct {
				Limit int64 `json:"limit"`
			} `json:"memory"`
			CPU struct {
				Shares int `json:"shares"`
			} `json:"cpu"`
		} `json:"resources"`
		CgroupsPath     string `json:"cgroupsPath"`
		ApparmorProfile string `json:"apparmorProfile"`
		Seccomp         struct {
			DefaultAction string `json:"defaultAction"`
			Unixs         []struct {
				Names  []string `json:"names"`
				Action string   `json:"action"`
			} `json:"unixs"`
		} `json:"seccomp"`
	} `json:"linux"`
	Hooks struct {
		Prestart []struct {
			Path string   `json:"path"`
			Args []string `json:"args"`
			Env  []string `json:"env"`
		} `json:"prestart"`
		Poststart []struct {
			Path string   `json:"path"`
			Args []string `json:"args"`
			Env  []string `json:"env"`
		} `json:"poststart"`
		Poststop []struct {
			Path string   `json:"path"`
			Args []string `json:"args"`
			Env  []string `json:"env"`
		} `json:"poststop"`
	} `json:"hooks"`
}
