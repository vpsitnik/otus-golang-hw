package main

import (
	"errors"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// open file for read
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Panicf("File %v is not exist\n", fromPath)
		}
		log.Panic(err)
		return ErrUnsupportedFile
	}
	defer fileFrom.Close()

	// create file for copy
	fileTo, err := os.Create(toPath)
	if err != nil {
		log.Panicf("Failed to create file %v with error: %v\n", toPath, err)
		return ErrUnsupportedFile
	}
	defer fileTo.Close()

	// get FileInfo
	fi, err := fileFrom.Stat()
	if err != nil {
		log.Panic(err)
		return ErrUnsupportedFile
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// get file size
	size := fi.Size()

	// compare offset and file size
	if offset > size {
		log.Panicf("Offset value: %v is bigger than file size: %v\n", offset, size)
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit+offset > size {
		limit = size - offset
	}

	if _, err := fileFrom.Seek(offset, io.SeekStart); err != nil {
		log.Panic(err)
		return ErrUnsupportedFile
	}

	if _, err := io.CopyN(fileTo, fileFrom, limit); err != nil {
		log.Panic(err)
		return ErrUnsupportedFile
	}

	return nil
}
