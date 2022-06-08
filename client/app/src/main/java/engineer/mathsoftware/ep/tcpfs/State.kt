// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import engineer.mathsoftware.ep.tcpfs.Process.*

enum class Process {
    START,
    DATA,
    STREAM,
    EOF,
    DONE,
    ERROR
}

class State {
    private var state = START

    fun isInProgress() = !isOnHold()

    fun isOnHold() = when (state) {
        START, DONE, ERROR -> true
        else               -> false
    }
}
