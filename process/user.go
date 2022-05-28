// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package process

import (
	"errors"
	"fs"
	"fs/files"
	"log"
)

// User Contains all the FSM implementation details.
type User struct {
	file     fs.OsFile
	channel  Channel
	osFsRoot string
	size     uint64
	count    int64
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

func (u *User) start(payload StartPayload) error {
	u.channel = payload.Channel
	u.count = 0
	err := u.setFile()
	if err != nil {
		return err
	}
	err = u.createChannelIfNotExists()
	if err != nil {
		return err
	}
	err = u.startAction(payload)
	if err != nil {
		return err
	}
	return nil
}

func (u User) setFile() error {
	file, err := u.getOsFile()
	if err != nil {
		return err
	}
	u.file = file
	return nil
}

func (u User) startAction(payload StartPayload) error {
	switch payload.Action {
	case ActionUpload:
		err := u.startActionUpload(payload)
		if err != nil {
			return err
		}
	case ActionDownload:
		err := u.startActionDownload()
		if err != nil {
			return err
		}
	}
	return nil
}

func (u User) startActionUpload(payload StartPayload) error {
	u.size = payload.Size
	if u.size <= 0 {
		return errors.New("file sent is empty")
	}
	err := u.createFile()
	if err != nil {
		log.Println(err)
		return errors.New("fail to create file")
	}
	return nil
}

func (u User) startActionDownload() error {
	exists, err := files.Exists(u.file)
	if err != nil {
		log.Println(err)
		return errors.New("fail to read file exists")
	}
	if !exists {
		log.Println(err)
		return errors.New("requested file does not exist")
	}
	size, err := files.ReadSize(u.file)
	if err != nil {
		log.Println(err)
		return errors.New("fail to read file size")
	}
	u.size = uint64(size)
	return nil
}

func (u User) getOsFile() (fs.OsFile, error) {
	fsFile, err := fs.NewFileFromString(u.channel.Name) // channel/
	if err != nil {
		log.Println(err)
		return fs.OsFile{}, errors.New("invalid channel name: " + u.channel.Name)
	}
	err = fsFile.Append(u.file.Value) // channel/file.txt
	if err != nil {
		log.Println(err)
		return fs.OsFile{}, errors.New("invalid file: " + u.file.Value)
	}
	return fsFile.ToOsFile(u.osFsRoot), nil
}

func (u User) createChannelIfNotExists() error {
	channel, err := u.channel.File()
	if err != nil {
		log.Println(err)
		return errors.New("invalid channel")
	}
	channelFile := channel.ToOsFile(u.osFsRoot)
	err = files.CreateIfNotExists(channelFile)
	if err != nil {
		log.Println(err)
		return errors.New("fail to read StartPayload Path/Create channel")
	}
	return nil
}

func (u User) createFile() error {
	return files.Create(u.file)
}

func (u User) processChunk(chunk []byte) error {
	if u.overflows(chunk) {
		return errors.New("overflow")
	}
	if len(chunk) == 0 {
		return errors.New("underflow")
	}
	err := files.WriteBuf(u.file, chunk)
	if err != nil {
		log.Println(err)
		return errors.New("fail to write chunk")
	}
	u.count += int64(len(chunk))
	return nil
}

func (u User) overflows(chunk []byte) bool {
	return u.count+int64(len(chunk)) > int64(u.size)
}

func (u User) stream(size uint, f func(buf []byte)) error {
	err := files.Stream(u.file, size, f)
	if err != nil {
		return err
	}
	return nil
}
