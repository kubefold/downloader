package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/kubefold/downloader/internal/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type DownloadService interface {
	Download(dataset types.Dataset, destination string, rate int) error
}

type downloadService struct {
}

func newDownloadService() DownloadService {
	return &downloadService{}
}

const baseUrl = "https://storage.googleapis.com/alphafold-databases/v3.0/"

func (d downloadService) Download(dataset types.Dataset, destination string, rate int) error {
	destination = fmt.Sprintf("%s/%s", destination, string(dataset))
	url := baseUrl + string(dataset) + ".zst"

	if _, err := os.Stat(destination); err == nil {
		return nil
	}

	destDir := filepath.Dir(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server returned non-success status: %d %s", resp.StatusCode, resp.Status)
	}

	var reader io.Reader = resp.Body
	if rate > 0 {
		reader = newRateLimitedReader(resp.Body, rate*1024)
	}

	zReader := zstd.NewReader(reader)
	defer zReader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go d.trackDownloadProgress(ctx, &wg, destFile, dataset)

	_, err = io.Copy(destFile, zReader)
	cancel()
	wg.Wait()

	if err != nil {
		return fmt.Errorf("failed to download and decompress data: %w", err)
	}

	fileInfo, err := destFile.Stat()
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"dataset": dataset,
			"size":    fileInfo.Size(),
			"unit":    "bytes",
			"hash":    d.hashFile(destFile),
			"type":    "download",
			"total":   dataset.Size(),
		}).Info("Download completed")
	}

	return nil
}

func (d downloadService) trackDownloadProgress(ctx context.Context, wg *sync.WaitGroup, file *os.File, dataset types.Dataset) {
	defer wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fileInfo, err := file.Stat()
			if err == nil {
				logrus.WithFields(logrus.Fields{
					"dataset": dataset,
					"size":    fileInfo.Size(),
					"unit":    "bytes",
					"type":    "download",
					"total":   dataset.Size(),
				}).Info("Download progress")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (d downloadService) hashFile(file *os.File) string {
	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

type rateLimitedReader struct {
	reader  io.Reader
	limiter *rate.Limiter
	ctx     context.Context
}

func newRateLimitedReader(reader io.Reader, bytesPerSec int) io.Reader {
	limiter := rate.NewLimiter(rate.Limit(bytesPerSec), bytesPerSec)
	return &rateLimitedReader{
		reader:  reader,
		limiter: limiter,
		ctx:     context.Background(),
	}
}

func (r *rateLimitedReader) Read(p []byte) (n int, err error) {
	toRead := len(p)

	err = r.limiter.WaitN(r.ctx, toRead)
	if err != nil {
		return 0, err
	}

	return r.reader.Read(p)
}
