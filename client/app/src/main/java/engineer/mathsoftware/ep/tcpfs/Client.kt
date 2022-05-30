package engineer.mathsoftware.ep.tcpfs

import android.util.Log
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
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
            }catch (e: ConnectException) {
                Log.d("UPLOAD", e.message.toString())
            }
        }
    }

    fun disconnect() {
        socket.close()
    }
}
