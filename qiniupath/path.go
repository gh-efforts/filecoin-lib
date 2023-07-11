package qiniupath

import (
	"net/url"
	"strings"
)

const (
	// QiniuProtocol is the protocol of Qiniu.
	// A file path should start with "qiniu://" or "qiniu:".
	// 1. qiniu://{host}/file/path
	// 2. qiniu:/file/path
	// Check https://en.wikipedia.org/wiki/File_URI_scheme
	QiniuProtocol = "qiniu"
)

func SplitQiniuPath(s string, trimPrefix bool) string {
	u, err := url.Parse(s)
	if err != nil {
		// todo: should panic or handle error?
		return s
	}
	if u.Scheme != QiniuProtocol {
		return s
	}
	if trimPrefix {
		return strings.TrimPrefix(u.Path, "/")
	}
	return u.Path
}

func IsQiniuPath(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		// todo: should panic or handle error?
		return false
	}
	return u.Scheme == QiniuProtocol
}
