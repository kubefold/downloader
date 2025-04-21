package service

import (
	"fmt"
	"github.com/kubefold/downloader/internal/types"
)

type DownloadService interface {
	Download(dataset types.Dataset, destination string, rate int) error
}

type downloadService struct {
}

func newDownloadService() DownloadService {
	return &downloadService{}
}

func (d downloadService) Download(dataset types.Dataset, destination string, rate int) error {
	url := fmt.Sprintf("https://storage.googleapis.com/alphafold-databases/v3.0/%s", string(dataset))
	_ = url
	//TODO implement me
	panic("implement me")
}
