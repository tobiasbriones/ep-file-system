// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
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
