// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

import { Client } from './tcp.mjs';

const client = Client(handle);

client.connect();

function handle(msg) {
  console.log(msg);
}
