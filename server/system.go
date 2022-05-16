// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"errors"
)

type State uint

const (
	Start State = iota
	Data
	Stream
	Eof
	Error
	Done
)

func (s State) String() string {
	return States()[s]
}

func ToState(i uint) (State, error) {
	if int(i) >= len(States()) {
		return State(0), errors.New("invalid state")
	}
	return State(i), nil
}

func States() []string {
	return []string{
		"start",
		"data",
		"stream",
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

type Message struct {
	State
	Payload
}

type Payload struct {
	Data []byte
}

func NewPayloadFrom(p any) (Payload, error) {
	ser, err := json.Marshal(p)
	return Payload{Data: ser}, err
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

// StreamPayload Returns the computed attribute for an assumed StreamPayload
// data.
func (p Payload) StreamPayload() (StreamPayload, error) {
	payload := StreamPayload{}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}

type StartPayload struct {
	Action
	FileInfo
}

type StreamPayload struct {
	FileInfo
}

type Channel struct {
	Name string
}
