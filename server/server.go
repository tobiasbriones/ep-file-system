// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"log"
	"net"
)

type Response int

const (
	Connect = iota
	Quit
	Update
	Ok
)

func listen(server net.Listener) {
	osFsRoot := loadRoot()
	hub := NewHub()

	log.Println("Server running on:", osFsRoot)
	go hub.run()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Fail to accept client")
			continue
		}
		client := newClient(
			conn,
			osFsRoot,
			hub.register,
			hub.unregister,
			hub.change,
		)
		go client.run()
	}
}

func loadRoot() string {
	osFsRoot, err := getOsFsRoot()
	if err != nil {
		panic("fail to load OS FS Root")
	}
	return osFsRoot
}
