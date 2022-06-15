// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"log"
	"strconv"
)

type Hub struct {
	clients         map[uint]*Client
	register        chan *Client
	unregister      chan *Client
	quit            chan struct{}
	change          chan struct{} // Rudimentary signal to test broadcast
	list            chan *Client  // Signal to send the list of connected clients
	cid             uint          // Current ID for clients on this server instance
	clientHubChange chan struct{} // When a client regs or unregs
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[uint]*Client),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		quit:            make(chan struct{}),
		change:          make(chan struct{}),
		list:            make(chan *Client),
		cid:             0,
		clientHubChange: make(chan struct{}),
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
			go h.broadcastChange()
		case c := <-h.list:
			go h.listClients(c)
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
	go func() {
		h.clientHubChange <- struct{}{}
	}()
	log.Println("Registering client", id, "into the Hub")
}

func (h *Hub) unregisterAll() {
	for _, client := range h.clients {
		h.unregisterClient(client)
	}
}

func (h *Hub) unregisterClient(c *Client) {
	delete(h.clients, c.id)
	go func() {
		h.clientHubChange <- struct{}{}
	}()
	log.Println("Unregistering client", c.id, "from the Hub")
}

func (h *Hub) broadcastChange() {
	payload := UpdatePayload{Change: true}
	for _, client := range h.clients {
		client.notify <- payload
	}
}

func (h *Hub) listClients(c *Client) {
	var list []string
	for _, client := range h.clients {
		list = append(list, strconv.Itoa(int(client.id)))
	}
	c.sendList(list)
}
