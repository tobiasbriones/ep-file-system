// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"testing"
)

// Side effect test. Requires a file "file.pdf" into the server's file system
// directory.
func TestReceiveSend(t *testing.T) {
	path := "file.pdf"
	size, err := GetFileSize(path)

	if err != nil {
		t.Fatal(err)
	}
	downloaded := make([]byte, 0, size)
	ds := newDataStream(path, bufSize, func(buf []byte) {
		downloaded = append(downloaded, buf...)
	})

	StreamFile(&ds) // blocking

	// Upload the file back
	newPath := "new-file.pdf"
	CreateFile(newPath)
	for i := 0; i < cap(downloaded); i += bufSize {
		end := i + bufSize

		if end >= cap(downloaded) {
			end = cap(downloaded) - 1
		}
		chunk := downloaded[i:end]

		// Mimic sending to remote server
		WriteBuf(newPath, chunk)
	}
}
