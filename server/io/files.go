// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import (
	"errors"
	"fs"
	"os"
)

func Exists(file fs.File) (bool, error) {
	path, err := AbsolutePath(file)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, nil
}

func Create(file fs.File) error {
	path, err := AbsolutePath(file)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	f.Close()
	return err
}

func CreateIfNotExists(file fs.File) error {
	path, err := AbsolutePath(file)
	if err != nil {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

func ReadSize(file fs.File) (int64, error) {
	path, err := AbsolutePath(file)
	if err != nil {
		return 0, err
	}
	return ReadFileSize(path)
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
