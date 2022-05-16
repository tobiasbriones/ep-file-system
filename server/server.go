// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"log"
	"net"
)

func listen(server net.Listener) {
	for {
		conn, err := server.Accept()
		requireNoError(err)
		client := newClient(conn)
		go client.run()
	}
}

func readChunk(conn net.Conn) []byte {
	b := make([]byte, bufSize)
	n, err := conn.Read(b)

	if err != nil {
		if err.Error() != "EOF" {
			log.Println("Error reading chunk:", err)
			requireNoError(err)
		}
	}
	return b[:n]
}

func writeState(state State, conn net.Conn) {
	msg := Message{
		State: state,
	}
	enc := json.NewEncoder(conn)
	err := enc.Encode(msg)
	requireNoError(err)
}

func readMessage(conn net.Conn) (Message, error) {
	var msg Message
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	return msg, err
}
