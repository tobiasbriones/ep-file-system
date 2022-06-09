// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.net.Uri
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONException
import org.json.JSONObject
import java.net.*

enum class Action {
    UPLOAD,
    DOWNLOAD
}

// TODO temp. response values, server is still in beta
const val UPDATE = 2
const val OK = 3

const val PORT: Int = 8080

data class Input(
    val onChannelList: ((channels: List<String>) -> Unit)? = null,
    val onFileList: ((channels: List<String>) -> Unit)? = null,
    val onCID: ((cid: Int) -> Unit)? = null,
    val onUpdate: (() -> Unit)? = null,
)

class Client(
    private val socket: Socket,
    private val conn: Conn,
    private val input: Input,
    output: Output
) {
    companion object {
        suspend fun newInstance(
            host: String,
            input: Input,
            output: Output
        ): Client? {
            return withContext(Dispatchers.IO) {
                try {
                    val address = InetAddress.getByName(host)
                    val socket = Socket(address, PORT)
                    Client(socket, Conn(socket), input, output)
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

    private val state = State(conn, output)
    var file: String = ""
    var channel: String = "test"

    suspend fun disconnect() {
        println("Disconnecting...")
        withContext(Dispatchers.IO) {
            socket.close()
        }
    }

    suspend fun listen() {
        withContext(Dispatchers.IO) {
            while (socket.isConnected) {
                try {
                    val data = conn.readNext()
                    onData(data)
                }
                catch (e: SocketException) {
                }
                catch (e: Exception) {
                    println("ERROR: $e")
                }
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
        if (response == UPDATE) {
            input.onUpdate?.invoke()
            return
        }
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

    suspend fun upload(bytes: ByteArray) {
        withContext(Dispatchers.IO) {
            val payload = StartPayload(file, channel, bytes.size)
            state.startUpload(bytes, payload)
        }
    }

    suspend fun download(uri: Uri) {
        withContext(Dispatchers.IO) {
            val payload = StartPayload(file, channel)
            state.startDownload(payload, uri)
        }
    }
}
