package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %q is not exist", fromPath)
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
		return fmt.Errorf("create file '%v': %w", toPath, err)
	}

	defer func() {
		if e := dstFile.Close(); e != nil {
			log.Fatalf("failed to close file %v, error: %q", toPath, e)
		}
	}()

	sf, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("get file info: %w", err)
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

	fmt.Printf("Coping file %v to file %v\n", fromPath, toPath)

	tmpl := `{{rtime . "%s remain"}} {{bar . "<" "oOo" "|" "~" ">"}} {{speed . | rndcolor }} {{percent .}}`

	bar := pb.ProgressBarTemplate(tmpl).Start64(limit)
	barReader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()

	return nil
}
