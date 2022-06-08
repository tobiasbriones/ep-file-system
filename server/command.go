// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"errors"
	"fs/files"
	"fs/process"
	"log"
	"net"
	"strconv"
)

type req string

const (
	CreateChannel  req = "CREATE_CHANNEL"
	ListChannels   req = "LIST_CHANNELS"
	ListFiles      req = "LIST_FILES"
	CID            req = "CID"
	ConnectedUsers req = "CONNECTED_USERS"
)

type command struct {
	conn net.Conn
	commandClient
}

func newCommand(conn net.Conn, client commandClient) command {
	return command{conn: conn, commandClient: client}
}

func (c command) execute(cmd map[string]string) error {
	req := req(cmd["REQ"])

	switch req {
	case CreateChannel:
		return c.createChannel(cmd)
	case ListChannels:
		return c.listChannels()
	case ListFiles:
		return c.listFiles(cmd)
	case CID:
		return c.sendCID()
	case ConnectedUsers:
		// Send a signal to send the list of users to this client
		c.requestClientList()
	default:
		return errors.New("invalid command request")
	}
	return nil
}

func (c command) createChannel(cmd map[string]string) error {
	channelName := cmd["CHANNEL"]
	file, err := getFsRootFile()
	if err != nil {
		log.Println(err)
		return errors.New("server error")
	}
	err = file.Append(channelName)
	if err != nil {
		return errors.New("invalid channel")
	}
	err = files.CreateIfNotExists(file)
	if err != nil {
		log.Println(err)
		return errors.New("server error")
	}
	return nil
}

func (c command) listChannels() error {
	err := writeChannels(c.conn)
	if err != nil {
		return errors.New("fail to send list of channels")
	}
	return nil
}

func (c command) listFiles(cmd map[string]string) error {
	// TODO channel := c.process.User().Channel()
	channelName := cmd["CHANNEL"]
	channel := process.NewChannel(channelName)
	err := writeFiles(c.conn, channel)
	if err != nil {
		return errors.New("fail to send list of files")
	}
	return nil
}

func (c command) sendCID() error {
	_, err := c.conn.Write([]byte(strconv.Itoa(int(c.cid())) + "\n"))
	if err != nil {
		return errors.New("fail to send client ID")
	}
	return nil
}

type commandClient interface {
	cid() uint
	requestClientList()
}
