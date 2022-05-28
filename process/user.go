// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package process

import (
	"errors"
	"fs"
	"fs/files"
)

// User Contains all the FSM implementation details.
type User struct {
	file     fs.OsFile
	channel  Channel
	osFsRoot string
	size     uint64
	count    uint64
}

func newUser(osFsRoot string) User {
	return User{
		osFsRoot: osFsRoot,
	}
}

func (u User) FileInfo() fs.FileInfo {
	return fs.FileInfo{
		File: u.file.File,
		Size: u.size,
	}
}

func (u User) File() fs.OsFile {
	return u.file
}

func (u User) Size() uint64 {
	return u.size
}

func (u *User) set(payload StartPayload) error {
	u.channel = payload.Channel
	file, err := u.getOsFile()
	if err != nil {
		return err
	}
	u.file = file
	u.count = 0
	return nil
}

func (u User) getOsFile() (fs.OsFile, error) {
	fsFile, err := fs.NewFileFromString(u.channel.Name) // channel/
	if err != nil {
		return fs.OsFile{}, errors.New("invalid channel name: " + u.channel.Name)
	}
	err = fsFile.Append(u.file.Value) // channel/file.txt
	if err != nil {
		return fs.OsFile{}, errors.New("invalid file: " + u.file.Value)
	}
	return fsFile.ToOsFile(u.osFsRoot), nil
}

func (u User) init() error {
	return u.createChannelIfNotExists()
}

func (u User) createChannelIfNotExists() error {
	channel, err := u.channel.File()
	if err != nil {
		return errors.New("invalid channel")
	}
	channelFile := channel.ToOsFile(u.osFsRoot)
	err = files.CreateIfNotExists(channelFile)
	if err != nil {
		return errors.New("fail to read StartPayload Path/Create channel")
	}
	return nil
}
