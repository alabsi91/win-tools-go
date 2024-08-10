package app

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type utils struct{}

var Utils = &utils{}

func (utils) IsPathExists(str string) bool {
	path := strings.Trim(str, "\r\n")
	_, err := os.Stat(path)
	return err == nil
}

func (utils) CopyFile(srcPath string, dstPath string) error {
	// Open the source file for reading
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file for writing
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
