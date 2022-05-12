// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package server

import (
	"fmt"
	"os"
)

type DataStream struct {
	path    string
	bufSize uint
	handler Handler
}

func newDataStream(relPath string, bufSize uint, handler Handler) DataStream {
	path := getFilePath(relPath)
	return DataStream{path, bufSize, handler}
}

type Handler func(buf []byte)

func getFilePath(relPath string) string {
	return fmt.Sprintf("%v%v%v", fsRootPath, os.PathSeparator, relPath)
}
