// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

// Package io(utils) implements basic OS file operations related to this IO
// module.
// Author Tobias Briones

import (
	"fs"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(path string) error {
	f, err := os.Create(path)
	f.Close()
	return err
}

func WriteBuf(path string, buf []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	f.Close()
	return err
}

func MakeDirIfNotExists(path fs.Path) error {
	return os.MkdirAll(path.Value, os.ModePerm)
}

func getExecPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath := filepath.Dir(ex)
	exePath = strings.ReplaceAll(exePath, "\\", fs.Separator)
	return exePath, nil
}
