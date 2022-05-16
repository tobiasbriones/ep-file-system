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

func TestNewComposedPath(t *testing.T) {
	composed, err := NewPathFrom(Root)
	requireNoError(t, err)

	if composed.value != Root {
		t.Fatal("Wrong root composed path")
	}

	composed, err = NewPathFrom(
		"fs",
		"dir",
	)
	requireNoError(t, err)

	if composed.value != "fs/dir" {
		t.Fatal("Wrong composed path")
	}

	composed, err = NewPathFrom(
		"fs",
		"dir",
		"file.txt",
	)
	requireNoError(t, err)

	if composed.value != "fs/dir/file.txt" {
		t.Fatal("Wrong composed path")
	}

	composed, err = NewPathFrom(
		"fs",
		"/dir",
		"file.txt",
	)
	requireError(
		t,
		err,
		"Composed paths must not have tokens containing the separator char",
	)
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

func requireError(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatal(msg)
	}
}
