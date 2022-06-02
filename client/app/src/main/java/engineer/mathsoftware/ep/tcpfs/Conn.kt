// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.util.Base64
import org.json.JSONArray
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.Socket
import java.nio.charset.StandardCharsets

class Conn(private val socket: Socket) {
    private val reader = BufferedReader(
        InputStreamReader(
            socket.getInputStream()
        )
    )

    fun readChannels(): List<String> {
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
        return channels.toList()
    }

    fun readFiles(channel: String): List<String> {
        val os = socket.getOutputStream()
        val msg = JSONObject()
        val cmd = JSONObject()

        cmd.put("REQ", "LIST_FILES")
        cmd.put("CHANNEL", channel)
        msg.put("Command", cmd)

        os.write(
            msg.toString()
                .toByteArray()
        )
        val res = reader.readLine()
        if (res == null || res == "null") {
            return ArrayList()
        }
        val jsonArray = JSONArray(res)
        val channels = Array(jsonArray.length()) {
            jsonArray.getString(it)
        }
        return channels.toList()
    }

    fun readCID(): Int {
        val os = socket.getOutputStream()
        val msg = JSONObject()
        val cmd = JSONObject()

        cmd.put("REQ", "CID")
        msg.put("Command", cmd)

        os.write(
            msg.toString()
                .toByteArray()
        )
        val res = reader.readLine()
        return Integer.parseInt(res)
    }

    fun stream(bytes: ByteArray, l: (progress: Float) -> Unit) {
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
            l(getPercentage(count, size))
        }
    }

    fun readState(): String {
        val ser = readMessage()
        return ser.get("State")
            .toString()
    }

    fun readMessage(): JSONObject {
        val res = reader.readLine()
        return JSONObject(res)
    }

    fun writeMessage(msg: JSONObject) {
        val os = socket.getOutputStream()
        os.write(
            msg.toString()
                .toByteArray()
        )
    }

    fun readData(msg: JSONObject): JSONObject {
        val data = msg["Data"].toString()
        val str = Base64.decode(data, Base64.DEFAULT)
            .toString(StandardCharsets.UTF_8)
        return JSONObject(str)
    }

    fun downstream(size: Int, l: (progress: Float) -> Unit): ByteArray {
        var array = ByteArray(0)
        var count = 0
        while (count < size) {
            val chunk = readChunk()
            array += chunk
            count += chunk.size
            l(getPercentage(count, size))
        }
        return array
    }

    fun readChunk(): ByteArray {
        val chunk = ByteArray(SERVER_BUF_SIZE)
        val n = socket.getInputStream()
            .read(chunk)
        return chunk.sliceArray(IntRange(0, n - 1))
    }
}

private fun getPercentage(count: Int, size: Int): Float {
    return if (count >= size) 1.0f
    else count.toFloat() / size.toFloat()
}
