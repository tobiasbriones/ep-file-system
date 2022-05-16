// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	fsRootPath = "fs"
)

type DataStream struct {
	path    string
	bufSize uint
	handle  Handle
}

func newDataStream(relPath string, bufSize uint, handler Handle) DataStream {
	path := getFilePath(relPath)
	return DataStream{path, bufSize, handler}
}

type Handle func(buf []byte)

type FileInfo struct {
	RelPath string
	Size    int64
}

// ReadFileSize Returns the file size read from the server file system.
func (i *FileInfo) readFileSize() (int64, error) {
	return ReadFileSize(i.getPath())
}

// Returns the file path in the server file system.
func (i *FileInfo) getPath() string {
	return getFilePath(i.RelPath)
}

func ReadFileSize(path string) (int64, error) {
	f, err := os.Open(path)
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

func StreamFile(ds *DataStream) {
	f, err := os.Open(ds.path)
	if err != nil {
		log.Fatalf("Fail to read file %v: %v", ds.path, err.Error())
	}
	buf := make([]byte, 0, ds.bufSize)
	reader := bufio.NewReader(f)
	bytesNumber, chunksNumber := stream(reader, buf, ds.handle)

	log.Println(
		"Streaming completed.\n",
		"File:",
		ds.path,
		"Bytes:",
		bytesNumber,
		"Chunks:", chunksNumber,
	)
}

func StreamLocalFile(path string, bufSize uint, handle Handle) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Fail to read file %v: %v", path, err.Error())
	}
	defer f.Close()
	buf := make([]byte, 0, bufSize)
	reader := bufio.NewReader(f)
	bytesNumber, chunksNumber := stream(reader, buf, handle)

	log.Println(
		"Streaming completed.\n",
		"File:",
		path,
		"Bytes:",
		bytesNumber,
		"Chunks:", chunksNumber,
	)
}

func stream(reader *bufio.Reader, buf []byte, handle Handle) (int64, int64) {
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
			log.Fatal(err)
		}
		chunksNumber++
		bytesNumber += int64(len(buf))

		handle(buf)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	return bytesNumber, chunksNumber
}

func CreateFile(relPath string) {
	path := getFilePath(relPath)
	CreateLocalFile(path)
}

func CreateLocalFile(path string) {
	f, err := os.Create(path)
	requireNoError(err)
	f.Close()
}

func WriteBuf(relPath string, buf []byte) {
	path := getFilePath(relPath)
	WriteLocalBuf(path, buf)
}

func WriteLocalBuf(path string, buf []byte) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	requireNoError(err)
	_, err = f.Write(buf)
	requireNoError(err)
	f.Close()
}

func getFilePath(relPath string) string {
	return fmt.Sprintf("%v%v%v", fsRootPath, string(os.PathSeparator), relPath)
}
