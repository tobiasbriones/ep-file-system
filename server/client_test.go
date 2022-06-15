// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs"
	"fs/files"
	"fs/process"
	"fs/utils"
	"log"
	"testing"
)

// Requires client file ".../.test_fs/client/file.pdf" to upload it to the
// server.
func TestUpload(t *testing.T) {
	info := newTestFileInfo()
	osFile := info.ToOsFile(testFsClientRoot) // .../.test_fs/client/file.pdf
	err := loadFileSize(&info, osFile)
	utils.RequirePassCase(t, err, "Fail to read file info")
	conn := initiateConn(t, process.ActionUpload, info)
	defer conn.Close()

	res := readResponseMsg(t, conn)
	if res.State != process.Data {
		t.Fatal("Fail to get state=DATA")
	}
	log.Println("State=DATA")
	upload(t, conn, osFile)
	log.Println("Uploaded")

	res = readResponseMsg(t, conn)
	if res.State != process.Eof {
		t.Fatal("Fail to get state=EOF")
	}

	log.Println("State=EOF")
	eof(t, conn)
	res = readResponseMsg(t, conn)
	log.Println(res.State)
}

// Requires the file testFile = "file.pdf" in the server FS at channel "test",
//and will write it to "download.pdf" into this source code directory.
func TestDownload(t *testing.T) {
	info := newTestFileInfo()
	conn := initiateConn(t, process.ActionDownload, info)
	defer conn.Close()

	// Receive state STREAM with payload
	res := readResponseMsg(t, conn)
	if res.State != process.Stream {
		t.Fatal("Fail to get state=STREAM")
	}
	payload, err := res.StreamPayload()
	utils.RequirePassCase(t, err, "Fail to read StreamPayload")

	// Write state STREAM to confirm
	err = writeState(process.Stream, conn)
	utils.RequirePassCase(t, err, "Fail to write state=STREAM")

	// Get local file ready
	file, _ := fs.NewFileFromString("download.pdf")
	osFile := file.ToOsFile(testFsClientRoot) // .../.fs/client/download.pdf
	err = files.Create(osFile)
	utils.RequirePassCase(t, err, "Fail to create file download.pdf")
	size := payload.Size
	count := uint64(0)

	writeChunk := func(chunk []byte) {
		err = files.WriteBuf(osFile, chunk)
		utils.RequirePassCase(t, err, "Fail to write chunk to file")
	}

	// Read stream chunks from the server
	for {
		b := make([]byte, bufSize)
		n, err := conn.Read(b)
		utils.RequirePassCase(t, err, "Fail to read chunk from server")
		count += uint64(n)

		chunk := b[:n]
		writeChunk(chunk)
		if n == 0 { // TODO Underflow should be handled differently now :/
			t.Fatal("Underflow!")
		}
		if count >= size {
			break
		}
	}

	if count != size {
		t.Fatal("Overflow")
	}

	// Confirm received with state EOF
	err = writeState(process.Eof, conn)
	if err != nil {
		t.Fatal("fail to write STATUS=EOF")
	}

	// And then get state DONE
	msg, err := readMessage(conn, readTimeOut)
	if err != nil {
		t.Fatal("fail to read STATUS=DONE")
	}
	if msg.State != process.Done {
		t.Fatal("fail to read STATUS=EOF, it might be overflow")
	}
}

// Requires not to have a file "not-exists.txt" in the server test channel.
func TestDownloadIfNotExists(t *testing.T) {
	file, _ := fs.NewFileFromString("not-exists.txt")
	info := fs.FileInfo{File: file}
	conn := initiateConn(t, process.ActionDownload, info)
	defer conn.Close()
	res := readResponseMsg(t, conn)
	if res.State != process.Error {
		t.Fatal("Fail to get state=ERROR")
	}
}
