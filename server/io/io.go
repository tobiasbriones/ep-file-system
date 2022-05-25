// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

const (
	DefChannel = "main"
	fsRootPath = ".fs"
)

type Handle func(buf []byte)

// FileInfo Defines a DTO to transfer data from server to client.
type FileInfo struct {
	RelPath string
	Size    int64
}

// ReadFileSize Returns the file size read from the server file system.
func (i *FileInfo) ReadFileSize(channel string) (int64, error) {
	file, err := i.ToFile(channel)
	if err != nil {
		return 0, err
	}
	return ReadFileSize(file.value)
}

func (i *FileInfo) Stream(channel string, bufSize uint, handle Handle) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	return StreamLocalFile(file.value, bufSize, handle)
}

func (i *FileInfo) Exists(channel string) (bool, error) {
	file, err := i.ToFile(channel)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(file.value); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, nil
}

func (i *FileInfo) Create(channel string) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	return CreateFile(file.value)
}

func (i *FileInfo) WriteChunk(channel string, chunk []byte) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	return WriteBuf(file.value, chunk)
}

func (i *FileInfo) ChannelPath(channel string) (Path, error) {
	return getChannelPath(channel)
}

func (i *FileInfo) ToFile(channel string) (File, error) {
	path, err := getPath(i.RelPath, channel)
	if err != nil {
		return File{}, err
	}
	return File{Path: path}, nil
}

func StreamLocalFile(path string, bufSize uint, handle Handle) error {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Fail to read file %v: %v", path, err.Error())
	}
	defer f.Close()
	buf := make([]byte, 0, bufSize)
	reader := bufio.NewReader(f)
	bytesNumber, chunksNumber, err := stream(reader, buf, handle)

	if err != nil {
		log.Println(
			"Streaming completed.\n",
			"File:",
			path,
			"Bytes:",
			bytesNumber,
			"Chunks:", chunksNumber,
		)
	}
	return err
}

func stream(
	reader *bufio.Reader,
	buf []byte,
	handle Handle) (int64, int64, error) {
	bytesNumber := int64(0)
	chunksNumber := int64(0)

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
		chunksNumber++
		bytesNumber += int64(len(buf))

		handle(buf)

		if err != nil && err != io.EOF {
			return 0, 0, err
		}
	}
	return bytesNumber, chunksNumber, nil
}

func getChannelPath(channel string) (Path, error) {
	// TODO File needs to implement Parent()
	path, err := NewPathFrom(fsRootPath, channel)
	if err != nil {
		return Path{}, err
	}
	return path, nil
}

func getPath(relPath string, channel string) (Path, error) {
	path, err := NewPathFrom(fsRootPath, channel)
	if err != nil {
		return Path{}, err
	}
	children := strings.Split(relPath, Separator)
	err = path.Append(children...)
	if err != nil {
		return Path{}, err
	}
	return path, nil
}
