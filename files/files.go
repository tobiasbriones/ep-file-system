// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package files

import (
	"bufio"
	"errors"
	"fs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Exists(file fs.OsFile) (bool, error) {
	if _, err := os.Stat(file.Path()); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, nil
}

func Create(file fs.OsFile) error {
	f, err := os.Create(file.Path())
	f.Close()
	return err
}

func CreateIfNotExists(file fs.OsFile) error {
	return os.MkdirAll(file.Path(), os.ModePerm)
}

func DeleteIfExists(file fs.OsFile) error {
	return os.RemoveAll(file.Path())
}

func ReadSize(file fs.OsFile) (int64, error) {
	f, err := os.Open(file.Path())
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// ReadFileNames returns a list of file names that are children of the given
// file.
func ReadFileNames(file fs.OsFile) ([]string, error) {
	files, err := ioutil.ReadDir(file.Path())
	if err != nil {
		return nil, err
	}
	var list []string
	for _, f := range files {
		list = append(list, f.Name())
	}
	return list, nil
}

func GetExecPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath := filepath.Dir(ex)
	exePath = strings.ReplaceAll(exePath, "\\", fs.Separator)
	return exePath, nil
}

func WriteBuf(file fs.OsFile, buf []byte) error {
	path := file.Path()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	f.Close()
	return err
}

type Handle func(buf []byte)

func Stream(file fs.OsFile, bufSize uint, handle Handle) error {
	path := file.Path()
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Fail to read file %v: %v\n", path, err.Error())
		return errors.New("fail to read file")
	}
	defer f.Close()
	buf := make([]byte, 0, bufSize)
	reader := bufio.NewReader(f)
	bytesTotal, chunksTotal, err := stream(reader, buf, handle)

	if err != nil {
		log.Println(
			"Streaming completed.\n",
			"File:",
			path,
			"Bytes:",
			bytesTotal,
			"Chunks:", chunksTotal,
		)
	}
	return err
}

func stream(
	reader *bufio.Reader,
	buf []byte,
	handle Handle,
) (int64, int64, error) {
	bytesTotal := int64(0)
	chunksTotal := int64(0)

	for {
		n, err := reader.Read(buf[:cap(buf)])

		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		chunksTotal++
		bytesTotal += int64(len(buf))

		handle(buf)

		if err != nil && err != io.EOF {
			return 0, 0, err
		}
	}
	return bytesTotal, chunksTotal, nil
}
