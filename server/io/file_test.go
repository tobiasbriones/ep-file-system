// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import "testing"

func TestBasics(t *testing.T) {
	if Root != "" {
		t.Fatal("Invalid definition of root path")
	}
	if Separator != "/" {
		t.Fatal("Invalid definition of path separator")
	}
}

func TestNewPath(t *testing.T) {
	_, err := NewPath("")
	requireNoError(t, err)

	// Notice how everything is relative (no initial "/")
	// Although the path "/" is also valid
	_, err = NewPath("fs")
	requireNoError(t, err)

	_, err = NewPath("fs/file-1.txt")
	requireNoError(t, err)
}

func TestNewFileAndDirectory(t *testing.T) {
	// There is no difference between File and Directory so far...

	_, err := NewDirectoryFromString("")
	requireNoError(t, err)

	_, err = NewDirectoryFromString("fs")
	requireNoError(t, err)

	_, err = NewFileFromString("fs/file-1.txt")
	requireNoError(t, err)
}

func requireNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}
