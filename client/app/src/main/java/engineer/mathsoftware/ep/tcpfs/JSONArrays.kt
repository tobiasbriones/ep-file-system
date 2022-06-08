// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import org.json.JSONArray

fun JSONArray.toStringList(): List<String> {
    val values = Array(length()) {
        getString(it)
    }
    return values.toList()
}
