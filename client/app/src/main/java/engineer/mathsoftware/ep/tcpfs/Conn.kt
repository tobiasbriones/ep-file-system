// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.Socket

class Conn(private val socket: Socket) {
    private val reader = BufferedReader(
        InputStreamReader(
            socket.getInputStream()
        )
    )

    fun readChannels(): Array<String> {
        val os = socket.getOutputStream()
        val msg = JSONObject()
        val cmd = JSONObject()

        cmd.put("REQ", "LIST_CHANNELS")
        msg.put("Command", cmd)

        os.write(
            msg.toString()
                .toByteArray()
        )
        val res = reader.readLine()
        val jsonArray = JSONArray(res)
        val channels = Array(jsonArray.length()) {
            jsonArray.getString(it)
        }
        return channels
    }

    fun stream(bytes: ByteArray) {
        val size = bytes.size
        val os = socket.getOutputStream()
        var count = 0

        while (count < size) {
            var end = count + SERVER_BUF_SIZE - 1
            end = if (end >= size) size - 1 else end
            val chunk = bytes.sliceArray(
                IntRange(count, end)
            )
            os.write(chunk)
            count += SERVER_BUF_SIZE
        }
        println("Finished sending chunks: $count")
    }

    fun readState(): String {
        val res = reader.readLine()
        val ser = JSONObject(res)
        return ser.get("State")
            .toString()
    }

    fun writeMessage(msg: JSONObject) {
        val os = socket.getOutputStream()
        os.write(
            msg.toString()
                .toByteArray()
        )
    }
}
