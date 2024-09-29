package internal_test

import (
	"testing"

	"github.com/ei-sugimoto/soybeans/internal"
)

var mapping internal.Mapping

func init() {
	// テストの前処理
	mapping = internal.NewMapping()

}

func TestMapping(t *testing.T) {
	tc := []byte(`{
  "ociVersion": "1.0.2",
  "process": {
    "terminal": false,
    "user": {
      "uid": 0,
      "gid": 0
    },
    "args": [
      "sh"
    ],
    "env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "cwd": "/",
    "capabilities": {
      "bounding": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "effective": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "inheritable": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "permitted": ["CAP_CHOWN", "CAP_DAC_OVERRIDE"],
      "ambient": []
    },
    "noNewPrivileges": true
  },
  "root": {
    "path": "rootfs",
    "readonly": true
  },
  "hostname": "my-container",
  "mounts": [
    {
      "destination": "/proc",
      "type": "proc",
      "source": "proc",
      "options": ["nosuid", "noexec", "nodev"]
    },
    {
      "destination": "/dev",
      "type": "tmpfs",
      "source": "tmpfs",
      "options": ["nosuid", "strictatime", "mode=755", "size=65536k"]
    }
  ],
  "linux": {
    "namespaces": [
      {
        "type": "pid"
      },
      {
        "type": "network"
      },
      {
        "type": "ipc"
      },
      {
        "type": "uts"
      },
      {
        "type": "mount"
      }
    ],
    "resources": {
      "memory": {
        "limit": 536870912
      },
      "cpu": {
        "shares": 1024
      }
    },
    "cgroupsPath": "/my-container",
    "apparmorProfile": "docker-default",
    "seccomp": {
      "defaultAction": "SCMP_ACT_ERRNO",
      "syscalls": [
        {
          "names": ["clone", "execve", "exit"],
          "action": "SCMP_ACT_ALLOW"
        }
      ]
    }
  }
}
`)
	if err := mapping.Unmarshal(tc); err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if got, want := mapping.Version(), "1.0.2"; got != want {
		t.Errorf("Version() = %q; want %q", got, want)
	}
}
