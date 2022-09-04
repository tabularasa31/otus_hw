package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			continue
		}

		fpath := filepath.Join(dir, file.Name())

		val, err := ReadFile(fpath)
		if err != nil {
			return nil, err
		}

		fi, err := os.Stat(fpath)
		if err != nil {
			return nil, err
		}

		if fi.Size() == 0 {
			env[file.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
		}

		env[file.Name()] = EnvValue{
			Value: val,
		}
	}

	return env, nil
}

func ReadFile(fpath string) (string, error) {

	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Fatalf("failed to close file %v, error: %q", fpath, e)
		}
	}()

	reader := bufio.NewReader(f)
	val, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	strings.ReplaceAll(val, "\x00", "\n")
	val = strings.TrimRight(val, " \n\t\r")

	return val, nil
}
