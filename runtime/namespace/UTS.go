package namespace

import (
	"fmt"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func uts(spec *specs.Spec) error {
	if err := syscall.Setdomainname([]byte(spec.Domainname)); err != nil {
		return fmt.Errorf("Failed to set domain name to be \"%s\": %v", spec.Domainname, err)
	}
	if err := syscall.Sethostname([]byte(spec.Hostname)); err != nil {
		return fmt.Errorf("Failed to set host name to be \"%s\": %v", spec.Hostname, err)
	}
	return nil
}
