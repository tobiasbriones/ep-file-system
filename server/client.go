// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"log"
	"net"
)

type Client struct {
	conn   net.Conn
	status Status
	req    StartPayload
	count  int64
}

func newClient(
	conn net.Conn,
) *Client {
	return &Client{
		conn:   conn,
		status: Start,
	}
}

func (c *Client) run() {
	defer c.conn.Close()
	for {
		switch c.status {
		default:
			c.listenMessage()
		case Data:
			c.listenData()
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
	requireNoError(err)
	c.onMessage(msg)
}

func (c *Client) onMessage(msg Message) {
	log.Println("Message received:", msg)
	switch msg.Status {
	case Start:
		c.start(msg)
	default:
		c.error("Wrong message status")
	}
}

func (c *Client) start(msg Message) {
	payload, err := msg.StartPayload()
	requireNoError(err)
	c.req = payload
	c.status = Data
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
	CreateFile(payload.RelPath)
	log.Println("Payload saved, writing status=DATA", payload)
	writeStatus(Data, c.conn)
}

func (c *Client) startDownload(msg StartPayload) {
	// TODO
}

func (c *Client) listenData() {
	chunk := readChunk(c.conn)
	c.processChunk(chunk)
	if c.count == c.req.Size {
		c.status = Eof
		log.Println("File saved, writing status EOF")
		writeStatus(Eof, c.conn)
	}
}

func (c *Client) listenEof() {
	log.Println("Listening for EOF")
	msg, err := readMessage(c.conn)
	requireNoError(err)
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
	WriteBuf(c.req.RelPath, chunk)
	c.count += int64(len(chunk))
}

func (c *Client) eof(msg Message) {
	if msg.Status != Eof {
		c.error("Expecting EOF")
		return
	}
	log.Println("DONE!")
	c.status = Done
	writeStatus(Done, c.conn)
}

func (c *Client) overflows(chunk []byte) bool {
	return c.count+int64(len(chunk)) > c.req.Size
}

func (c *Client) error(msg string) {
	// TODO update func to accept msg
	log.Println("ERROR:", msg)
	c.status = Error
	writeStatus(Error, c.conn)
}
