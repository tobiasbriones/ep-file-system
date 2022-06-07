// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"encoding/json"
	"fs/utils"
	"log"
	"net"
	"testing"
	"time"
)

func TestCommandId(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "CID"
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	buf := make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal("Fail to read CID: ", err)
	}

	// Check
	log.Println("Client ID:", string(buf[:n]))
}

func TestCommandListOfClients(t *testing.T) {
	runOnHold := func() {
		tcpAddr, _ := net.ResolveTCPAddr(network, getServerAddress())
		net.DialTCP(network, nil, tcpAddr)
		time.Sleep(2 * time.Second)
	}
	go runOnHold()
	go runOnHold()
	go runOnHold()

	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "CONNECTED_USERS"
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	var users []string
	dec := json.NewDecoder(conn)
	err = dec.Decode(&users)
	if err != nil {
		t.Fatal("Fail to read CONNECTED_USERS: ", err)
	}

	// Check
	log.Println("Connected users:", users)
}
