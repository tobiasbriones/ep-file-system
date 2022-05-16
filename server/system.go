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
	Payload
}

type Payload struct {
	Data []byte
}

func NewPayload(v any) (Payload, error) {
	payload, err := json.Marshal(v)
	return Payload{Data: payload}, err
}

// StartPayload Returns the computed attribute for an assumed StartPayload
// data.
func (p Payload) StartPayload() (StartPayload, error) {
	payload := StartPayload{}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}

type StartPayload struct {
	Action
	FileInfo
}

type Status uint

const (
	Start Status = iota
	Ok
	Data
	Eof
	Error
	Done
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
	ActionUpload Action = iota
	ActionDownload
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
