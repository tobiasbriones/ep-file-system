// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import "log"

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	quit       chan struct{}
	change     chan struct{} // Rudimentary signal to test broadcast
	cid        uint          // Current ID for clients on this server instance
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		quit:       make(chan struct{}),
		change:     make(chan struct{}),
		cid:        0,
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.registerClient(c)
		case c := <-h.unregister:
			h.unregisterClient(c)
		case <-h.change:
			h.broadcastChange()
		case <-h.quit:
			h.unregisterAll()
			return
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	id := h.cid
	client.id = id
	h.clients[id] = client
	h.cid++
	log.Println("Registering client", id, "into the Hub")
}

func (h *Hub) unregisterAll() {
	for _, client := range h.clients {
		h.unregisterClient(client)
	}
}

func (h *Hub) unregisterClient(c *Client) {
	delete(h.clients, c.id)
	log.Println("Unregistering client", c.id, "from the Hub")
}

func (h *Hub) broadcastChange() {
	payload := UpdatePayload{Change: true}
	for _, client := range h.clients {
		client.change <- payload
	}
}
