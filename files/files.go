// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package files

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
	f, err := os.Open(file.Path())
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

func WriteBuf(file fs.OsFile, buf []byte) error {
	path := file.Path()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	f.Close()
	return err
}
