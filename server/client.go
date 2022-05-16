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

}

func (c *Client) listenData() {

}

func (c *Client) listenEof() {

}

func (c *Client) error(msg string) {
	// TODO update func to accept msg
	writeStatus(Error, c.conn)
}
