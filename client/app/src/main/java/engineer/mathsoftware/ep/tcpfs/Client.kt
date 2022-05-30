package engineer.mathsoftware.ep.tcpfs

import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONObject
import java.io.PrintWriter
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
const val URI: String = "http://$HOST:$PORT"

class Client {
    var file: String = ""
    private var channel: String = "test"
    private lateinit var socket: Socket

    suspend fun connect() {
        withContext(Dispatchers.IO) {
            val address = InetAddress.getByName(HOST)
            try {
                socket = Socket(address, PORT)
                Log.d("UPLOAD", "Connection established")
            }
            catch (e: ConnectException) {
                Log.d("UPLOAD", e.message.toString())
            }
        }
    }

    fun disconnect() {
        socket.close()
    }

    suspend fun upload(bytes: ByteArray) {
        val msg = getStartMessage(Action.UPLOAD, bytes.size)
        withContext(Dispatchers.IO) {
            PrintWriter(socket.getOutputStream()).use {
                it.print(msg)
                Log.d("UPLOAD", msg.toString())
            }
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
                "Value": $file,
                "Size": $size,
                "Channel": {
                    "Name": $channel
                }
            }
        """.trimIndent()
        var json = JSONObject(payload)
        val arr = JSONArray()
        json.toString().toByteArray().forEach {
            arr.put(it)
        }
        return arr
    }
}
