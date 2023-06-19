package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
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
			log.Printf("File %v is not exist\n", fromPath)
		}
		return ErrUnsupportedFile
	}
	defer fileFrom.Close()

	// create file for copy
	fileTo, err := os.Create(toPath)
	if err != nil {
		log.Printf("Failed to create file %v with error: %v\n", toPath, err)
		return ErrUnsupportedFile
	}
	defer fileTo.Close()

	// get FileInfo
	fi, err := fileFrom.Stat()
	if err != nil {
		log.Println(err)
		return ErrUnsupportedFile
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// get file size
	size := fi.Size()

	// compare offset and file size
	if offset > size {
		log.Printf("Offset value: %v is bigger than file size: %v\n", offset, size)
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit+offset > size {
		limit = size - offset
	}

	// start new bar
	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	bar.SetRefreshRate(time.Nanosecond)
	bar.Set(pb.SIBytesPrefix, true)

	if _, err := fileFrom.Seek(offset, io.SeekStart); err != nil {
		log.Println(err)
		return ErrUnsupportedFile
	}

	// create proxy reader
	barReader := bar.NewProxyReader(fileFrom)

	if _, err := io.CopyN(fileTo, barReader, limit); err != nil {
		log.Println(err)
		return ErrUnsupportedFile
	}

	// finish bar
	bar.Finish()

	return nil
}
