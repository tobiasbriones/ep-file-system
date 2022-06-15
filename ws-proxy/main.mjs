// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

import { Client } from './tcp.mjs';
import { WSS } from './wss.mjs';

const client = Client(handle);
const wss = WSS();

client.connect();
wss.init();

function handle(msg) {
  wss.broadcast(msg);
}
