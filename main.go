package main

import (
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	filename := os.Args[1]
	flag := detectCompressionFlag(filename)

	if flag == "" {
		fmt.Printf("Unsupported file format: %s\n", filename)
		os.Exit(1)
	}

	cmd := buildTarCommand(flag, filename)
	fmt.Printf("Executing: %s\n", cmd)

	if err := executeCommand(cmd); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

// Searches for a compression flag based on the file extension
func findFlagByExtension(ext string) string {
	ext = strings.ToLower(ext)
	for _, c := range compressionTypes {
		if containsString(c.Extensions, ext) {
			return c.Flag
		}
	}
	return ""
}

// Searches for a compression flag based on the MIME type
func findFlagByMimeType(mimeType string) string {
	if mimeType == "" {
		return ""
	}

	for _, c := range compressionTypes {
		if containsString(c.MimeTypes, mimeType) {
			return c.Flag
		}
	}
	return ""
}

// Processes special file extensions that need specific handling
func handleSpecialCases(ext string) string {
	switch ext {
	case ".tgz":
		return "z"
	case ".tbz2":
		return "j"
	default:
		return ""
	}
}

// Checks if a string slice contains a specific target string
func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// Determines the appropriate compression flag for the given filename
func detectCompressionFlag(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check special cases first
	if flag := handleSpecialCases(ext); flag != "" {
		return flag
	}

	// Try extension-based detection
	if flag := findFlagByExtension(ext); flag != "" {
		return flag
	}

	// Try MIME type-based detection
	mimeType := mime.TypeByExtension(ext)
	return findFlagByMimeType(mimeType)
}

// Constructs the tar command string with the appropriate flag
func buildTarCommand(flag, filename string) string {
	return fmt.Sprintf("tar xf%s %s", flag, filename)
}

// Runs the shell command and returns any error that occurred
func executeCommand(cmd string) error {
	return exec.Command("sh", "-c", cmd).Run()
}
