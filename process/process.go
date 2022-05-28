// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package process

import (
	"errors"
	"fs"
)

var Valid = struct{}{}

type State string

const (
	Start  State = "START"
	Data   State = "DATA"
	Stream State = "STREAM"
	Eof    State = "EOF"
	Error  State = "ERROR"
	Done   State = "DONE"
)

var stateStrings = map[string]struct{}{
	"start":  Valid,
	"data":   Valid,
	"stream": Valid,
	"eof":    Valid,
	"error":  Valid,
	"done":   Valid,
}

func ToState(value string) (State, error) {
	if _, isValid := stateStrings[value]; !isValid {
		return "", errors.New("invalid state value: " + value)
	}
	return State(value), nil
}

// TODO Update
const (
	// Connect Next states are not related to the main FSM. Initiates the
	// server/client connection, it's sent by the client.
	Connect = iota
	// Quit Sent by a client to exit.
	Quit
	// Update This state will be used to send broadcast notifications to
	// clients.
	Update
	// Ok Sent by the server to confirm a client request.
	Ok
)

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

type Process struct {
	state  State
	action Action
	file   fs.OsFile
	size   uint64
}

func NewProcess(action Action, info fs.FileInfo, fsOsRoot string) Process {
	return Process{
		state:  Start,
		action: action,
		file:   info.ToOsFile(fsOsRoot),
		size:   info.Size,
	}
}

type Channel struct {
	Name string
}

func NewChannel(name string) Channel {
	return Channel{Name: name}
}

func (c Channel) File() (fs.File, error) {
	return fs.NewFileFromString(c.Name)
}
