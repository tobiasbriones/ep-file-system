// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

export function newConn(url, handleConnected, handleMessage) {
  const socket = new WebSocket(url);
  socket.onopen = handleConnected;
  socket.onmessage = handleMessage;
  return {
    sendCommandConnectedUsers() {
      const msg = {
        Command: {
          REQ: 'CONNECTED_USERS'
        }
      };
      socket.send(JSON.stringify(msg));
    }
  };
}
