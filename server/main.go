// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

// Entry point for the file system server.
//
// Author: Tobias Briones
package main

import (
	"fmt"
	"fs/utils"
	"net"
)

const (
	port    = 8080
	network = "tcp"
	bufSize = 1024
)

func main() {
	server, err := net.Listen(network, getServerAddress())

	defer server.Close()
	utils.RequireNoError(err)
	listen(server)
}

func getServerAddress() string {
	return fmt.Sprintf("0.0.0.0:%v", port)
}
