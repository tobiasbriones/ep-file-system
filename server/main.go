// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

// Entry point for the file system server.
//
// Author: Tobias Briones
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	port    = 8080
	network = "tcp"
	bufSize = 1024
)

type Status string

const (
	START Status = "start"
	OK    Status = "ok"
	DATA  Status = "data"
	EOF   Status = "eof"
	ERROR Status = "error"
)

type Message struct {
	Status
	Action  string
	Payload string
}

func main() {
	server, err := net.Listen(network, getServerAddress())

	defer server.Close()
	requireNoError(err)
	listen(server)
}

func listen(server net.Listener) {
	for {
		conn, err := server.Accept()
		requireNoError(err)
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)

	var msg Message

	err := dec.Decode(&msg)

	if err != nil {
		log.Println("Skipped: ", conn)
		return
	}
	log.Println(msg.Status)
	switch msg.Status {
	case START:
		handleStatusStart(conn, msg)
		break
	default:
		writeStatus(ERROR, conn)
		break
	}
}

func handleStatusStart(conn net.Conn, msg Message) {
	log.Println(msg.Action)
	switch msg.Action {
	case "upload":
		handleUpload(conn, msg)
		break
	case "download":
		handleDownload(conn, msg)
		break
	default:
		writeStatus(ERROR, conn)
		break
	}
}

func handleDownload(conn net.Conn, msg Message) {
	writeStatus(OK, conn)
	// TODO
}

func handleUpload(conn net.Conn, msg Message) {
	writeStatus(OK, conn)
	// TODO
}

func writeStatus(status Status, conn net.Conn) {
	msg := Message{
		Status: status,
	}
	b, err := json.Marshal(msg)

	requireNoError(err)
	_, err = conn.Write(b)

	requireNoError(err)
}

func getServerAddress() string {
	return fmt.Sprintf("localhost:%v", port)
}
