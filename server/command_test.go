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
	data := readData(t, conn)

	// Check
	res := Message{}
	_ = json.Unmarshal(data, &res)

	if res.Response != Ok {
		t.Fatal("Response was not OK")
	}
	if res.Command["REQ"] != "CID" {
		t.Fatal("Invalid request response")
	}
	log.Println("Client ID:", res.Command["PAYLOAD"])
}

func TestCommandListOfClients(t *testing.T) {
	// TODO EXPERIMENTAL, don't use it yet, returns raw array instead of Message
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
	data := readData(t, conn)
	res := Message{}
	_ = json.Unmarshal(data, &res)
	log.Println(res)
	if res.Response != Ok {
		t.Fatal("Response was not OK")
	}
	// It used to be CONNECTED_USERS but both are about the same
	// TODO it should be well defined
	if res.Command["REQ"] != "SUBSCRIBE_TO_LIST_CONNECTED_USERS" {
		t.Fatal("Invalid request response")
	}

	// Check
	var users []string
	err = json.Unmarshal([]byte(res.Command["PAYLOAD"]), &users)

	if err != nil {
		t.Fatal("Fail to read payload")
	}

	// Check
	log.Println("Connected users:", users)
}

func TestCommandChannelList(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "LIST_CHANNELS"
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	data := readData(t, conn)
	res := Message{}
	_ = json.Unmarshal(data, &res)

	if res.Response != Ok {
		t.Fatal("Response was not OK")
	}
	if res.Command["REQ"] != "LIST_CHANNELS" {
		t.Fatal("Invalid request response")
	}

	// Check
	var channels []string
	err = json.Unmarshal([]byte(res.Command["PAYLOAD"]), &channels)

	if err != nil {
		t.Fatal("Fail to read payload")
	}

	// Check at least has the main, and test channels
	if !utils.StringSliceContains(channels, "main") {
		t.Fatal("Channel list does not contain channel: main")
	}
	if !utils.StringSliceContains(channels, "test") {
		t.Fatal("Channel list does not contain channel: test")
	}
	log.Println("Channels:", channels)
}

func TestCommandFileList(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr(network, getServerAddress())
	utils.RequirePassCase(t, err, "Fail to resolve TCP address")
	conn, err := net.DialTCP(network, nil, tcpAddr)
	defer conn.Close()
	utils.RequirePassCase(t, err, "Fail to establish connection")

	// Send command
	cmd := make(map[string]string)
	cmd["REQ"] = "LIST_FILES"
	cmd["CHANNEL"] = testChannel
	msg := Message{
		Command: cmd,
	}
	b, err := json.Marshal(msg)
	_, err = conn.Write(b)
	utils.RequirePassCase(t, err, "Fail to write command to the server")

	// Receive response
	data := readData(t, conn)
	res := Message{}
	_ = json.Unmarshal(data, &res)

	if res.Response != Ok {
		t.Fatal("Response was not OK")
	}
	if res.Command["REQ"] != "LIST_FILES" {
		t.Fatal("Invalid request response")
	}

	// Check
	var fileList []string
	err = json.Unmarshal([]byte(res.Command["PAYLOAD"]), &fileList)

	if err != nil {
		t.Fatal("Fail to read payload")
	}

	// Check at least has the file.pdf file (required for side effect tests)
	if !utils.StringSliceContains(fileList, "file.pdf") {
		t.Fatal("Channel does not contain file: file.pdf")
	}
	log.Println("Files on channel test:", fileList)
}

func readData(t *testing.T, conn net.Conn) []byte {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal("Fail to read CID: ", err)
	}
	return buf[:n]
}
