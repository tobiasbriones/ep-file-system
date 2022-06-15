// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package main

import (
	"fs/process"
	"io"
	"log"
	"net"
)

type state struct {
	conn    net.Conn
	process process.Process
	channel process.Channel
	quit    func()
	change  chan struct{}
}

func newState(
	conn net.Conn,
	osFsRoot string,
	quit func(),
	change chan struct{},
) state {
	return state{
		conn:    conn,
		process: process.NewProcess(osFsRoot),
		channel: process.Channel{},
		quit:    quit,
		change:  change,
	}
}

func (s *state) next() {
	if s.isOnHold() {
		return
	}
	switch s.process.State() {
	case process.Data:
		s.listenData()
	case process.Stream:
		s.listenStream()
	case process.Eof:
		s.handleEof()
	case process.Error:
		s.quit()
	}
}

// Returns true iff the process is not on hold. That is, iff isOnHold is false.
func (s state) isInProgress() bool {
	return !s.isOnHold()
}

// Returns true iff the process is not involved into any action in progress.
func (s state) isOnHold() bool {
	state := s.process.State()
	return state == process.Start || state == process.Done || state == process.Error
}

func (s *state) start(msg Message) {
	payload, err := msg.StartPayload()
	if err != nil {
		s.error("fail to read StartPayload")
		return
	}
	err = s.process.Start(payload)
	if err != nil {
		s.error(err.Error())
		return
	}
	// TODO check breaks backward compatibility
	//if s.process.User().Channel().Name != s.channel.Name {
	//	s.error("Client channel doesn't match")
	//	return
	//}
	log.Println("Accepting request:", payload)
	s.onProcessStarted()
}

func (s *state) onProcessStarted() {
	switch s.process.Action() {
	case process.ActionUpload:
		s.onActionUploadStarted()
	case process.ActionDownload:
		s.onActionDownloadStarted()
	}
}

func (s *state) onActionUploadStarted() {
	err := writeState(process.Data, s.conn)
	if err != nil {
		s.error("Fail to write state=DATA")
		return
	}
	log.Println("State DATA sent")
}

func (s *state) onActionDownloadStarted() {
	s.writeStreamState(process.StreamPayload{FileInfo: s.process.User().FileInfo()})
}

func (s *state) listenData() {
	chunk, err := readChunk(s.conn)
	if err != nil {
		s.handleReadError(err, "fail to read chunk")
		return
	}
	err = s.process.Data(chunk)
	if err != nil {
		s.error(err.Error())
		return
	}
}

func (s *state) onChunkProcessed() {
	if s.process.State() == process.Eof {
		err := writeState(process.Eof, s.conn)
		if err != nil {
			s.error("Fail to write state=EOF")
			return
		}
	}
}

func (s *state) handleEof() {
	err := s.writeEofState()
	if err != nil {
		s.error("fail to write EOF state")
		return
	}
	log.Println("State EOF sent, waiting for EOF message")
	msg, err := readMessage(s.conn, readTimeOut)
	if err != nil {
		s.handleReadError(err, "fail to read EOF message")
		return
	}
	s.eof(msg)
}

func (s *state) writeEofState() error {
	err := writeState(process.Eof, s.conn)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) eof(msg Message) {
	if msg.State != process.Eof {
		s.error("expecting EOF")
		return
	}
	log.Println("DONE!")
	err := s.process.Done()
	if err != nil {
		s.error("fail to write state=DONE on server")
		return
	}
	err = writeState(process.Done, s.conn)
	if err != nil {
		s.error("fail to write state=DONE")
		return
	}

	// If a file was uploaded, notify
	if s.process.Action() == process.ActionUpload {
		log.Println("File was uploaded, sending notification")
		s.change <- struct{}{}
	}
}

func (s *state) writeStreamState(payload process.StreamPayload) {
	p, err := NewPayloadFrom(payload)
	if err != nil {
		s.error("Fail to read payload from StreamPayload")
		return
	}
	msg := Message{
		State:   process.Stream,
		Payload: p,
	}
	err = writeMessage(msg, s.conn)
	if err != nil {
		s.error("Fail to write state=STREAM")
		return
	}
	log.Println("Payload sent, writing state=STREAM", payload)
}

func (s *state) listenStream() {
	log.Println("Listening for client STREAM signal")
	msg, err := readMessage(s.conn, readTimeOut)
	if err != nil {
		s.handleReadError(err, "fail to read status STREAM")
		return
	}
	if msg.State != process.Stream {
		s.error("wrong client state, state=STREAM was expected")
		return
	}
	s.stream()
}

func (s *state) stream() {
	err := s.process.Stream(
		bufSize,
		func(buf []byte) {
			_, err := s.conn.Write(buf)
			if err != nil {
				// TODO Fix StreamLocalFile paradigm
				s.error("Fail to write chunk")
				return
			}
		},
	)
	if err != nil {
		s.error("fail to stream file: " + err.Error())
		return
	}
	log.Println("File sent to client, waiting for client state EOF")
	msg, err := readMessage(s.conn, readTimeOut)
	if err != nil {
		s.error("Server error, fail to read state=EOF")
		return
	}
	if msg.State != process.Eof {
		s.error("Fail to read state=EOF")
		return
	}

	log.Println("Sending state DONE")
	err = writeState(process.Done, s.conn)
	if err != nil {
		s.error("Fail to write state=DONE")
		return
	}
}

func (s *state) handleReadError(err error, msg string) {
	if err == io.EOF {
		log.Println("Communication closed by the client")
		s.quit()
		return
	}
	log.Println(msg, err)
	s.error(msg)
}

func (s *state) error(msg string) {
	log.Println("ERROR:", msg)
	s.process.Error()
	writeErrorState(msg, s.conn)
}
