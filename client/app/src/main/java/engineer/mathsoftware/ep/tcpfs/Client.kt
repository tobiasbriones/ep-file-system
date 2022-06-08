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
import java.net.UnknownHostException

enum class Action {
    UPLOAD,
    DOWNLOAD
}

// TODO temp. response values, server is still in beta
const val OK = 3

const val PORT: Int = 8080

data class Input(
    val onChannelList: ((channels: List<String>) -> Unit)? = null,
    val onFileList: ((channels: List<String>) -> Unit)? = null,
    val onCID: ((cid: Int) -> Unit)? = null,
)

class Client(
    private val socket: Socket,
    private val conn: Conn,
    private val input: Input
) {
    companion object {
        suspend fun newInstance(host: String, input: Input): Client? {
            return withContext(Dispatchers.IO) {
                try {
                    val address = InetAddress.getByName(host)
                    val socket = Socket(address, PORT)
                    Client(socket, Conn(socket), input)
                }
                catch (e: ConnectException) {
                    println("ERROR: " + e.message.toString())
                    null
                }
                catch (e: UnknownHostException) {
                    println("ERROR: " + e.message.toString())
                    null
                }
            }
        }
    }

    private val state = State()
    var file: String = ""
    var channel: String = "test"

    suspend fun disconnect() {
        withContext(Dispatchers.IO) {
            socket.close()
        }
    }

    suspend fun listen() {
        withContext(Dispatchers.IO) {
            try {
                while (socket.isConnected) {
                    val data = conn.readNext()
                    onData(data)
                }
            }
            catch (e: Exception) {
                println("ERROR: $e")
            }
        }
    }

    private suspend fun onData(data: ByteArray) {
        if (state.isInProgress()) {
            state.parse(data)
        }
        else {
            parseResponse(data)
        }
    }

    private suspend fun parseResponse(data: ByteArray) {
        val msg = data.parseMessage()
        onMessage(msg)
    }

    private suspend fun onMessage(msg: JSONObject) {
        val response = msg.getInt("Response")
println(msg)
        if (response != OK) {
            // TODO handle
            println("ERROR: Response not OK")
        }
        val command = msg.getJSONObject("Command")
        when (command["REQ"].toString()) {
            "LIST_CHANNELS" -> onListChannelsResponse(command)
            "LIST_FILES"    -> onListFilesResponse(command)
            "CID"           -> onCIDResponse(command)
        }
    }

    private suspend fun onListChannelsResponse(command: JSONObject) {
        val payload = command.getString("PAYLOAD")
        val channels = JSONArray(payload).toStringList()
        withContext(Dispatchers.Main) {
            input.onChannelList?.invoke(channels)
        }
    }

    private suspend fun onListFilesResponse(command: JSONObject) {
        val payload = command.getString("PAYLOAD")
        val files = JSONArray(payload).toStringList()
        withContext(Dispatchers.Main) {
            input.onFileList?.invoke(files)
        }
    }

    private suspend fun onCIDResponse(command: JSONObject) {
        val cid = command.getInt("PAYLOAD")
        withContext(Dispatchers.Main) {
            input.onCID?.invoke(cid)
        }
    }

    suspend fun readChannels() {
        withContext(Dispatchers.IO) {
            conn.writeCommandListChannels()
        }
    }

    suspend fun readFiles() {
        withContext(Dispatchers.IO) {
            conn.writeCommandListFiles(channel)
        }
    }

    suspend fun readCID() {
        return withContext(Dispatchers.IO) {
            conn.writeCommandCID()
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
                val size = payload["Size"].toString()
                    .toInt()
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
        ser.put("State", Process.START)
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
        ser.put("State", Process.EOF)
        ser.put("Data", null)
        return ser
    }

    private fun getStreamMessage(): JSONObject {
        val ser = JSONObject()
        ser.put("State", Process.STREAM)
        ser.put("Data", null)
        return ser
    }
}
