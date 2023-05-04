package qiniureader

import (
	"fmt"
	"io"
	"net/http"

	"github.com/service-sdk/go-sdk-qn/syncdata/operation"
)

type QiniuReader struct {
	Key    string
	closed bool
	body   io.ReadCloser
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
		resp, err := operation.NewDownloaderV2().DownloadRaw(reader.Key, http.Header{})
		if err != nil {
			return 0, err
		}

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf(resp.Status)
		}

		reader.body = resp.Body
	}

	return reader.body.Read(p)
}
