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
	req      req
	file     fs.OsFile
	osFsRoot string
	count    int64
}

func newUser(osFsRoot string) User {
	return User{
		osFsRoot: osFsRoot,
	}
}

func (u User) FileInfo() fs.FileInfo {
	return u.req.info
}

func (u User) File() fs.OsFile {
	return u.file
}

func (u *User) start(payload StartPayload) error {
	u.req.set(payload)
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

func (u *User) setFile() error {
	file, err := u.req.file()
	if err != nil {
		return err
	}
	u.file = file.ToOsFile(u.osFsRoot)
	return nil
}

func (u User) startAction(payload StartPayload) error {
	switch payload.Action {
	case ActionUpload:
		err := u.startActionUpload()
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

func (u User) startActionUpload() error {
	if u.req.info.Size <= 0 {
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
	u.req.setFileSize(uint64(size))
	return nil
}

func (u User) createChannelIfNotExists() error {
	channel, err := u.req.channel.File()
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

func (u *User) processChunk(chunk []byte) error {
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
	return u.count+int64(len(chunk)) > int64(u.req.info.Size)
}

func (u User) stream(size uint, f func(buf []byte)) error {
	err := files.Stream(u.file, size, f)
	if err != nil {
		return err
	}
	return nil
}

type req struct {
	info    fs.FileInfo
	channel Channel
}

func (r *req) set(payload StartPayload) {
	r.info = payload.FileInfo
	r.channel = payload.Channel
}

func (r req) setFileSize(size uint64) {
	r.info.Size = size
}

func (r req) file() (fs.File, error) {
	f, err := fs.NewFileFromString(r.channel.Name) // {channel}
	if err != nil {
		log.Println(err)
		return fs.File{}, errors.New("invalid channel name: " + r.channel.Name)
	}
	err = f.Append(r.info.Value) // {channel}/{file.txt}
	if err != nil {
		log.Println(err)
		return fs.File{}, errors.New("invalid file: " + r.info.Value)
	}
	return f, nil
}
