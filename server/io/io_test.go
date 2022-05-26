// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package io

import (
	"fs/utils"
	"testing"
)

func TestGetPath(t *testing.T) {
	path, err := getPath("file.txt", DefChannel)
	utils.RequirePassCase(t, err, "")
	if path.Value != "main/file.txt" {
		t.Fatal("Computed path is wrong")
	}

	path, err = getPath("dir1/dir2/file.txt", DefChannel)
	utils.RequirePassCase(t, err, "")
	if path.Value != "main/dir1/dir2/file.txt" {
		t.Fatal("Computed path is wrong")
	}
}
