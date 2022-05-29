// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package utils

import "testing"

func RequirePassCase(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatal(msg, "Error:", err.Error())
	}
}

func RequireFailureCase(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatal(msg)
	}
}

func RequireNoError(e error) {
	if e != nil {
		panic(e)
	}
}
