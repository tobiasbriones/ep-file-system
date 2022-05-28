// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package process

import "fs"

// User Contains all the FSM implementation details.
type User struct {
	file     fs.OsFile
	channel  Channel
	osFsRoot string
	size     uint64
	count    uint64
}
