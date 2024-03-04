package utils

import (
	"fmt"
	"os"
)

func CustomPath(addPath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	return dir + addPath, nil
}
