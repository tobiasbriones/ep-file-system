// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import (
	"strings"
)

const (
	DefChannel = "main"
	fsRootPath = "fs"
)

func getPath(relPath string, channel string) (Path, error) {
	path, err := NewPathFrom(fsRootPath, channel)
	if err != nil {
		return Path{}, err
	}
	children := strings.Split(relPath, Separator)
	err = path.Append(children...)
	if err != nil {
		return Path{}, err
	}
	return path, nil
}
