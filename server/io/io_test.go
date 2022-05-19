// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package io

import "testing"

func TestGetPath(t *testing.T) {
	path, err := getPath("file.txt", DefChannel)
	requireNoError(t, err)
	if path.value != "fs/main/file.txt" {
		t.Fatal("Computed path is wrong")
	}

	path, err = getPath("dir1/dir2/file.txt", DefChannel)
	requireNoError(t, err)
	if path.value != "fs/main/dir1/dir2/file.txt" {
		t.Fatal("Computed path is wrong")
	}
}
