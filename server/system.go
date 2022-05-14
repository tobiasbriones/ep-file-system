// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

type Status string

const (
	START Status = "start"
	OK    Status = "ok"
	DATA  Status = "data"
	EOF   Status = "eof"
	ERROR Status = "error"
)

type Message struct {
	Status
	Action  string
	Payload string
	Data    []byte
}
