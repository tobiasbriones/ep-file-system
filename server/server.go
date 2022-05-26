// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"log"
	"net"
)

func listen(server net.Listener) {
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Fail to accept client")
			continue
		}
		client := newClient(conn)
		go client.run()
	}
}
