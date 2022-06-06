<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# Development

In this article, I will document the personal experience I had with Go and this
project by also employing my previous thoughts and experiences on some design
decisions.

## Domain Model

I had some loose file system logic, some weird data structures, and the utility
functions. This is a design smell, so I decided to build a basic domain model
based on my
other [file system](https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system/tree/main/model).

This basic model just adjusts to the project requirements, so I don't have to
implement other features like `getParent`, or tree traversal, etc.

The Android and Web clients need to read these models too, building the base
model for the file system (`File`, `Directory`, `Node`, operations like file
name/extension, etc.) would be insanely expensive.

To address this problem, the domain model in the server has to be adequate for
this project since it's the most important module. The clients can read just
primitive data types to consume the content.

This way, the loose logic is correctly coupled, and testable.

I always test the domain modules as much as I can to avoid propagating errors
forward to presentation layers, and avoid system downtimes or failures.

## Navigation

- [Enums in Go](enums-in-go)
- [EOF Issues](eof-issues)
- [Go Dependency System](go-dependency-system)
- [A Development Cycle Can Be Shorter](a-development-cycle-can-be-shorter)
- [Refactorize Before It's Too Late](refactorize-before-it's-too-late)
- [Measure Cohesiveness](measure-cohesiveness)
- [Plain TCP vs Web Sockets](plain-tcp-vs-web-sockets)
