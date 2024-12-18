package main

import (
	"os"
	"testing"
)

type MockFileSystem struct {
	files map[string]bool
	dirs  map[string]bool
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string]bool),
		dirs:  make(map[string]bool),
	}
}

func (m *MockFileSystem) FileExists(path string) (bool, error) {
	return m.files[path], nil
}

func (m *MockFileSystem) DirExists(path string) (bool, error) {
	return m.dirs[path], nil
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.dirs[path] = true
	return nil
}

func TestDetermineOperation(t *testing.T) {
	tests := []struct {
		name        string
		archivePath string
		fileExists  bool
		dirExists   bool
		want        Operation
		wantErr     bool
	}{
		{
			name:        "Archive exists - should extract",
			archivePath: "test.tar.gz",
			fileExists:  true,
			dirExists:   false,
			want:        Extract,
			wantErr:     false,
		},
		{
			name:        "Directory exists - should compress",
			archivePath: "test.tar.gz",
			fileExists:  false,
			dirExists:   true,
			want:        Compress,
			wantErr:     false,
		},
		{
			name:        "Neither exists - should error",
			archivePath: "missing.tar.gz",
			fileExists:  false,
			dirExists:   false,
			want:        Unknown,
			wantErr:     true,
		},
		{
			name:        "Both archive and directory exist - should error",
			archivePath: "test.tar.gz",
			fileExists:  true,
			dirExists:   true,
			want:        Unknown,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewMockFileSystem()
			fs.files[tt.archivePath] = tt.fileExists

			handler := NewArchiveHandler(fs, nil, nil, ".")
			dirPath := handler.getDirectoryPath(tt.archivePath)
			fs.dirs[dirPath] = tt.dirExists

			got, err := handler.determineOperation(tt.archivePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("determineOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDirectoryPath(t *testing.T) {
	tests := []struct {
		name        string
		archivePath string
		workingDir  string
		want        string
	}{
		{
			name:        "Single extension tar.gz",
			archivePath: "test.tar.gz",
			workingDir:  "/tmp",
			want:        "/tmp/test",
		},
		{
			name:        "Special case tgz",
			archivePath: "test.tgz",
			workingDir:  "/tmp",
			want:        "/tmp/test",
		},
		{
			name:        "With path components",
			archivePath: "/path/to/archive.tar.xz",
			workingDir:  "/tmp",
			want:        "/tmp/archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewArchiveHandler(nil, nil, nil, tt.workingDir)
			if got := handler.getDirectoryPath(tt.archivePath); got != tt.want {
				t.Errorf("getDirectoryPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
