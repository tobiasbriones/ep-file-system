// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

// Package io Models a file system according to
// https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system/tree/main/model
// It doesn't have to be too granular, as long as it can read the format of this
// file system.
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

// CommonFile Defines a generic file sum type: File or Directory.
type CommonFile interface{}

// File is just a simple Path for this system.
// It's open to extension with more properties.
type File struct {
	Path
}

// Directory is just a simple Path for this system.
// It's open to extension with more properties.
type Directory struct {
	Path
}

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
