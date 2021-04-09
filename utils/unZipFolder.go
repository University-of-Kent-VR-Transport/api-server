package utils

import (
	"archive/zip"
	"os"
	"fmt"
	"bytes"
)

func UnZipFile(zippedFolder []byte) ([]*zip.File, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zippedFolder), int64(len(zippedFolder)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create zip reader", err)
		return nil, err
	}

	return zipReader.File, nil
}