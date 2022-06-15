// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

const HOST = 'localhost';
const PORT = 8081;
const URL = `ws://${ HOST }:${ PORT }`;

export function Client(handleUsers) {
  const socket = new WebSocket(URL);
  const handleData = data => readData(data).apply(handleUsers);
  const handleClose = e => console.log(`Connection closed, code=${ e.code } reason=${ e.reason }`);
  const handleError = e => console.log(`ERROR: ${ e }`);
  return {
    init() {
      socket.onmessage = e => handleData(e.data);
      socket.onclose = e => (
        e.wasClean ? handleClose(e) : handleError('Connection died')
      );
      socket.onerror = handleError;
    }
  };
}

function readData(data) {
  const array = JSON.parse(data);
  return {
    apply(f) {
      f(array);
    }
  };
}
