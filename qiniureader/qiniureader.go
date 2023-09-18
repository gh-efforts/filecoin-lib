package qiniureader

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/service-sdk/go-sdk-qn/v2/operation"
)

type QiniuReader struct {
	Key    string
	Offset *int64
	Size   *int64
	closed bool
	body   io.ReadCloser
}

func NewQiniuReader(key string, offset *int64, size *int64) *QiniuReader {
	return &QiniuReader{
		Key:    key,
		Offset: offset,
		Size:   size,
	}
}

func (reader *QiniuReader) SeekStart() error {
	return nil
}

func (reader *QiniuReader) Seek(_ int64, _ int) (int64, error) {
	return 0, nil
}

func (reader *QiniuReader) Close() error {
	if !reader.closed {
		reader.closed = true

		if reader.body != nil {
			return reader.body.Close()
		}
	}
	return nil
}

func (reader *QiniuReader) Read(p []byte) (n int, err error) {
	if reader.closed {
		return 0, fmt.Errorf("file reader closed")
	}
	if reader.body == nil {
		var dl *operation.Downloader
		cfgPath := os.Getenv("QINIU_READER_CONFIG_PATH")
		if cfgPath != "" {
			configurable, err := operation.Load(cfgPath)
			if err != nil {
				return 0, fmt.Errorf("load qiniu config failed: %s", err)
			}
			dl = operation.NewDownloader(configurable)
		} else {
			dl = operation.NewDownloaderV2()
		}
		if reader.Offset == nil {
			resp, err := dl.DownloadRaw(reader.Key, http.Header{})
			if err != nil {
				return 0, err
			}
			if resp.StatusCode != http.StatusOK {
				return 0, fmt.Errorf(resp.Status)
			}
			reader.body = resp.Body
		} else {
			_, rc, err := dl.DownloadRangeReader(reader.Key, *reader.Offset, *reader.Size)
			if err != nil {
				return 0, err
			}
			reader.body = rc
		}
	}

	return reader.body.Read(p)
}
