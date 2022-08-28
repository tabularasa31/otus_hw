package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %v is not exist", fromPath)
		}
		return fmt.Errorf("failed to open file %v, error: %w", fromPath, err)
	}

	defer func() {
		if e := srcFile.Close(); e != nil {
			log.Fatalf("failed to close file %v, error: %q", fromPath, e)
			return
		}
	}()

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("unable to create file %v, error: %w", toPath, err)
	}

	defer func() {
		if e := dstFile.Close(); e != nil {
			log.Fatalf("failed to close file %v, error: %q", toPath, e)
			return
		}
	}()

	sf, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if sf.Size() < offset { // offset больше, чем размер файла - невалидная ситуация
		return ErrOffsetExceedsFileSize
	}
	if !sf.Mode().IsRegular() || sf.Size() == 0 { // если не файл или длина неизвестна - невалидная ситуация
		return ErrUnsupportedFile
	}

	if limit == 0 || limit > sf.Size() {
		limit = sf.Size()
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.CopyN(dstFile, srcFile, limit)
	if err != nil {
		return err
	}

	return nil
}
