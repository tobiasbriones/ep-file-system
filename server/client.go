// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs/process"
	"log"
	"net"
)

type Client struct {
	conn    net.Conn
	process process.Process
}

func newClient(
	conn net.Conn,
	osFsRoot string,
) *Client {
	return &Client{
		conn:    conn,
		process: process.NewProcess(osFsRoot),
	}
}

func (c *Client) run() {
	defer c.conn.Close()
	log.Println("Client connected")
	for {
		switch c.process.State() {
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
		log.Println("Fail to read message:", err)
		c.error("fail to read message")
		return
	}
	c.onMessage(msg)
}

func (c *Client) onMessage(msg process.Message) {
	log.Println("Message received with state:", msg.State)
	switch msg.State {
	case process.Start:
		c.start(msg)
	default:
		c.error("wrong message state")
	}
}

func (c *Client) start(msg process.Message) {
	payload, err := msg.StartPayload()
	if err != nil {
		c.error("fail to read StartPayload")
		return
	}
	err = c.process.Start(payload)
	if err != nil {
		c.error(err.Error())
		return
	}
	c.onProcessStarted()
}

func (c *Client) onProcessStarted() {
	switch c.process.Action() {
	case process.ActionUpload:
		c.onActionUploadStarted()
	case process.ActionDownload:
		c.onActionDownloadStarted()
	}
}

func (c *Client) onActionUploadStarted() {
	err := writeState(process.Data, c.conn)
	if err != nil {
		c.error("Fail to write state=DATA")
		return
	}
}

func (c *Client) onActionDownloadStarted() {
	c.writeStreamState(process.StreamPayload{FileInfo: c.process.User().FileInfo()})
}

func (c *Client) listenData() {
	chunk, _ := readChunk(c.conn)
	err := c.process.Data(chunk)
	if err != nil {
		c.error(err.Error())
		return
	}
	c.onChunkProcessed()
}

func (c *Client) onChunkProcessed() {
	if c.process.State() == process.Eof {
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
		c.error("fail to read message")
		return
	}
	c.eof(msg)
}

func (c *Client) eof(msg process.Message) {
	if msg.State != process.Eof {
		c.error("expecting EOF")
		return
	}
	log.Println("DONE!")
	err := c.process.Done()
	if err != nil {
		c.error("fail to write state=DONE on server")
		return
	}
	err = writeState(process.Done, c.conn)
	if err != nil {
		c.error("fail to write state=DONE")
		return
	}
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
		c.error("fail to read message")
		return
	}
	if msg.State != process.Stream {
		c.error("wrong client state, state=STREAM was expected")
		return
	}
	c.stream()
}

func (c *Client) stream() {
	err := c.process.Stream(
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
		c.error("fail to stream file: " + err.Error())
		return
	}
	log.Println("File sent to client, changing state to DONE")
	err = writeState(process.Done, c.conn)
	if err != nil {
		c.error("Fail to write state=DONE")
		return
	}
}

func (c *Client) error(msg string) {
	// TODO update func to accept msg
	log.Println("ERROR:", msg)
	c.process.Error()
	writeState(process.Error, c.conn)
}
