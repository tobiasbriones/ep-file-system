// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"errors"
)

type Message struct {
	Status
	Payload []byte
}

func (m Message) StartPayload() (StartPayload, error) {
	if m.Status != Start {
		return StartPayload{}, errors.New("message status is not START")
	}
	payload := StartPayload{}
	err := json.Unmarshal(m.Payload, &payload)
	return payload, err
}

type StartPayload struct {
	Action
	FileInfo
}

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
		return Status(0), errors.New("invalid status")
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

type Action uint

const (
	ActionUpload   Action = 0
	ActionDownload Action = 1
)

func ToAction(i uint) (Action, error) {
	if int(i) >= len(Actions()) {
		return Action(0), errors.New("invalid action")
	}
	return Action(i), nil
}

func Actions() []string {
	return []string{
		"upload",
		"download",
	}
}
