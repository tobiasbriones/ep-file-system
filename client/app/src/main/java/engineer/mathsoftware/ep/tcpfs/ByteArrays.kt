// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import engineer.mathsoftware.ep.tcpfs.DataType.*
import org.json.JSONArray
import org.json.JSONObject

enum class DataType {
    MESSAGE,
    ARRAY,
    RAW
}

// It parses this data as a JSONObject assuming this is a generic server
// response message.
fun ByteArray.parseMessage(): JSONObject {
    return JSONObject(String(this))
}


// It parses this data as a JSONArray assuming this is a generic server
// response array.
fun ByteArray.parseArray(): JSONArray {
    return JSONArray(String(this))
}

// This is a low-level response parsing, a high-level construct like states
// must be used to understand the actual meaning of the server.
//
// For example: It can be a response that is a data chunk for a small JSON
// text file, but this can be interpreted as a server message or as a data
// chunk, it depends on the machine state.
fun ByteArray.type(): DataType {
    val str = String(this)
    if (str.isEmpty()) return RAW
    return when (str[0]) {
        '{'  -> MESSAGE
        '['  -> ARRAY
        else -> RAW
    }
}
