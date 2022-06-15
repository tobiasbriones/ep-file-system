<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# WS Client

## Introduction

Raw TCP servers and clients are low-level implementations of sockets as said
along this project documentation. They send streams instead of
standalone packages and have a blocking nature while Web Sockets are
implemented on top of TCP but are higher-leve, this including sending
asynchronous packages of messages.

The main problem with the underlying issue is that while TCP being low-level
it needs to be implemented close to the OS layer with traditional sockets
like Java/Go/C/C++/Rust/Android TCP sockets.

Raw TCP sockets can't be added to web browsers due to security
vulnerabilities so Web Sockets are the ones that can (barely) run in the web
browser.

## Possibilities

Web clients are required for any application, at least they should be the 
default as they run on virtually any device. But does that mean that we need 
to use only Web Sockets?

One can use Web Sockets by default, but raw TCP Sockets are better for 
high-performance applications like file systems.

That is due to the extra overhead added by Web Sockets to upgrade the 
connection from HTTP to WS as well as other \*minor*\ inconveniences.
