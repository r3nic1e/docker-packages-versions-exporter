package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3nic1e/docker-packages-versions-exporter/analyze"
	log "github.com/sirupsen/logrus"
)

var dockerAuth types.AuthConfig

var (
	dockerToken = kingpin.Flag("docker-token", "Docker token").Envar("DOCKER_TOKEN").String()
	outFile     = kingpin.Flag("output-file", "Output file for metrics").Envar("OUTPUT_FILE").Required().String()
	promURL     = kingpin.Flag("prometheus-url", "Prometheus URL").Envar("PROMETHEUS_URL").Required().URL()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Panic(err)
	}

	if *dockerToken != "" {
		dockerAuth.Username = "bla"
		dockerAuth.Password = *dockerToken
	}

	images, err := getRunningImages(*promURL)
	if err != nil {
		log.Panic(err)
	}
	log.WithField("images", images).Debug("Got running images")

	analyzer := analyze.NewAnalyzer(docker, dockerAuth)
	for _, image := range images {
		analyzer.GetImagePackages(image)
	}

	err = prometheus.WriteToTextfile(*outFile, prometheus.DefaultGatherer)
	if err != nil {
		log.Panic(err)
	}
}
