// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"fs/server/io"
	"fs/utils"
	"log"
	"net"
	"testing"
)

const (
	testFile      = "file.pdf"
	testLocalFile = "C:\\file.pdf"
	testChannel   = "test"
)

// Side effect test. Requires a file "file.pdf" into the server's file system
// directory. It tests the server file system for write and read.
func TestReceiveSend(t *testing.T) {
	serverFileInfo := newTestFileInfo()
	size, err := serverFileInfo.ReadFileSize(io.DefChannel)
	utils.RequirePassCase(t, err, "Fail to load test file info")

	downloaded := make([]byte, 0, size)
	err = serverFileInfo.Stream(
		io.DefChannel,
		bufSize,
		func(buf []byte) {
			downloaded = append(downloaded, buf...)
		},
	)
	utils.RequirePassCase(t, err, "Fail to stream file")

	// Upload the file back
	newPath := "new-file.pdf"
	err = io.CreateFile(newPath)
	utils.RequirePassCase(t, err, "Fail to create file new-file.pdf")
	for i := 0; i < cap(downloaded); i += bufSize {
		end := i + bufSize

		if end >= cap(downloaded) {
			end = cap(downloaded) - 1
		}
		chunk := downloaded[i:end]

		// Mimic sending to remote server
		err = io.WriteBuf(newPath, chunk)
		utils.RequirePassCase(t, err, "Fail to write chunk")
	}
}

// Makes a request to the server. It can be either upload or download. After the
// initial request (state START), the server will respond with state OK.
func TestTcpConn(t *testing.T) {
	info := newTestFileInfo()
	info.Size = 0 // Don't upload anything, just initiate a connection and wait
	conn := initiateConn(t, ActionUpload, info)
	defer conn.Close()

	res := readResponseMsg(t, conn)
	if res.State != Error { // The file sent is empty, ERROR must be responded.
		t.Fatal("Fail to establish the TCP connection to the server")
	}
}

// Side effect. Requires testLocalFile = "C:\\file.pdf".
func TestUpload(t *testing.T) {
	info, err := newTestLocalFileInfo()
	utils.RequirePassCase(t, err, "Fail to read file info")
	conn := initiateConn(t, ActionUpload, info)
	defer conn.Close()

	res := readResponseMsg(t, conn)
	if res.State != Data {
		t.Fatal("Fail to get state=DATA")
	}
	log.Println("State=DATA")
	upload(t, conn, testLocalFile)
	log.Println("Uploaded")

	res = readResponseMsg(t, conn)
	if res.State != Eof {
		t.Fatal("Fail to get state=EOF")
	}

	log.Println("State=EOF", res)
	eof(t, conn)
	res = readResponseMsg(t, conn)
	log.Println(res.State)
}

// Requires the file testFile = "file.pdf" in the server FS at channel "test",
//and will write it to "download.pdf" into this source code directory.
func TestDownload(t *testing.T) {
	info := newTestFileInfo()
	conn := initiateConn(t, ActionDownload, info)
	defer conn.Close()
	res := readResponseMsg(t, conn)
	if res.State != Stream {
		t.Fatal("Fail to get state=STREAM")
	}
	payload, err := res.StreamPayload()
	utils.RequirePassCase(t, err, "Fail to read StreamPayload")
	err = writeState(Stream, conn)
	utils.RequirePassCase(t, err, "Fail to write state=STREAM")
	path := "download.pdf"
	err = io.CreateFile(path)
	utils.RequirePassCase(t, err, "Fail to create file download.pdf")
	size := uint64(payload.Size)
	count := uint64(0)
	log.Println(size)
	for {
		if count >= size {
			break
		}
		b := make([]byte, bufSize)
		n, err := conn.Read(b)
		utils.RequirePassCase(t, err, "Fail to read chunk from server")
		chunk := b[:n]
		err = io.WriteBuf(path, chunk)
		utils.RequirePassCase(t, err, "Fail to write chunk to file")
		count += uint64(n)
		if n == 0 {
			t.Fatal("Underflow!")
		}
	}
	log.Println(count)
	if count != size {
		// TODO The download works, but extra bytes are written
		t.Fatal("Overflow!")
	}
}

// Requires not to have a file "not-exists.txt" in the server utils.
func TestDownloadIfNotExists(t *testing.T) {
	info := io.FileInfo{
		RelPath: "not-exists",
		Size:    0,
	}
	conn := initiateConn(t, ActionDownload, info)
	defer conn.Close()
	res := readResponseMsg(t, conn)
	if res.State != Error {
		t.Fatal("Fail to get state=ERROR")
	}
}

func upload(t *testing.T, conn *net.TCPConn, path string) {
	log.Println("Streaming file to server:", path)
	err := io.StreamLocalFile(path, bufSize, func(buf []byte) {
		_, err := conn.Write(buf)
		utils.RequirePassCase(t, err, "Fail to write chunk to server")
	})
	utils.RequirePassCase(t, err, "Fail to stream file")
}

func eof(t *testing.T, conn *net.TCPConn) {
	err := writeState(Eof, conn)
	utils.RequirePassCase(t, err, "Fail to write EOF")
}

func initiateConn(t *testing.T, action Action, info io.FileInfo) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")

	conn, err := net.DialTCP(network, nil, tcpAddr)
	utils.RequirePassCase(t, err, "Fail to establish connection")

	body := StartPayload{
		Action:   action,
		FileInfo: info,
		Channel:  NewChannel(testChannel),
	}
	utils.RequirePassCase(t, err, "Fail to load test FileInfo")

	payload, err := NewPayload(body)
	utils.RequirePassCase(t, err, "Fail to load create payload")

	msg := Message{
		State:   Start,
		Payload: payload,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write state=START to the server")
	return conn
}

func readResponseMsg(t *testing.T, conn net.Conn) Message {
	var msg Message
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	utils.RequirePassCase(t, err, "Fail to read response from server")
	return msg
}

func newTestFileInfo() io.FileInfo {
	i := io.FileInfo{
		RelPath: testFile,
		Size:    0,
	}
	return i
}

func newTestLocalFileInfo() (io.FileInfo, error) {
	i := io.FileInfo{
		RelPath: testFile,
		Size:    0,
	}
	size, err := io.ReadFileSize(testLocalFile)
	i.Size = size
	return i, err
}
