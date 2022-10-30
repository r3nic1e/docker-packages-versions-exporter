package analyze

import (
	"github.com/r3nic1e/docker-packages-versions-exporter/analyze/packages"
	log "github.com/sirupsen/logrus"
)

func getPackages(os string) []string {
	packageGetter, err := packages.NewPackageGetter(os)
	if err != nil {
		log.WithError(err).Error("Failed to get packages getter")
		return []string{}
	}
	return packageGetter.GetPackages()
}
