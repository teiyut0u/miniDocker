package namespace

import (
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func uts(specs *specs.Spec) error {
	return syscall.Sethostname([]byte(specs.Hostname))
}
