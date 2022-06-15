// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

import WebSocket, { WebSocketServer } from 'ws';

const PORT = 8081;

export function WSS() {
  const wss = new WebSocketServer({ port: PORT });
  const isClientReady = c => c.readyState === WebSocket.OPEN;
  const sendMessage = (c, msg) => c.send(msg);
  return {
    init() {
      wss.on('connection', ws => {
        console.log('WSS connection started');
        ws.on('message', msg => {
          console.log(msg);
        });
      });
      wss.on('error', (error) => {
        console.log(error);
      });
    },
    broadcast(msg) {
      console.log(msg);
      wss.clients
         .forEach(c => {
           if (isClientReady(c)) {
             sendMessage(c, JSON.stringify(msg));
           }
         });
    }
  };
}

