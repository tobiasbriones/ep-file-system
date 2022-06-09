// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.util.Base64
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
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

    fun writeCommandListChannels() {
        val os = socket.getOutputStream()
        val msg = JSONObject()
        val cmd = JSONObject()

        cmd.put("REQ", "LIST_CHANNELS")
        msg.put("Command", cmd)

        os.write(
            msg.toString()
                .toByteArray()
        )
    }

    fun writeCommandListFiles(channel: String) {
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
    }

    fun writeCommandCID() {
        val os = socket.getOutputStream()
        val msg = JSONObject()
        val cmd = JSONObject()

        cmd.put("REQ", "CID")
        msg.put("Command", cmd)

        os.write(
            msg.toString()
                .toByteArray()
        )
    }

    suspend fun stream(bytes: ByteArray, l: (progress: Float) -> Unit): Int {
        val size = bytes.size
        val os = socket.getOutputStream()
        var count = 0
        var chunksTotal = 0

        while (count < size) {
            var end = count + SERVER_BUF_SIZE - 1
            end = if (end >= size) size - 1 else end
            val chunk = bytes.sliceArray(
                IntRange(count, end)
            )
            os.write(chunk)
            count += SERVER_BUF_SIZE
            chunksTotal++
            withContext(Dispatchers.Main) {
                l(getPercentage(count, size))
            }
        }
        return chunksTotal
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

    suspend fun downstream(size: Int, l: (progress: Float) -> Unit): ByteArray {
        var array = ByteArray(0)
        var count = 0
        while (count < size) {
            val chunk = readChunk()
            array += chunk
            count += chunk.size
            withContext(Dispatchers.Main) {
                l(getPercentage(count, size))
            }
        }
        return array
    }

    fun readChunk(): ByteArray {
        val chunk = ByteArray(SERVER_BUF_SIZE)
        val n = socket.getInputStream()
            .read(chunk)
        return chunk.sliceArray(IntRange(0, n - 1))
    }

    fun readNext(): ByteArray {
        // TODO give a big enough buffer to avoid while loop for now and
        //  message end char
        val buff = ByteArray(SERVER_BUF_SIZE*4)
        socket.getInputStream()
            .read(buff)
        return buff
    }
}

private fun getPercentage(count: Int, size: Int): Float {
    return if (count >= size) 1.0f
    else count.toFloat() / size.toFloat()
}
