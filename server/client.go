// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs/process"
	"fs/server/io"
	"log"
	"net"
)

type Client struct {
	conn  net.Conn
	state process.State
	req   process.StartPayload
	count int64
}

func newClient(
	conn net.Conn,
) *Client {
	return &Client{
		conn:  conn,
		state: process.Start,
	}
}

func (c *Client) run() {
	defer c.conn.Close()
	for {
		switch c.state {
		default:
			c.listenMessage()
		case process.Data:
			c.listenData()
		case process.Stream:
			c.listenStream()
		case process.Eof:
			c.listenEof()
		case process.Done, process.Error:
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

func (c *Client) onMessage(msg process.Message) {
	log.Println("Message received:", msg)
	switch msg.State {
	case process.Start:
		c.start(msg)
	default:
		c.error("Wrong message state")
	}
}

func (c *Client) start(msg process.Message) {
	payload, err := msg.StartPayload()
	if err != nil {
		c.error("Fail to read StartPayload")
		return
	}
	err = payload.CreateChannelIfNotExists(payload.Channel.Name)
	if err != nil {
		c.error("Fail to read StartPayload Path/Create channel")
		return
	}
	c.req = payload
	c.count = 0

	switch payload.Action {
	case process.ActionUpload:
		c.startUpload(payload)
	case process.ActionDownload:
		c.startDownload(payload)
	}
}

func (c *Client) startUpload(payload process.StartPayload) {
	if payload.Size <= 0 {
		c.error("File sent is empty")
		return
	}
	err := payload.Create(c.req.Channel.Name)
	if err != nil {
		c.error("Fail to create file")
		return
	}
	log.Println("Payload saved, writing state=DATA", payload)
	c.state = process.Data
	err = writeState(process.Data, c.conn)
	if err != nil {
		c.error("Fail to write state=DATA")
		return
	}
}

func (c *Client) listenData() {
	chunk, _ := readChunk(c.conn)
	c.processChunk(chunk)
	if c.count == c.req.Size {
		c.state = process.Eof
		log.Println("File saved, writing state EOF")
		err := writeState(process.Eof, c.conn)
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
	err := c.req.WriteChunk(c.req.Channel.Name, chunk)
	if err != nil {
		c.error("Fail to write chunk")
		return
	}
	c.count += int64(len(chunk))
}

func (c *Client) eof(msg process.Message) {
	if msg.State != process.Eof {
		c.error("Expecting EOF")
		return
	}
	log.Println("DONE!")
	c.state = process.Done
	err := writeState(process.Done, c.conn)
	if err != nil {
		c.error("Fail to write state=DONE")
		return
	}
}

func (c *Client) startDownload(payload process.StartPayload) {
	exists, err := c.req.Exists(c.req.Channel.Name)
	if err != nil {
		c.error("Fail to read file exists")
		return
	}
	if !exists {
		c.error("Requested file does not exist")
		return
	}
	size, err := c.req.ReadFileSize(c.req.Channel.Name)
	if err != nil {
		c.error("Fail to read file size")
		return
	}
	payload.Action = process.ActionDownload
	info := io.FileInfo{
		RelPath: payload.RelPath,
		Size:    size,
	}
	c.req = payload
	c.writeStreamState(process.StreamPayload{FileInfo: info})
}

func (c *Client) writeStreamState(payload process.StreamPayload) {
	p, err := process.NewPayloadFrom(payload)
	if err != nil {
		c.error("Fail to read payload from StreamPayload")
		return
	}
	msg := process.Message{
		State:   process.Stream,
		Payload: p,
	}
	c.state = process.Stream
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
	if msg.State != process.Stream {
		c.error("Wrong client state, state=STREAM was expected")
		return
	}
	c.stream()
}

func (c *Client) stream() {
	err := c.req.Stream(
		c.req.Channel.Name,
		bufSize,
		func(buf []byte) {
			_, err := c.conn.Write(buf)
			if err != nil {
				// TODO Fix StreamLocalFile paradigm
				c.error("Fail to write chunk")
				return
			}
		},
	)
	if err != nil {
		c.error("Fail to stream file")
		return
	}
	log.Println("File sent to client, changing state to DONE")
	c.state = process.Done
	err = writeState(process.Done, c.conn)
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
	c.state = process.Error
	writeState(process.Error, c.conn)
}
