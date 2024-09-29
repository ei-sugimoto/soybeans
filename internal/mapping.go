package internal

import (
	jsoniter "github.com/json-iterator/go"
)

type OCIConfig struct {
	OciVersion string `json:"ociVersion"`
	Process    struct {
		Terminal bool `json:"terminal"`
		User     struct {
			Uid int `json:"uid"`
			Gid int `json:"gid"`
		} `json:"user"`
		Args         []string `json:"args"`
		Env          []string `json:"env"`
		Cwd          string   `json:"cwd"`
		Capabilities struct {
			Bounding    []string `json:"bounding"`
			Effective   []string `json:"effective"`
			Inheritable []string `json:"inheritable"`
			Permitted   []string `json:"permitted"`
			Ambient     []string `json:"ambient"`
		} `json:"capabilities"`
		NoNewPrivileges bool `json:"noNewPrivileges"`
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
			Syscalls      []struct {
				Names  []string `json:"names"`
				Action string   `json:"action"`
			} `json:"syscalls"`
		} `json:"seccomp"`
	} `json:"linux"`
}

type Mapping interface {
	Unmarshal([]byte) error
	Version() string
	UID() int
	HostName() string
	NameSpaces() []string
	RootPath() string
	Args() []string
	Env() []string
}

type MappingImpl struct {
	Config OCIConfig
}

func NewMapping() Mapping {
	return &MappingImpl{}
}

func (m *MappingImpl) Unmarshal(data []byte) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &m.Config); err != nil {
		return err
	}
	return nil
}

func (m *MappingImpl) Version() string {
	return m.Config.OciVersion
}

func (m *MappingImpl) UID() int {
	return m.Config.Process.User.Uid
}

func (m *MappingImpl) HostName() string {
	return m.Config.Hostname
}

func (m *MappingImpl) NameSpaces() []string {
	var namespaces []string
	for _, ns := range m.Config.Linux.Namespaces {
		namespaces = append(namespaces, ns.Type)
	}
	return namespaces
}

func (m *MappingImpl) RootPath() string {
	return m.Config.Root.Path
}

func (m *MappingImpl) Args() []string {
	return m.Config.Process.Args
}

func (m *MappingImpl) Env() []string {
	return m.Config.Process.Env
}
