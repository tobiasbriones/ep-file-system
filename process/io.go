// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package process

import (
	"fs"
)

type StartPayload struct {
	Action
	fs.FileInfo
	Channel Channel
}

type StreamPayload struct {
	fs.FileInfo
}
