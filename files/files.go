// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package io

import (
	"errors"
	"fs"
	"os"
	"path/filepath"
	"strings"
)

func Exists(file fs.OsFile) (bool, error) {
	if _, err := os.Stat(file.Path()); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, nil
}

func Create(file fs.OsFile) error {
	f, err := os.Create(file.Path())
	f.Close()
	return err
}

func CreateIfNotExists(file fs.OsFile) error {
	return os.MkdirAll(file.Path(), os.ModePerm)
}

func ReadSize(file fs.OsFile) (int64, error) {
	return ReadFileSize(file.Path())
}

func ReadFileSize(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func GetExecPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath := filepath.Dir(ex)
	exePath = strings.ReplaceAll(exePath, "\\", fs.Separator)
	return exePath, nil
}
