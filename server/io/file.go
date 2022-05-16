// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

// Package io Models a file system according to
// https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system/tree/main/model
// Author Tobias Briones
package io

import (
	"errors"
	"regexp"
)

const (
	Root           = ""
	Separator      = "/"
	ValidPathRegex = "^$|\\w+/*\\.*-*"
)

type Path struct {
	value string
}

func NewPath(value string) (Path, error) {
	if !isValidPath(value) {
		return Path{}, errors.New("invalid path")
	}
	return Path{value: value}, nil
}

func isValidPath(value string) bool {
	r, _ := regexp.Compile(ValidPathRegex)
	return r.MatchString(value)
}
