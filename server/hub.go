// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	change     chan bool // Rudimentary signal to test broadcast
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		change:     make(chan bool),
	}
}
