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
		if err != nil {
			log.Println("Fail to accept client")
			continue
		}
		client := newClient(conn)
		go client.run()
	}
}

func readChunk(conn net.Conn) ([]byte, error) {
	b := make([]byte, bufSize)
	n, err := conn.Read(b)

	if err != nil {
		if err.Error() != "EOF" {
			return []byte{}, err
		}
	}
	return b[:n], nil
}

func writeState(state State, conn net.Conn) error {
	msg := Message{
		State: state,
	}
	return writeMessage(msg, conn)
}

func writeMessage(msg Message, conn net.Conn) error {
	enc := json.NewEncoder(conn)
	return enc.Encode(msg)
}

func readMessage(conn net.Conn) (Message, error) {
	var msg Message
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	return msg, err
}
