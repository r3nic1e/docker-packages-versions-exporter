package packages

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AlpinePackageGetter struct{}

func init() {
	registerPackageGetter(&AlpinePackageGetter{}, "alpine")
}

func (*AlpinePackageGetter) GetPackages() (pkgs []string) {
	installed, err := ioutil.ReadFile("/lib/apk/db/installed")
	if err != nil {
		log.WithError(err).Error("Failed to get installed apk packages")
		return
	}

	packages := strings.Split(string(installed), "\n\n")
	for _, pkgLines := range packages {
		var arch, name, version string
		for _, line := range strings.Split(pkgLines, "\n") {
			switch {
			case strings.HasPrefix(line, "A:"):
				arch = strings.TrimPrefix(line, "A:")
			case strings.HasPrefix(line, "P:"):
				name = strings.TrimPrefix(line, "P:")
			case strings.HasPrefix(line, "V:"):
				version = strings.TrimPrefix(line, "V:")
			}
		}

		if name != "" {
			pkgs = append(pkgs, fmt.Sprintf("%s %s %s", name, version, arch))
		}
	}

	return

}
