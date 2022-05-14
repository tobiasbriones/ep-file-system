// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import "errors"

type Status uint

const (
	Start Status = 0
	Ok    Status = 1
	Data  Status = 2
	Eof   Status = 3
	Error Status = 4
	Done  Status = 5
)

func (s Status) String() string {
	return Statuses()[s]
}

func ToStatus(i uint) (Status, error) {
	if int(i) >= len(Statuses()) {
		return -1, errors.New("invalid status")
	}
	return Status(i), nil
}

func Statuses() []string {
	return []string{
		"start",
		"ok",
		"data",
		"eof",
		"error",
		"done",
	}
}

type Message struct {
	Status
	Action  string
	Payload string
	Data    []byte
}

type MessageType uint

const (
	MsgFileInfo MessageType = 0
	MsgAction   MessageType = 1
	MsgData     MessageType = 2
)

func (t MessageType) String() string {
	return MessageTypes()[t]
}

func ToMessageType(i uint) (MessageType, error) {
	if int(i) >= len(MessageTypes()) {
		return -1, errors.New("invalid message type")
	}
	return MessageType(i), nil
}

func MessageTypes() []string {
	return []string{
		"file-info",
		"action",
		"data",
	}
}
