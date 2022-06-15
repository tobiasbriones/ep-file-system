// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"encoding/json"
	"errors"
	"fs/files"
	"fs/process"
	"log"
	"net"
	"strconv"
)

type req string

const (
	Subscribe                     req = "SUBSCRIBE"
	CreateChannel                 req = "CREATE_CHANNEL"
	ListChannels                  req = "LIST_CHANNELS"
	ListFiles                     req = "LIST_FILES"
	CID                           req = "CID"
	ConnectedUsers                req = "CONNECTED_USERS"
	SubscribeToListConnectedUsers req = "SUBSCRIBE_TO_LIST_CONNECTED_USERS"
)

type command struct {
	conn net.Conn
	commandClient
	clientHubChange chan struct{}
	quit            chan struct{}
}

func newCommand(
	conn net.Conn,
	client commandClient,
	clientHubChange chan struct{},
	quit chan struct{},
) command {
	return command{
		conn:            conn,
		commandClient:   client,
		clientHubChange: clientHubChange,
		quit:            quit,
	}
}

func (c command) execute(cmd map[string]string) error {
	req := req(cmd["REQ"])

	switch req {
	case Subscribe:
		return c.subscribe(cmd)
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
	case SubscribeToListConnectedUsers:
		return c.subscribeToListConnectedUsers()
	default:
		return errors.New("invalid command request")
	}
	return nil
}

func (c command) subscribe(cmd map[string]string) error {
	name := cmd["CHANNEL"]
	c.commandClient.subscribe(process.Channel{Name: name})
	return c.respond(Subscribe, Ok, "")
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
	return c.respond(CreateChannel, Ok, "")
}

func (c command) listChannels() error {
	channels, err := readChannels()
	if err != nil {
		return errors.New("fail to read list of channels")
	}
	ser, _ := json.Marshal(channels)
	return c.respond(ListChannels, Ok, string(ser))
}

func (c command) listFiles(cmd map[string]string) error {
	// TODO channel := c.process.User().Channel()
	channelName := cmd["CHANNEL"]
	channel := process.NewChannel(channelName)

	fileList, err := readFiles(channel)
	if err != nil {
		return errors.New("fail to read list of files")
	}
	ser, _ := json.Marshal(fileList)
	return c.respond(ListFiles, Ok, string(ser))
}

func (c command) sendCID() error {
	payload := strconv.Itoa(int(c.cid()))
	return c.respond(CID, Ok, payload)
}

func (c command) subscribeToListConnectedUsers() error {
	log.Println("Subscribing client to listen for connected users")
	go func() {
		for {
			select {
			case <-c.clientHubChange:
				c.requestClientList()
			case <-c.quit:
				return
			}
		}
	}()
	return nil
}

func (c command) respond(req req, res Response, payload string) error {
	cmd := make(map[string]string)
	cmd["REQ"] = string(req)
	cmd["PAYLOAD"] = payload
	msg := Message{
		Command:  cmd,
		Response: res,
	}
	return writeMessage(msg, c.conn)
}

type commandClient interface {
	cid() uint
	subscribe(channel process.Channel)
	requestClientList()
}

func readChannels() ([]string, error) {
	root, err := getFsRootFile()
	if err != nil {
		return nil, err
	}
	channels, err := files.ReadFileNames(root)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func readFiles(channel process.Channel) ([]string, error) {
	root, err := getFsRootFile()
	if err != nil {
		return nil, err
	}
	dir, _ := channel.File()
	channelFile := dir.ToOsFile(root.Path())
	fileList, err := files.ReadFileNames(channelFile)
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
