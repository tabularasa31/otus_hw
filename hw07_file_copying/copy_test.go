package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tmpfile, err := os.CreateTemp("testdata/", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if e := os.Remove(tmpfile.Name()); e != nil {
			log.Fatalf("failed to remove temp dir, error: %q", e)
			return
		}
	}()

	srcFile, _ := os.Open("testdata/input.txt")
	sf, _ := srcFile.Stat()

	t.Run("error: unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", tmpfile.Name(), 0, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unsupported file")
	})
	t.Run("error: offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpfile.Name(), 100000, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "offset exceeds file size")
	})
	t.Run("success: full copy", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile.Name(), 0, 0)
		require.NoError(t, err)
		dstFile, _ := os.Open(tmpfile.Name())
		df, _ := dstFile.Stat()
		require.Equal(t, sf.Size(), df.Size())
	})
	t.Run("success: limit exceeds file size", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile.Name(), 0, sf.Size()+10)
		dstFile, _ := os.Open(tmpfile.Name())
		df, _ := dstFile.Stat()
		require.Equal(t, sf.Size(), df.Size())
	})
	t.Run("offset=100 limit=1000", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile.Name(), 100, 1000)
		sampleFile, _ := os.Open("testdata/out_offset100_limit1000.txt")
		sample, _ := sampleFile.Stat()
		dstFile, _ := os.Open(tmpfile.Name())
		df, _ := dstFile.Stat()
		require.Equal(t, sample.Size(), df.Size())
	})
	t.Run("offset=6000 limit=1000", func(t *testing.T) {
		err = Copy("testdata/input.txt", tmpfile.Name(), 6000, 1000)
		sampleFile, _ := os.Open("testdata/out_offset6000_limit1000.txt")
		sample, _ := sampleFile.Stat()
		dstFile, _ := os.Open(tmpfile.Name())
		df, _ := dstFile.Stat()
		require.Equal(t, sample.Size(), df.Size())
	})
}
