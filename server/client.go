// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"errors"
	"log"
	"net"
	"os"
)

type Client struct {
	conn  net.Conn
	state State
	req   StartPayload
	count int64
}

func newClient(
	conn net.Conn,
) *Client {
	return &Client{
		conn:  conn,
		state: Start,
	}
}

func (c *Client) run() {
	defer c.conn.Close()
	for {
		switch c.state {
		default:
			c.listenMessage()
		case Data:
			c.listenData()
		case Stream:
			c.listenStream()
		case Eof:
			c.listenEof()
		case Done, Error:
			log.Println("Exiting client")
			return
		}
	}
}

func (c *Client) listenMessage() {
	log.Println("Listening for client message")
	msg, err := readMessage(c.conn)
	if err != nil {
		c.error("Fail to read message")
		return
	}
	c.onMessage(msg)
}

func (c *Client) onMessage(msg Message) {
	log.Println("Message received:", msg)
	switch msg.State {
	case Start:
		c.start(msg)
	default:
		c.error("Wrong message state")
	}
}

func (c *Client) start(msg Message) {
	payload, err := msg.StartPayload()
	if err != nil {
		c.error("Fail to read StartPayload")
		return
	}
	c.req = payload
	c.count = 0

	switch payload.Action {
	case ActionUpload:
		c.startUpload(payload)
	case ActionDownload:
		c.startDownload(payload)
	}
}

func (c *Client) startUpload(payload StartPayload) {
	if payload.Size <= 0 {
		c.error("File sent is empty")
		return
	}
	err := CreateFile(payload.RelPath)
	if err != nil {
		c.error("Fail to create file")
		return
	}
	log.Println("Payload saved, writing state=DATA", payload)
	c.state = Data
	err = writeState(Data, c.conn)
	if err != nil {
		c.error("Fail to write state=DATA")
		return
	}
}

func (c *Client) listenData() {
	chunk, _ := readChunk(c.conn)
	c.processChunk(chunk)
	if c.count == c.req.Size {
		c.state = Eof
		log.Println("File saved, writing state EOF")
		err := writeState(Eof, c.conn)
		if err != nil {
			c.error("Fail to write state=EOF")
			return
		}
	}
}

func (c *Client) listenEof() {
	log.Println("Listening for EOF")
	msg, err := readMessage(c.conn)
	if err != nil {
		c.error("Fail to read message")
		return
	}
	c.eof(msg)
}

func (c *Client) processChunk(chunk []byte) {
	if c.overflows(chunk) {
		c.error("Overflow!")
		return
	}
	if len(chunk) == 0 {
		c.error("Underflow!")
		return
	}
	err := WriteBuf(c.req.RelPath, chunk)
	if err != nil {
		c.error("Fail to write chunk")
		return
	}
	c.count += int64(len(chunk))
}

func (c *Client) eof(msg Message) {
	if msg.State != Eof {
		c.error("Expecting EOF")
		return
	}
	log.Println("DONE!")
	c.state = Done
	err := writeState(Done, c.conn)
	if err != nil {
		c.error("Fail to write state=DONE")
		return
	}
}

func (c *Client) startDownload(payload StartPayload) {
	if _, err := os.Stat(payload.getPath()); errors.Is(err, os.ErrNotExist) {
		c.error("Requested file does not exists")
		return
	}
	size, err := ReadFileSize(payload.getPath())
	if err != nil {
		c.error("Fail to read file size")
		return
	}
	info := FileInfo{
		RelPath: payload.RelPath,
		Size:    size,
	}
	c.req = StartPayload{
		Action:   ActionDownload,
		FileInfo: info,
	}
	c.writeStreamState(StreamPayload{FileInfo: info})
}

func (c *Client) writeStreamState(payload StreamPayload) {
	p, err := NewPayloadFrom(payload)
	if err != nil {
		c.error("Fail to read payload from StreamPayload")
		return
	}
	msg := Message{
		State:   Stream,
		Payload: p,
	}
	c.state = Stream
	err = writeMessage(msg, c.conn)
	if err != nil {
		c.error("Fail to write state=STREAM")
		return
	}
	log.Println("Payload sent, writing state=STREAM", payload)
}

func (c *Client) listenStream() {
	log.Println("Listening for client STREAM signal")
	msg, err := readMessage(c.conn)
	if err != nil {
		c.error("Fail to read message")
		return
	}
	if msg.State != Stream {
		c.error("Wrong client state, state=STREAM was expected")
		return
	}
	c.stream()
}

func (c *Client) stream() {
	err := StreamLocalFile(c.req.FileInfo.getPath(), bufSize, func(buf []byte) {
		_, err := c.conn.Write(buf)
		if err != nil {
			// TODO Fix StreamLocalFile paradigm
			c.error("Fail to write chunk")
			return
		}
	})
	if err != nil {
		c.error("Fail to stream file")
		return
	}
	log.Println("File sent to client, changing state to DONE")
	c.state = Done
	err = writeState(Done, c.conn)
	if err != nil {
		c.error("Fail to write state=DONE")
		return
	}
}

func (c *Client) overflows(chunk []byte) bool {
	return c.count+int64(len(chunk)) > c.req.Size
}

func (c *Client) error(msg string) {
	// TODO update func to accept msg
	log.Println("ERROR:", msg)
	c.state = Error
	writeState(Error, c.conn)
}
