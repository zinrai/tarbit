package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Operation int

const (
	Unknown Operation = iota
	Extract
	Compress
)

type ArchiveHandler struct {
	fs         FileSystem
	commander  Commander
	detector   CompressionDetector
	workingDir string
}

func NewArchiveHandler(fs FileSystem, commander Commander, detector CompressionDetector, workingDir string) *ArchiveHandler {
	return &ArchiveHandler{
		fs:         fs,
		commander:  commander,
		detector:   detector,
		workingDir: workingDir,
	}
}

// Determines and executes the appropriate archive operation
func (h *ArchiveHandler) Process(archivePath string) error {
	op, err := h.determineOperation(archivePath)
	if err != nil {
		return fmt.Errorf("failed to determine operation: %w", err)
	}

	switch op {
	case Extract:
		return h.extract(archivePath)
	case Compress:
		return h.compress(archivePath)
	default:
		return fmt.Errorf("unsupported operation for %s", archivePath)
	}
}

// Decides whether to extract or compress based on filesystem state
func (h *ArchiveHandler) determineOperation(archivePath string) (Operation, error) {
	exists, err := h.fs.FileExists(archivePath)
	if err != nil {
		return Unknown, err
	}

	if exists {
		return Extract, nil
	}

	// Check if a directory exists for compression
	dirPath := h.getDirectoryPath(archivePath)
	dirExists, err := h.fs.DirExists(dirPath)
	if err != nil {
		return Unknown, err
	}

	if dirExists {
		return Compress, nil
	}

	return Unknown, fmt.Errorf("neither archive nor directory exists: %s", archivePath)
}

// Handles archive extraction
func (h *ArchiveHandler) extract(archivePath string) error {
	flag := h.detector.DetectFlag(archivePath)
	if flag == "" {
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	// Obtain the absolute path of the file to be decompressed
	absArchivePath, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	targetDir := h.getDirectoryPath(archivePath)
	if err := h.fs.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	cmd := fmt.Sprintf("cd %s && tar xf%s %s", targetDir, flag, absArchivePath)
	fmt.Printf("Command: %s\n", cmd)
	if err := h.commander.Execute(cmd); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	return nil
}

// Handles directory compression
func (h *ArchiveHandler) compress(archivePath string) error {
	flag := h.detector.DetectFlag(archivePath)
	if flag == "" {
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	dirPath := h.getDirectoryPath(archivePath)
	cmd := fmt.Sprintf("tar cf%s %s %s", flag, archivePath, filepath.Base(dirPath))
	fmt.Printf("Command: %s\n", cmd)
	if err := h.commander.Execute(cmd); err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}

	return nil
}

// Returns the directory path for a given archive path
func (h *ArchiveHandler) getDirectoryPath(archivePath string) string {
	base := filepath.Base(archivePath)
	ext := filepath.Ext(base)
	dirName := strings.TrimSuffix(base, ext)
	if ext == ".gz" || ext == ".bz2" || ext == ".xz" {
		dirName = strings.TrimSuffix(dirName, ".tar")
	}
	return filepath.Join(h.workingDir, dirName)
}
