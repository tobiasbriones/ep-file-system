// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import (
	"fs"
	"strings"
)

const (
	fsRootPath = ".fs"
)

func getChannelPath(channel string) (fs.Path, error) {
	// TODO File needs to implement Parent()
	path, err := fs.NewPathFrom(fsRootPath, channel)
	if err != nil {
		return fs.Path{}, err
	}
	return path, nil
}

func getPath(relPath string, channel string) (fs.Path, error) {
	path, err := fs.NewPathFrom(channel)
	if err != nil {
		return fs.Path{}, err
	}
	children := strings.Split(relPath, fs.Separator)
	err = path.Append(children...)
	if err != nil {
		return fs.Path{}, err
	}
	return path, nil
}

func AbsolutePath(file fs.File) (string, error) {
	root, err := getExecPath()
	if err != nil {
		return "", err
	}
	return root + fs.Separator + fsRootPath + fs.Separator + file.Value, nil
}
