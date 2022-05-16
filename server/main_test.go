// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"net"
	"testing"
)

const (
	testFile = "file.pdf"
)

// Side effect test. Requires a file "file.pdf" into the server's file system
// directory. It tests the server file system for write and read.
func TestReceiveSend(t *testing.T) {
	serverFileInfo, err := newTestFileInfo()
	size := serverFileInfo.Size

	requirePassedTest(t, err, "Fail to load test file info")
	downloaded := make([]byte, 0, size)
	ds := newDataStream(testFile, bufSize, func(buf []byte) {
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

// Makes a request to the server. It can be either upload or download. After the
// initial request (status START), the server will respond with status OK.
func TestTcpConn(t *testing.T) {
	info, _ := newTestFileInfo()
	info.Size = 0 // Don't upload anything, just initiate a connection and wait
	conn := initiateConn(t, ActionUpload, info)
	defer conn.Close()

	res := readResponseMsg(t, conn)
	if res.Status != Error { // The file sent is empty, ERROR must be responded.
		t.Fatal("Fail to establish the TCP connection to the server")
	}
}

func initiateConn(t *testing.T, action Action, info FileInfo) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	requirePassedTest(t, err, "Fail to resolve TCP address")

	conn, err := net.DialTCP(network, nil, tcpAddr)
	requirePassedTest(t, err, "Fail to establish connection")

	body := StartPayload{
		Action:   action,
		FileInfo: info,
	}
	requirePassedTest(t, err, "Fail to load test FileInfo")

	payload, err := NewPayload(body)
	requirePassedTest(t, err, "Fail to load create payload")

	msg := Message{
		Status:  Start,
		Payload: payload,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	requirePassedTest(t, err, "Fail to write status=START to the server")
	return conn
}

func readResponseMsg(t *testing.T, conn net.Conn) Message {
	var msg Message
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	requirePassedTest(t, err, "Fail to read response from server")
	return msg
}

func newTestFileInfo() (FileInfo, error) {
	i := FileInfo{
		RelPath: testFile,
		Size:    0,
	}
	size, err := i.readFileSize()
	i.Size = size
	return i, err
}
