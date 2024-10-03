package hostname

import (
	"github.com/ei-sugimoto/soybeans/internal/util"
	"golang.org/x/sys/unix"
)

func SetHostname(hostname string) {
	util.Must(unix.Sethostname([]byte(hostname)))
}
