package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PathType int

const (
	Unknown PathType = iota
	File
	Directory
)

// IsDir checks if the provided path is a directory, file, or unknown (non-existent).
func IsDir(path string) PathType {
	info, err := os.Stat(path)
	if err != nil {
		return Unknown
	}
	if os.IsNotExist(err) {
		return Unknown
	}
	if info.IsDir() {
		return Directory
	}
	return File
}

// CopyFile copies a file from the source to the destination folder.
//
//   - The destination path should be a folder. If it doesn't exist, it will be created.
//   - This function overwrites the destination file if it already exists.
func CopyFile(source, destination string) error {

	// Construct the destination file path
	destinationFilePath := filepath.Join(destination, filepath.Base(source))

	// Open the source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("CopyFile failed to open source file: %s", source)
	}
	defer sourceFile.Close()

	// Check the destination path
	destinationType := IsDir(destination)

	// If the destination path is a file, return an error
	if destinationType == File {
		return errors.New("CopyFile destination path is a file, should be a directory")
	}

	// If the destination directory doesn't exist, create it
	if destinationType == Unknown {
		if err := os.MkdirAll(destination, os.ModePerm); err != nil {
			return fmt.Errorf("CopyFile failed to create destination directory: %w", err)
		}
	}

	// Create or overwrite the destination file
	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("CopyFile failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the content from the source file to the destination file
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("CopyFile failed to copy file content: %w", err)
	}

	return nil
}

// ListEntries lists all files and directories at the first level under the given root directory.
//
// It does not check if the root path exists or is a directory.
func ListEntries(root string) ([]os.DirEntry, error) {
	var entries []os.DirEntry

	dir, err := os.Open(root)
	if err != nil {
		return nil, fmt.Errorf("ListEntries failed to open directory: %w", err)
	}
	defer dir.Close()

	entries, err = dir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("ListEntries failed to read directory contents: %w", err)
	}

	return entries, nil
}

// CopyDirectory copies a directory from the source to the destination folder.
//
//   - Both source and destination should be directory paths
//   - When looping through the entries, if an error occurs, it will continue to the next entry and return all errors at the end
func CopyDirectory(source, destination string) error {

	// Check the destination path
	destinationPathType := IsDir(destination)
	if destinationPathType == File {
		return errors.New("CopyDirectory destination path is a file, should be a directory")
	}
	if destinationPathType == Unknown {
		if err := os.MkdirAll(destination, os.ModePerm); err != nil {
			return fmt.Errorf("CopyDirectory failed to create destination directory: %w", err)
		}
	}

	// Get the list of files in the source directory
	entries, err := ListEntries(source)
	if err != nil {
		return fmt.Errorf("CopyDirectory failed to list files in source directory: %w", err)
	}

	// Loop through the entries in the source directory
	var catchErrors []error
	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(destination, entry.Name())

		// Recursively copy directories
		if entry.IsDir() {
			if err := CopyDirectory(srcPath, destPath); err != nil {
				catchErrors = append(catchErrors, fmt.Errorf("CopyDirectory failed to copy directory '%s': %w", srcPath, err))
			}
			continue
		}

		// Copy files
		if err := CopyFile(srcPath, destination); err != nil {
			catchErrors = append(catchErrors, fmt.Errorf("CopyDirectory failed to copy file '%s': %s", srcPath, err))
		}

	}

	if len(catchErrors) > 0 {
		return errors.Join(catchErrors...)
	}

	return nil
}

func Copy(source, destination string) error {

	// Check the source path
	sourcePathType := IsDir(source)
	if sourcePathType == Unknown {
		return errors.New("Copy source file does not exist")
	}

	// copy a file
	if sourcePathType == File {
		return CopyFile(source, destination)
	}

	// copy a directory
	if sourcePathType == Directory {
		destination = filepath.Join(destination, filepath.Base(source))
		return CopyDirectory(source, destination)
	}

	return nil
}
