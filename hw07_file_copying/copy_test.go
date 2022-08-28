package main

import (
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	dir, err := os.MkdirTemp("testdata/", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if e := os.RemoveAll(dir); e != nil {
			log.Fatalf("failed to remove temp dir, error: %q", e)
			return
		}
	}()

	tmpfile := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(tmpfile, []byte("content"), 0666); err != nil {
		log.Fatal(err)
	}
	srcFile, _ := os.Open("testdata/input.txt")
	sf, _ := srcFile.Stat()

	t.Run("error: unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", tmpfile, 0, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unsupported file")
	})
	t.Run("error: offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpfile, 100000, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "offset exceeds file size")
	})
	t.Run("success: full copy", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile, 0, 0)
		require.NoError(t, err)
		dstFile, _ := os.Open(tmpfile)
		df, _ := dstFile.Stat()
		require.Equal(t, sf.Size(), df.Size())
	})
	t.Run("success: limit exceeds file size", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile, 0, sf.Size()+10)
		dstFile, _ := os.Open(tmpfile)
		df, _ := dstFile.Stat()
		require.Equal(t, sf.Size(), df.Size())
	})
	t.Run("offset=100 limit=1000", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile, 100, 1000)
		sampleFile, _ := os.Open("testdata/out_offset100_limit1000.txt")
		sample, _ := sampleFile.Stat()
		dstFile, _ := os.Open(tmpfile)
		df, _ := dstFile.Stat()
		require.Equal(t, sample.Size(), df.Size())
	})
	t.Run("offset=6000 limit=1000", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile, 6000, 1000)
		sampleFile, _ := os.Open("testdata/out_offset6000_limit1000.txt")
		sample, _ := sampleFile.Stat()
		dstFile, _ := os.Open(tmpfile)
		df, _ := dstFile.Stat()
		require.Equal(t, sample.Size(), df.Size())
	})
}
