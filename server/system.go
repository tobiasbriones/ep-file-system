// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

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
