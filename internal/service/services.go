package service

type Services interface {
	Download() DownloadService
}

type services struct {
	downloadService DownloadService
}

func NewServices() Services {
	return &services{
		downloadService: newDownloadService(),
	}
}

func (s services) Download() DownloadService {
	return s.downloadService
}
