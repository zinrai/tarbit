package main

import (
	"fmt"
	"os"
)

type Compression struct {
	Flag       string
	MimeTypes  []string
	Extensions []string
}

var compressionTypes = []Compression{
	{
		Flag:       "z",
		MimeTypes:  []string{"application/gzip", "application/x-gzip"},
		Extensions: []string{".gz", ".tgz"},
	},
	{
		Flag:       "j",
		MimeTypes:  []string{"application/x-bzip2"},
		Extensions: []string{".bz2", ".tbz2"},
	},
	{
		Flag:       "J",
		MimeTypes:  []string{"application/x-xz"},
		Extensions: []string{".xz", ".txz"},
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tarbit <filename>")
		os.Exit(1)
	}

	// Initialize dependencies
	fs := &DefaultFileSystem{}
	commander := &DefaultCommander{}
	detector := NewDefaultCompressionDetector(compressionTypes)

	// Get current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Create archive handler
	handler := NewArchiveHandler(fs, commander, detector, workingDir)

	// Process the archive
	if err := handler.Process(os.Args[1]); err != nil {
		fmt.Printf("Error processing archive: %v\n", err)
		os.Exit(1)
	}
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
