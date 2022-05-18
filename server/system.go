// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import (
	"encoding/json"
	"errors"
	"server/io"
)

type State uint

const (
	Start State = iota
	Data
	Stream
	Eof
	Error
	Done
	// Update This state will be used to send broadcast notifications to
	// clients. It is not related to the main FSM.
	Update
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

func (p Payload) UpdatePayload() (UpdatePayload, error) {
	payload := UpdatePayload{}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}

type StartPayload struct {
	Action
	io.FileInfo
	Channel Channel
}

type StreamPayload struct {
	io.FileInfo
}

type UpdatePayload struct {
	change bool // Rudimentary signal to test broadcast
}

type Channel struct {
	Name string
}

func NewChannel(name string) Channel {
	return Channel{Name: name}
}
