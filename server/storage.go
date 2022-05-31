// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs"
	"fs/files"
)

const (
	fsRoot = ".fs"
)

func getFsRootFile() (fs.OsFile, error) {
	path, err := files.GetExecPath()
	if err != nil {
		return fs.OsFile{}, err
	}
	f, _ := fs.NewFileFromString(fsRoot)
	return f.ToOsFile(path), nil
}

func getOsFsRoot() (string, error) {
	path, err := files.GetExecPath()
	if err != nil {
		return "", err
	}
	return path + fs.Separator + fsRoot, nil
}
