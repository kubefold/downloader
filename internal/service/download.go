package service

import (
	"archive/tar"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/kubefold/downloader/pkg/types"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/zstd"
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
	datasetStr := string(dataset)
	isTar := strings.HasSuffix(datasetStr, ".tar")

	finalDestination := destination
	if isTar {
		finalDestination = destination
	} else {
		finalDestination = filepath.Join(destination, datasetStr)
	}

	url := baseUrl + datasetStr + ".zst"

	if isTar {
		if dirExists, err := d.directoryExists(finalDestination); err == nil && dirExists {
			dirSize, err := d.calculateDirSize(finalDestination)
			if err == nil && dirSize == dataset.Size() {
				return nil
			}
		}
	} else {
		if fileInfo, err := os.Stat(finalDestination); err == nil {
			if fileInfo.Size() == dataset.Size() {
				return nil
			}
			if err := os.Remove(finalDestination); err != nil {
				return fmt.Errorf("failed to remove existing file with incorrect size: %w", err)
			}
		}
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

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
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	if isTar {
		go d.trackDirProgress(ctx, &wg, finalDestination, dataset)

		tarReader := tar.NewReader(zReader)
		if err := d.extractTar(tarReader, finalDestination); err != nil {
			return fmt.Errorf("failed to extract tar archive: %w", err)
		}

		dirSize, err := d.calculateDirSize(finalDestination)
		if err == nil {
			logrus.WithFields(logrus.Fields{
				"dataset": dataset,
				"size":    dirSize,
				"unit":    "bytes",
				"type":    "download",
				"total":   dataset.Size(),
			}).Info("Tar extraction completed")
		}
	} else {
		destFile, err := os.Create(finalDestination)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destFile.Close()

		go d.trackFileProgress(ctx, &wg, destFile, dataset)

		_, err = io.Copy(destFile, zReader)
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
	}

	return nil
}

func (d downloadService) extractTar(tarReader *tar.Reader, destination string) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return err
			}

			file.Close()
		}
	}

	return nil
}

func (d downloadService) directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func (d downloadService) calculateDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func (d downloadService) trackDirProgress(ctx context.Context, wg *sync.WaitGroup, dirPath string, dataset types.Dataset) {
	defer wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			size, err := d.calculateDirSize(dirPath)
			if err == nil {
				logrus.WithFields(logrus.Fields{
					"dataset": dataset,
					"size":    size,
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

func (d downloadService) trackFileProgress(ctx context.Context, wg *sync.WaitGroup, file *os.File, dataset types.Dataset) {
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
