// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs/process"
	"io"
	"log"
	"net"
)

type Client struct {
	conn       net.Conn
	process    process.Process
	id         uint // Current ID assigned by the Hub
	register   chan *Client
	unregister chan *Client
	change     chan UpdatePayload
	quit       chan struct{}
}

func newClient(
	conn net.Conn,
	osFsRoot string,
	register chan *Client,
	unregister chan *Client,
) *Client {
	return &Client{
		conn:       conn,
		process:    process.NewProcess(osFsRoot),
		register:   register,
		unregister: unregister,
		change:     make(chan UpdatePayload),
		quit:       make(chan struct{}),
	}
}

func (c *Client) run() {
	defer c.conn.Close()
	c.connect()
	log.Println("Client connected")
	for {
		select {
		case u := <-c.change:
			c.sendUpdate(u)
		case <-c.quit:
			c.unregister <- c
			return
		default:
			c.next()
		}
	}
}

func (c *Client) connect() {
	c.register <- c
}

func (c *Client) next() {
	switch c.process.State() {
	default:
		c.listenMessage()
	case process.Data:
		c.listenData()
	case process.Stream:
		c.listenStream()
	case process.Eof:
		c.listenEof()
	case process.Error:
		c.sendQuit()
	}
}

func (c *Client) listenMessage() {
	log.Println("Listening for client message")
	msg, err := readMessage(c.conn)
	if err != nil {
		c.handleReadError(err, "fail to read message")
		return
	}
	c.onMessage(msg)
}

func (c *Client) onMessage(msg Message) {
	log.Println("Message received with state:", msg.State)
	switch msg.State {
	case process.Start:
		c.start(msg)
	default:
		if msg.Command != nil {
			c.onCommand(msg.Command)
		} else {
			c.error("wrong message state")
		}
	}
}

func (c *Client) onCommand(cmd map[string]string) {
	req := cmd["REQ"]

	switch req {
	case "LIST_CHANNELS":
		err := writeChannels(c.conn)
		if err != nil {
			c.error("fail to send list of channels")
			return
		}
	default:
		c.error("invalid command request")
	}
}

func (c *Client) start(msg Message) {
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
	chunk, err := readChunk(c.conn)
	if err != nil {
		c.handleReadError(err, "fail to read chunk")
		return
	}
	err = c.process.Data(chunk)
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
		c.handleReadError(err, "fail to read EOF message")
		return
	}
	c.eof(msg)
}

func (c *Client) eof(msg Message) {
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
	p, err := NewPayloadFrom(payload)
	if err != nil {
		c.error("Fail to read payload from StreamPayload")
		return
	}
	msg := Message{
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
		c.handleReadError(err, "fail to read status STREAM")
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

func (c *Client) sendUpdate(u UpdatePayload) {
	p, err := NewPayloadFrom(u)
	if err != nil {
		log.Println(err)
		c.error("Fail to send update")
		return
	}
	msg := Message{
		Payload: p,
	}
	err = writeMessage(msg, c.conn)
	if err != nil {
		log.Println(err)
		c.error("Fail to send update")
		return
	}
}

func (c *Client) handleReadError(err error, msg string) {
	if err == io.EOF {
		log.Println("Communication closed by the client")
		c.sendQuit()
		return
	}
	log.Println(msg, err)
	c.error(msg)
}

func (c *Client) sendQuit() {
	go func() {
		c.quit <- struct{}{}
	}()
}

func (c *Client) error(msg string) {
	log.Println("ERROR:", msg)
	c.process.Error()
	writeErrorState(msg, c.conn)
}
