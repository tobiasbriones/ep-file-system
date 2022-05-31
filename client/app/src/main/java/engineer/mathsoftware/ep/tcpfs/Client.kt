package engineer.mathsoftware.ep.tcpfs

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONArray
import org.json.JSONException
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader
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

class Client(private val socket: Socket) {
    companion object {
        suspend fun newInstance(): Client? {
            return withContext(Dispatchers.IO) {
                try {
                    val address = InetAddress.getByName(HOST)
                    val socket = Socket(address, PORT)
                    Client(socket)
                }
                catch (e: ConnectException) {
                    println("ERROR: " + e.message.toString())
                    null
                }
            }
        }
    }

    var file: String = ""
    private var channel: String = "test"

    suspend fun disconnect() {
        withContext(Dispatchers.IO) {
            socket.close()
        }
    }

    suspend fun upload(bytes: ByteArray) {
        withContext(Dispatchers.IO) {
            try {
                val reader = BufferedReader(
                    InputStreamReader(
                        socket.getInputStream()
                    )
                )
                var msg = getStartMessage(Action.UPLOAD, bytes.size)

                // writer.println(msg)
                socket.getOutputStream()
                    .write(
                        msg.toString()
                            .toByteArray()
                    )
                println("Start message sent")
                var res = reader.readLine()

                println("response $res")
                var ser = JSONObject(res)
                var state = ser.get("State")
                println("STATE: $state")

                // Upload
                val size = bytes.size
                var count = 0

                while (count < size) {
                    var end = count + SERVER_BUF_SIZE - 1
                    end = if (end >= size) size - 1 else end
                    val chunk = bytes.sliceArray(
                        IntRange(count, end)
                    )
                    socket.getOutputStream()
                        .write(chunk)
                    count += SERVER_BUF_SIZE
                }
                println("Finished sending chunks: $count")

                res = reader.readLine()
                ser = JSONObject(res)
                state = ser.get("State")
                println("Received status: $state")


                msg = getEofMessage()
                println(msg)
                // writer.println(msg)
                socket.getOutputStream()
                    .write(
                        msg.toString()
                            .toByteArray()
                    )

                res = reader.readLine()
                ser = JSONObject(res)
                state = ser.get("State")
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
}
