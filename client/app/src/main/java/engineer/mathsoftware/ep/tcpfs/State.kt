// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import engineer.mathsoftware.ep.tcpfs.Process.*
import org.json.JSONArray
import org.json.JSONObject

enum class Process {
    START,
    DATA,
    STREAM,
    EOF,
    DONE,
    ERROR
}

data class StartPayload(
    val file: String,
    val channel: String,
    val size: Int = 0
)

interface Output {
    fun updateUploadProgress(progress: Float)
    fun uploadDone(file: String, chunksTotal: Int)
}

class State(private val conn: Conn, private val output: Output) {
    private var state = START
    private var action = Action.UPLOAD
    private var file = ""
    private var data = ByteArray(0)
    private var chunksTotal = 0

    fun isInProgress() = !isOnHold()

    fun isOnHold() = when (state) {
        START, ERROR -> true
        else         -> false
    }

    suspend fun parse(data: ByteArray) {
        when (state) {
            DATA -> readStateData(data.parseMessage())
            EOF  -> readStateEOF(data.parseMessage())
            DONE -> readStateDone(data.parseMessage())
        }
    }

    fun startUpload(bytes: ByteArray, p: StartPayload) {
        if (isInProgress()) {
            return
        }
        var msg = getStartMessage(Action.UPLOAD, p)
        conn.writeMessage(msg)
        action = Action.UPLOAD
        file = p.file
        data = bytes
        chunksTotal = 0
        state = DATA
        println("Start message sent")
    }

    private suspend fun readStateData(msg: JSONObject) {
        if (msg["State"] != "DATA") {
            state = ERROR
            print("ERROR: Fail to read state DATA: $msg")
            return
        }
        println("STATE=DATA confirmed")
        sendData()
    }

    private suspend fun sendData() {
        chunksTotal = conn.stream(data, output::updateUploadProgress)
        state = EOF
    }

    private fun readStateEOF(msg: JSONObject) {
        if (msg["State"] != "EOF") {
            state = ERROR
            print("ERROR: Fail to read state EOF: $msg")
            return
        }
        println("STATE=EOF confirmed")
        sendEof()
    }

    private fun sendEof() {
        val msg = getEofMessage()
        conn.writeMessage(msg)
        state = DONE
        println("EOF message sent")
    }

    private fun readStateDone(msg: JSONObject) {
        if (msg["State"] != "DONE") {
            state = ERROR
            print("ERROR: Fail to read state EOF: $msg")
            return
        }
        done()
    }

    private fun done() {
        state = START
        output.uploadDone(file, chunksTotal)
        println("Done!")
    }
}

private fun getStartMessage(action: Action, p: StartPayload): JSONObject {
    val payload = getStartPayload(action, p)
    val ser = JSONObject()
    ser.put("State", START)
    ser.put("Data", payload)
    return ser
}

private fun getStartPayload(action: Action, p: StartPayload): JSONArray {
    val payload = """
            {
                "Action": ${action.ordinal},
                "Value": "${p.file}",
                "Size": ${p.size},
                "Channel": {
                    "Name": "${p.channel}"
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
    ser.put("State", EOF)
    ser.put("Data", null)
    return ser
}

private fun getStreamMessage(): JSONObject {
    val ser = JSONObject()
    ser.put("State", STREAM)
    ser.put("Data", null)
    return ser
}
