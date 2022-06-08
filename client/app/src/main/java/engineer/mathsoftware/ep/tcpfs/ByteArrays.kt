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

