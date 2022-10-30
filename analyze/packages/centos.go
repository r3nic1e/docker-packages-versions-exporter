package packages

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CentOSPackageGetter struct{}

func init() {
	registerPackageGetter(&CentOSPackageGetter{}, "centos")
}

func (*CentOSPackageGetter) GetPackages() (pkgs []string) {
	out, err := exec.Command("rpm", "-qa").Output()
	if err != nil {
		log.WithError(err).Error("Failed to run rpm command")
		return
	}

	pkgs = strings.Split(strings.TrimSpace(string(out)), "\n")
	return
}
