package main

import (
	"github.com/kubefold/downloader/pkg/types"
	"os"
	"strconv"

	"github.com/kubefold/downloader/internal/service"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	dataset := os.Getenv("DATASET")
	if dataset == "" {
		var possibleDatasets string
		for _, dataset := range types.Datasets {
			possibleDatasets += " " + string(dataset)
		}
		logrus.Fatalf("dataset %s is invalid. Choose one of possible values: %s", dataset, possibleDatasets)
	}

	destination := os.Getenv("DESTINATION")
	if destination == "" {
		logrus.Fatalf("destination %s is invalid", destination)
	}

	rate := os.Getenv("RATE")
	if rate == "" {
		rate = "0"
	}
	parsedRate, err := strconv.ParseInt(rate, 10, 64)
	if err != nil {
		logrus.Fatalf("rate %s is invalid", rate)
	}

	services := service.NewServices()
	err = services.Download().Download(types.Dataset(dataset), destination, int(parsedRate))
	if err != nil {
		logrus.Fatalf("failed to download dataset: %s", dataset)
	}
}
