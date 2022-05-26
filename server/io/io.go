// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package io

import (
	"bufio"
	"fs"
	"io"
	"log"
	"os"
)

const (
	DefChannel = "main"
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
	return ReadSize(file)
}

func (i *FileInfo) Stream(channel string, bufSize uint, handle Handle) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	path, err := AbsolutePath(file)
	if err != nil {
		return err
	}
	return StreamLocalFile(path, bufSize, handle)
}

func (i *FileInfo) Exists(channel string) (bool, error) {
	file, err := i.ToFile(channel)
	if err != nil {
		return false, err
	}
	return Exists(file)
}

func (i *FileInfo) Create(channel string) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	return Create(file)
}

func (i *FileInfo) WriteChunk(channel string, chunk []byte) error {
	file, err := i.ToFile(channel)
	if err != nil {
		return err
	}
	path, err := AbsolutePath(file)
	if err != nil {
		return err
	}
	return WriteBuf(path, chunk)
}

func (i *FileInfo) CreateChannelIfNotExists(channel string) error {
	file, err := fs.NewFileFromString(channel)
	if err != nil {
		return err
	}
	return CreateIfNotExists(file)
}

func (i *FileInfo) ToFile(channel string) (fs.File, error) {
	path, err := getPath(i.RelPath, channel)
	if err != nil {
		return fs.File{}, err
	}
	return fs.File{Path: path}, nil
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
