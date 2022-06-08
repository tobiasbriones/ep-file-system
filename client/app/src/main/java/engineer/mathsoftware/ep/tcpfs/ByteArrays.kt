// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import engineer.mathsoftware.ep.tcpfs.DataType.*
import org.json.JSONArray
import org.json.JSONObject

sealed interface DataType {
    data class Message(val value: JSONObject) : DataType
    data class Array(val value: JSONArray) : DataType
    data class Raw(val value: ByteArray) : DataType
}

// This is a low-level response parsing, a high-level construct like states
// must be used to understand the actual meaning of the server.
//
// For example: It can be a response that is a data chunk for a small JSON
// text file, but this can be either interpreted as a server message since
// server messages are serialized as JSON objects or as a data chunk, it depends
// on the machine state.
fun ByteArray.parse(): DataType { // TODO return Result with error instead
    val str = String(this)
    if (str.isEmpty()) return Raw(this)
    return when (str[0]) {
        '{'  -> Message(parseMessage())
        '['  -> Array(parseArray())
        else -> Raw(this)
    }
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
