package analyze

import (
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3nic1e/docker-packages-versions-exporter/docker_image"
	log "github.com/sirupsen/logrus"
)

type Analyzer struct {
	docker     *client.Client
	dockerAuth types.AuthConfig
}

func NewAnalyzer(docker *client.Client, dockerAuth types.AuthConfig) *Analyzer {
	return &Analyzer{
		docker:     docker,
		dockerAuth: dockerAuth,
	}
}

var (
	metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "package_version",
	}, []string{"image", "image_id", "os", "package"})
)

func init() {
	prometheus.MustRegister(metric)
}

func (analyzer *Analyzer) GetImagePackages(image string) {
	defer func() {
		if err := recover(); err != nil {
			log.WithField("image", image).Error("Failed to get packages list")
		}
	}()

	manager := docker_image.NewDockerImageManager(analyzer.docker, analyzer.dockerAuth, image)

	imageID, err := manager.GetImageID()
	if err != nil {
		log.WithError(err).WithField("image", image).Error("Failed to retrieve image ID")
		return
	}

	var packages []string
	var os string
	os, packages = analyzer.analyzeImage(manager, image)

	log.WithField("image", image).WithField("packages", len(packages)).Info("Got packages list")
	log.WithField("image", image).WithField("packages", packages).Debug("Got packages list")

	for _, pkg := range packages {
		metric.WithLabelValues(image, imageID, os, pkg).Set(1)
	}
}

func (analyzer *Analyzer) analyzeImage(manager *docker_image.DockerImageManager, image string) (string, []string) {
	dir := manager.GetImageContents()
	defer os.RemoveAll(dir)

	exit, err := chroot(dir)
	if err != nil {
		log.Panic(err)
	}

	os.Chdir("/")

	os := detectOs()
	log.WithField("image", image).WithField("os", os).Info("Detected OS")

	packages := getPackages(os)
	log.WithField("image", image).WithField("packages", packages).Debug("Got packages")

	// exit from the chroot
	if err := exit(); err != nil {
		log.Panic(err)
	}
	return os, packages
}
