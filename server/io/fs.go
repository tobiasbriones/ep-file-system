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
	"strings"
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

func NewFileFromString(value string) (File, error) {
	path, err := NewPath(value)
	return File{Path: path}, err
}

// Directory is just a simple Path for this system.
// It's open to extension with more properties.
type Directory struct {
	Path
}

func NewDirectoryFromString(value string) (Directory, error) {
	path, err := NewPath(value)
	return Directory{Path: path}, err
}

type Path struct {
	Value string
}

func (p *Path) Append(values ...string) error {
	end, err := NewPathFrom(values...)
	if err != nil {
		return err
	}
	var newValue string
	if p.IsRoot() {
		newValue = end.Value
	} else {
		newValue = p.Value + Separator + end.Value
	}
	p.Value = newValue
	return nil
}

func (p *Path) IsRoot() bool {
	return p.Value == Root
}

// NewPathFrom constructs a Path from the given tokens. Tokens must be
// independent, e.g. not containing the separator character, one at a time.
func NewPathFrom(values ...string) (Path, error) {
	str := ""
	for _, value := range values {
		if strings.Contains(value, Separator) {
			msg := "invalid path token, it contains the separator character"
			return Path{}, errors.New(msg)
		}
		str += value + Separator
	}
	// Remove last separator
	str = str[:len(str)-1]
	return NewPath(str)
}

func NewPath(value string) (Path, error) {
	if !isValidPath(value) {
		return Path{}, errors.New("invalid path")
	}
	return Path{Value: value}, nil
}

func isValidPath(value string) bool {
	r, _ := regexp.Compile(ValidPathRegex)
	return r.MatchString(value)
}
