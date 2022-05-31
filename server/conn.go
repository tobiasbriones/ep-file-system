// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"encoding/json"
	"fs/files"
	"fs/process"
	"net"
	"time"
)

const (
	readTimeOut = 20 * time.Second
)

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

func writeResponse(res Response, conn net.Conn) error {
	msg := Message{
		Response: res,
	}
	return writeMessage(msg, conn)
}

func writeState(state process.State, conn net.Conn) error {
	msg := Message{
		State: state,
	}
	return writeMessage(msg, conn)
}

func writeErrorState(errorMsg string, conn net.Conn) error {
	payload := ErrorPayload{Message: errorMsg}
	p, _ := NewPayloadFrom(payload)
	msg := Message{
		State:   process.Error,
		Payload: p,
	}
	return writeMessage(msg, conn)
}

func writeMessage(msg Message, conn net.Conn) error {
	enc := json.NewEncoder(conn)
	return enc.Encode(msg)
}

func readMessage(conn net.Conn) (Message, error) {
	err := conn.SetReadDeadline(time.Now().Add(readTimeOut))
	if err != nil {
		return Message{}, err
	}
	var msg Message
	dec := json.NewDecoder(conn)
	err = dec.Decode(&msg)
	return msg, err
}

func writeChannels(conn net.Conn) error {
	root, err := getFsRootFile()
	if err != nil {
		return err
	}
	channels, err := files.ReadDirectories(root)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(conn)
	return enc.Encode(channels)
}
