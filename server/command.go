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

type command struct {
	conn net.Conn
	commandClient
}

func newCommand(conn net.Conn, client commandClient) command {
	return command{conn: conn, commandClient: client}
}

func (c command) execute(cmd map[string]string) error {
	req := cmd["REQ"]

	switch req {
	case "CREATE_CHANNEL":
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
	case "LIST_CHANNELS":
		err := writeChannels(c.conn)
		if err != nil {
			return errors.New("fail to send list of channels")
		}
	case "LIST_FILES":
		// TODO channel := c.process.User().Channel()
		channelName := cmd["CHANNEL"]
		channel := process.NewChannel(channelName)
		err := writeFiles(c.conn, channel)
		if err != nil {
			return errors.New("fail to send list of files")
		}
	case "CID":
		_, err := c.conn.Write([]byte(strconv.Itoa(int(c.cid())) + "\n"))
		if err != nil {
			return errors.New("fail to send client ID")
		}
	case "CONNECTED_USERS":
		// Send a signal to send the list of users to this client
		c.requestClientList()
	default:
		return errors.New("invalid command request")
	}
	return nil
}

type commandClient interface {
	cid() uint
	requestClientList()
}
