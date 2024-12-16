package main

import (
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
	MkdirAll(path string, perm os.FileMode) error
}

type Commander interface {
	Execute(cmd string) error
}

type CompressionDetector interface {
	DetectFlag(filename string) string
}

type DefaultFileSystem struct{}

func (fs *DefaultFileSystem) FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func (fs *DefaultFileSystem) DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (fs *DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

type DefaultCommander struct{}

func (c *DefaultCommander) Execute(cmd string) error {
	return exec.Command("sh", "-c", cmd).Run()
}

type DefaultCompressionDetector struct {
	compressionTypes []Compression
}

func NewDefaultCompressionDetector(types []Compression) *DefaultCompressionDetector {
	return &DefaultCompressionDetector{compressionTypes: types}
}

func (d *DefaultCompressionDetector) DetectFlag(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check special cases first
	if flag := d.handleSpecialCases(ext); flag != "" {
		return flag
	}

	// Try extension-based detection
	if flag := d.findFlagByExtension(ext); flag != "" {
		return flag
	}

	// Try MIME type-based detection
	mimeType := mime.TypeByExtension(ext)
	return d.findFlagByMimeType(mimeType)
}

func (d *DefaultCompressionDetector) findFlagByExtension(ext string) string {
	for _, c := range d.compressionTypes {
		if containsString(c.Extensions, ext) {
			return c.Flag
		}
	}
	return ""
}

func (d *DefaultCompressionDetector) findFlagByMimeType(mimeType string) string {
	if mimeType == "" {
		return ""
	}

	for _, c := range d.compressionTypes {
		if containsString(c.MimeTypes, mimeType) {
			return c.Flag
		}
	}
	return ""
}

func (d *DefaultCompressionDetector) handleSpecialCases(ext string) string {
	switch ext {
	case ".tgz":
		return "z"
	case ".tbz2":
		return "j"
	default:
		return ""
	}
}
