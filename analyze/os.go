package analyze

import (
	"os"

	"github.com/cobaugh/osrelease"
)

func detectOs() string {
	if _, err := os.Stat("/etc/alpine-release"); err == nil {
		return "alpine"
	}

	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}

	if _, err := os.Stat("/etc/centos-release"); err == nil {
		return "centos"
	}

	if _, err := os.Stat("/var/lib"); os.IsNotExist(err) {
		return "scratch"
	}

	osrelease, err := osrelease.Read()
	if err != nil {
		return ""
	}

	return osrelease["ID"]
}
