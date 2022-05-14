// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package main

import "testing"

func requireNoError(e error) {
	if e != nil {
		panic(e)
	}
}

func requirePassedTest(t *testing.T, e error, msg string) {
	if e != nil {
		t.Fatal(msg, "\n Error:", e)
	}
}
