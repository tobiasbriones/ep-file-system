// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

// This test suite consists of side effect test cases for using the TCP file
// system. Test file system requirements have to be provided, and the results
// can also be checked in that test file system. The test file system is
// defined as the server/.test_fs directory and the root app file system into
// its server directory. Channels are direct children of the FS root.
// Author: Tobias Briones

import (
	"encoding/json"
	"fs"
	"fs/files"
	"fs/process"
	"fs/utils"
	"log"
	"net"
	"testing"
	"time"
)

const (
	testChannel      = "test"
	testFile         = "file.pdf"
	testDir          = "C:/Users/tobi/go/src/github.com/tobiasbriones/ep-tcp-file-system/server/.test_fs/"
	testFsRoot       = testDir + "server"
	testFsClientRoot = testDir + "client"
)

// Requires a file "main/file.pdf" into the server's test file system.
// It tests the server physical file system for writing and reading.
// It first mimics downloading a file into a monolithic memory buffer, and then
// it writes that file back to the server as a new file.
func TestReceiveSend(t *testing.T) {
	// Get a FileInfo recording the underlying user file
	serverFileInfo := newTestFileInfo()

	// Channels are direct directories, get the FS File of the file
	file, _ := fs.NewFileFromString(
		process.DefChannel + fs.Separator + serverFileInfo.Value,
	)

	// Get the physical file from the test directory
	osFile := file.ToOsFile(testFsRoot)
	size, err := files.ReadSize(osFile)
	utils.RequirePassCase(t, err, "Fail to load test file info")

	// Start downloading the file from the test server
	downloaded := make([]byte, 0, size)
	err = files.Stream(
		osFile,
		bufSize,
		func(buf []byte) {
			downloaded = append(downloaded, buf...)
		},
	)
	utils.RequirePassCase(t, err, "Fail to stream file")

	// Upload the in-memory file back
	newFile, _ := fs.NewFileFromString("main/file-uploaded-back.pdf")
	newOsFile := newFile.ToOsFile(testFsRoot)
	err = files.Create(newOsFile)
	utils.RequirePassCase(t, err, "Fail to create file file-uploaded-back.pdf")
	for i := 0; i < cap(downloaded); i += bufSize {
		end := i + bufSize

		if end >= cap(downloaded) {
			end = cap(downloaded) - 1
		}
		chunk := downloaded[i:end]

		// Mimic sending to remote server
		err = files.WriteBuf(newOsFile, chunk)
		utils.RequirePassCase(t, err, "Fail to write chunk")
	}
}

// Makes a request to the server. It can be either upload or download. After the
// initial request (state START) the server will respond with state ERROR
// because the file sent is empty.
func TestTcpConn(t *testing.T) {
	info := newTestFileInfo()
	info.Size = 0 // Don't upload anything, just initiate a connection and wait
	conn := initiateConn(t, process.ActionUpload, info)
	defer conn.Close()

	res := readResponseMsg(t, conn)
	if res.State != process.Error { // The file sent is empty, ERROR must be responded.
		t.Fatal("Fail to establish the TCP connection to the server")
	}
	resPayload := Payload{Data: res.Data}

	// res.State is Error so this conversion is safe
	payload, _ := resPayload.ErrorPayload()
	if payload.Message != "file sent is empty" {
		t.Fatal("Fail to get error message", string(res.Data))
	}
}

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

func TestChannelList(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "LIST_CHANNELS"
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	var channels []string
	dec := json.NewDecoder(conn)
	err = dec.Decode(&channels)
	if err != nil {
		return
	}

	// Check at least has the main, and test channels
	if !utils.StringSliceContains(channels, "main") {
		t.Fatal("Channel list does not contain channel: main")
	}
	if !utils.StringSliceContains(channels, "test") {
		t.Fatal("Channel list does not contain channel: test")
	}
}

func TestFileList(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "LIST_FILES"
	cmd["CHANNEL"] = testChannel
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	var fileList []string
	dec := json.NewDecoder(conn)
	err = dec.Decode(&fileList)
	if err != nil {
		return
	}

	// Check
	log.Println("Files on channel test:", fileList)
}

func TestCommandId(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "CID"
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	buf := make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal("Fail to read CID: ", err)
	}

	// Check
	log.Println("Client ID:", string(buf[:n]))
}

// Makes a simple test when a goroutine uploads a file while another client is
// on hold, so it receives an update message. It might not pass sometimes due
// to the side effects.
func TestBroadcast(t *testing.T) {
	var update = make(chan struct{})
	var fail = make(chan struct{})

	go func() {
		// Connect and keep on hold state START
		tcpAddr, _ := net.ResolveTCPAddr(network, getServerAddress())
		conn, _ := net.DialTCP(network, nil, tcpAddr)

		log.Println("HOLD: Connection established")
		msg, _ := readMessage(conn, readTimeOut)
		log.Println("Received msg:", msg)
		if msg.Response != Update {
			fail <- struct{}{}
			return
		}
		update <- struct{}{}
	}()
	go func() {
		time.Sleep(1 * time.Second)

		// Upload a file to change the FS, channel test
		info := newTestFileInfo()
		osFile := info.ToOsFile(testFsClientRoot) // .../.test_fs/client/file.pdf
		loadFileSize(&info, osFile)
		conn := initiateConn(t, process.ActionUpload, info)
		defer conn.Close()
		readResponseMsg(t, conn) // State DATA
		upload(t, conn, osFile)
		readResponseMsg(t, conn) // State EOF
		eof(t, conn)
		readResponseMsg(t, conn) // State DONE
		log.Println("UPLOAD: File sent")
	}()

	select {
	case <-update:
		log.Println("Update received")
	case <-fail:
		t.Fatal("Fail to receive update from server")
	}
}

// Tests if the server closes the connection after the read timeout is consumed.
func TestTcpTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skip long running test: TCP timeout")
	}
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

	// Server is waiting for client chunks ...
	time.Sleep(readTimeOut + 1)

	res = readResponseMsg(t, conn)
	if res.State != process.Error {
		t.Fatal("fail to get state ERROR after timeout")
	}
}

func upload(t *testing.T, conn *net.TCPConn, file fs.OsFile) {
	log.Println("Streaming file to server:", file.Path())
	err := files.Stream(file, bufSize, func(buf []byte) {
		_, err := conn.Write(buf)
		utils.RequirePassCase(t, err, "Fail to write chunk to server")
	})
	utils.RequirePassCase(t, err, "Fail to stream file")
}

func eof(t *testing.T, conn *net.TCPConn) {
	err := writeState(process.Eof, conn)
	utils.RequirePassCase(t, err, "Fail to write EOF")
}

func initiateConn(
	t *testing.T,
	action process.Action,
	info fs.FileInfo,
) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")

	conn, err := net.DialTCP(network, nil, tcpAddr)
	utils.RequirePassCase(t, err, "Fail to establish connection")

	body := process.StartPayload{
		Action:   action,
		FileInfo: info,
		Channel:  process.NewChannel(testChannel),
	}
	utils.RequirePassCase(t, err, "Fail to load test FileInfo")

	payload, err := NewPayload(body)
	utils.RequirePassCase(t, err, "Fail to load create payload")

	msg := Message{
		State:   process.Start,
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

func newTestFileInfo() fs.FileInfo {
	f, _ := fs.NewFileFromString(testFile)
	i := fs.FileInfo{
		File: f,
		Size: 0,
	}
	return i
}

func loadFileSize(info *fs.FileInfo, file fs.OsFile) error {
	size, err := files.ReadSize(file)
	info.Size = uint64(size)
	return err
}
