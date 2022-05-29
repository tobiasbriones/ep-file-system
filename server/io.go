// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"encoding/json"
	"fs"
	"fs/process"
)

type Message struct {
	Response
	process.State
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
func (p Payload) StartPayload() (process.StartPayload, error) {
	payload := process.StartPayload{}
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

func (p Payload) ErrorPayload() (ErrorPayload, error) {
	payload := ErrorPayload{}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}

type StreamPayload struct {
	fs.FileInfo
}

type UpdatePayload struct {
	Change bool // Rudimentary signal to test broadcast
}

type ErrorPayload struct {
	Message string
}
