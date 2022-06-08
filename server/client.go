// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"encoding/json"
	"fs/process"
	"log"
	"net"
)

type Client struct {
	conn       net.Conn
	command    command
	state      state
	id         uint // Current ID assigned by the Hub
	register   chan *Client
	unregister chan *Client
	notify     chan UpdatePayload
	list       chan *Client
	quit       chan struct{}
}

func newClient(
	conn net.Conn,
	osFsRoot string,
	register chan *Client,
	unregister chan *Client,
	change chan struct{},
	list chan *Client,
) *Client {
	client := &Client{
		conn:       conn,
		register:   register,
		unregister: unregister,
		list:       list,
		notify:     make(chan UpdatePayload),
		quit:       make(chan struct{}),
	}
	client.command = newCommand(client.conn, client)
	client.state = newState(client.conn, osFsRoot, client.sendQuit, change)
	return client
}

func (c *Client) run() {
	defer c.conn.Close()
	c.connect() // TODO synchronize, wait for completing signal register
	log.Println("Client connected")

	go c.runNotification()
	for {
		select {
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

func (c *Client) runNotification() {
	for {
		select {
		case u := <-c.notify:
			c.sendUpdate(u)
		case <-c.quit:
			return
		}
	}
}

func (c *Client) next() {
	if c.state.isInProgress() {
		c.state.next()
	} else {
		c.listenMessage()
	}
}

func (c *Client) listenMessage() {
	log.Println("Listening for client message")
	msg, err := readMessage(c.conn, longReadTimeOut)
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
		c.state.start(msg)
	default:
		c.handleCommand(msg)
	}
}

func (c *Client) handleCommand(msg Message) {
	if msg.Command != nil {
		err := c.command.execute(msg.Command)
		if err != nil {
			writeErrorState(err.Error(), c.conn)
		}
	} else {
		writeErrorState("wrong message state", c.conn)
	}
}

func (c *Client) sendUpdate(u UpdatePayload) {
	if c.state.isInProgress() {
		return
	}
	p, err := NewPayloadFrom(u)
	if err != nil {
		log.Println(err)
		c.state.error("Fail to send update")
		return
	}
	msg := Message{
		Response: Update,
		Payload:  p,
	}
	err = writeMessage(msg, c.conn)
	if err != nil {
		log.Println(err)
		c.state.error("Fail to send update")
		return
	}
}

func (c *Client) sendList(clients []string) {
	enc := json.NewEncoder(c.conn)
	enc.Encode(clients)
}

func (c *Client) handleReadError(err error, msg string) {
	log.Println(msg, err)
	c.sendQuit()
}

func (c *Client) sendQuit() {
	go func() {
		c.quit <- struct{}{}
	}()
}

func (c Client) cid() uint {
	return c.id
}

func (c *Client) requestClientList() {
	c.list <- c
}
