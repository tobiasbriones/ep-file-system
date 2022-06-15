// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

import Net from 'net';

const PORT = 8080;
const HOST = 'localhost';
const OK = 3;

run();

function run() {
  const client = new Net.Socket();

  client.connect({ port: PORT, host: HOST }, () => handleConnect(client));
  client.on('data', handleData);

  client.on('end', () => {
    console.log('Ending the connection');
  });
}

function handleConnect(client) {
  console.log('TCP connection established with the FS server.');
  const msg = {
    Command: {
      REQ: 'SUBSCRIBE_TO_LIST_CONNECTED_USERS'
    }
  };
  client.write(JSON.stringify(msg));
}

function handleData(data) {
  const str = data.toString();
  const messages = [];
  let start = 0;
  for (let i = 0; i < str.length; i++) {
    const char = str[i];
    if (char === '\n') {
      const msg = JSON.parse(str.substring(start, i));
      messages.push(msg);
      start = i;
    }
  }
  messages.forEach(readMessage);
}

function readMessage(msg) {
  if (msg.Response === OK) {
    readCommandResponse(msg.Command);
  }
  else {
    console.log('Response not OK');
  }
}

function readCommandResponse(cmd) {
  if (cmd['REQ'] === 'SUBSCRIBE_TO_LIST_CONNECTED_USERS') {
    readPayload(cmd['PAYLOAD']);
  }
  else {
    console.log('Command request not expected');
  }
}

function readPayload(payload) {
  const users = JSON.parse(payload);
  console.log(users);
}
