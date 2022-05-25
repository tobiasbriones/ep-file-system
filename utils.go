// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package fs

import "testing"

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func RequireError(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatal(msg)
	}
}
