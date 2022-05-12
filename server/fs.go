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

func GetFileSize(path string) (int64, error) {
	f, err := os.Open(getFilePath(path))
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func StreamFile(stream *DataStream) {
	f, err := os.Open(stream.path)
	if err != nil {
		log.Fatalf("Fail to read file %v: %v", stream.path, err.Error())
	}

	bytesNumber := int64(0)
	chunksNumber := int64(0)
	reader := bufio.NewReader(f)
	buf := make([]byte, 0, stream.bufSize)

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

		stream.handle(buf)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	log.Println(
		"Streaming completed.\n",
		"File:",
		stream.path,
		"Bytes:",
		bytesNumber,
		"Chunks:", chunksNumber,
	)
}

func getFilePath(relPath string) string {
	return fmt.Sprintf("%v%v%v", fsRootPath, string(os.PathSeparator), relPath)
}
