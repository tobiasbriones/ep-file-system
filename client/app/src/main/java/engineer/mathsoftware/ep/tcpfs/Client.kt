// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONException
import org.json.JSONObject
import java.net.ConnectException
import java.net.InetAddress
import java.net.Socket

enum class Action {
    UPLOAD,
    DOWNLOAD
}

enum class State {
    START,
    DATA,
    STREAM,
    EOF,
    DONE,
    ERROR
}

const val HOST: String = "10.0.2.2" // This localhost IP works on the emulator
const val PORT: Int = 8080

class Client(private val socket: Socket, private val conn: Conn) {
    companion object {
        suspend fun newInstance(): Client? {
            return withContext(Dispatchers.IO) {
                try {
                    val address = InetAddress.getByName(HOST)
                    val socket = Socket(address, PORT)
                    Client(socket, Conn(socket))
                }
                catch (e: ConnectException) {
                    println("ERROR: " + e.message.toString())
                    null
                }
            }
        }
    }

    var file: String = ""
    var channel: String = "test"

    suspend fun disconnect() {
        withContext(Dispatchers.IO) {
            socket.close()
        }
    }

    suspend fun readChannels(): List<String> {
        return withContext(Dispatchers.IO) {
            conn.readChannels()
        }
    }

    suspend fun readFiles(): List<String> {
        return withContext(Dispatchers.IO) {
            conn.readFiles(channel)
        }
    }

    suspend fun readCID(): Int {
        return withContext(Dispatchers.IO) {
            conn.readCID()
        }
    }

    suspend fun upload(bytes: ByteArray, l: (progress: Float) -> Unit) {
        withContext(Dispatchers.IO) {
            try {
                var msg = getStartMessage(Action.UPLOAD, bytes.size)
                conn.writeMessage(msg)
                println("Start message sent")

                var state = conn.readState()
                println("STATE: $state")

                conn.stream(bytes, l)
                state = conn.readState()
                println("Received status: $state")

                msg = getEofMessage()
                conn.writeMessage(msg)
                println("EOF message sent")

                state = conn.readState()
                println("Received status: $state")
            }
            catch (e: JSONException) {
                println("ERROR: fail to read server response: " + e.message)
            }
            catch (e: NoSuchElementException) {
                println("Connection closed: $e.message")
            }
        }
    }

    suspend fun download(l: (progress: Float) -> Unit): ByteArray {
        return withContext(Dispatchers.IO) {
            try {
                var msg = getStartMessage(Action.DOWNLOAD)
                conn.writeMessage(msg)
                println("Start message sent")

                var res = conn.readMessage()
                var payload = conn.readData(res)
                val size = payload["Size"].toString().toInt()
                println("Payload: $payload")

                msg = getStreamMessage()
                conn.writeMessage(msg)

                val array = conn.downstream(size, l)

                if (array.size != size) {
                    println("ERROR: Overflow")
                }

                msg = getEofMessage()
                conn.writeMessage(msg)

                val done = conn.readMessage()
                println("State: ${done["State"]}")

                array
            }
            catch (e: JSONException) {
                println("ERROR: fail to read server response: " + e.message)
                ByteArray(0)
            }
            catch (e: NoSuchElementException) {
                println("Connection closed: $e.message")
                ByteArray(0)
            }
        }
    }

    suspend fun createChannel(channel: String) {
        withContext(Dispatchers.IO) {
            val msg = JSONObject()
            val cmd = JSONObject()
            cmd.put("REQ", "CREATE_CHANNEL")
            cmd.put("CHANNEL", channel)
            msg.put("Command", cmd)
            conn.writeMessage(msg)
        }
    }

    private fun getStartMessage(action: Action, size: Int = 0): JSONObject {
        val payload = getStartPayload(action, size)
        val ser = JSONObject()
        ser.put("State", State.START)
        ser.put("Data", payload)
        return ser
    }

    private fun getStartPayload(action: Action, size: Int = 0): JSONArray {
        val payload = """
            {
                "Action": ${action.ordinal},
                "Value": "$file",
                "Size": $size,
                "Channel": {
                    "Name": $channel
                }
            }
        """.trimIndent()
        var json = JSONObject(payload)
        val arr = JSONArray()
        json.toString()
            .toByteArray()
            .forEach {
                arr.put(it)
            }
        return arr
    }

    private fun getEofMessage(): JSONObject {
        val ser = JSONObject()
        ser.put("State", State.EOF)
        ser.put("Data", null)
        return ser
    }

    private fun getStreamMessage(): JSONObject {
        val ser = JSONObject()
        ser.put("State", State.STREAM)
        ser.put("Data", null)
        return ser
    }
}
