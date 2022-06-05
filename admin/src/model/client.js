// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

import { newConn } from '@/model/conn';

const HOST = 'localhost';
const PORT = 8080;
const URL = `ws://${ HOST }:${ PORT }`;

export function newClient(handleConnected) {
  function handleMessage(event) {
    console.log(event.data);
  }

  const conn = newConn(URL, handleConnected, handleMessage);
  return {
    readUsers() {
      conn.sendCommandConnectedUsers();
    }
  };
}