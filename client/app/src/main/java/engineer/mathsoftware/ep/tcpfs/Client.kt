package engineer.mathsoftware.ep.tcpfs

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

const val HOST: String = "localhost"
const val PORT: Int = 8080
const val URI: String = "http://$HOST:$PORT"

class Client {

}
