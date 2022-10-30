package packages

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tadasv/go-dpkg"
)

type DebianPackageGetter struct{}

func init() {
	registerPackageGetter(&DebianPackageGetter{}, "debian")
}

func (*DebianPackageGetter) GetPackages() (pkgs []string) {
	packages, err := dpkg.ReadPackagesFromFile("/var/lib/dpkg/status")
	if err != nil {
		log.WithError(err).Error("Failed to read dpkg status file")
		return
	}

	for _, pkg := range packages {
		if pkg.Status != "install ok installed" {
			continue
		}

		pkgs = append(pkgs, fmt.Sprintf("%s %s %s", pkg.Package, pkg.Version, pkg.Architecture))
	}

	return
}
