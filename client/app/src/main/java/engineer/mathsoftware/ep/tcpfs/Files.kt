package engineer.mathsoftware.ep.tcpfs

import android.content.ContentResolver
import android.net.Uri
import java.io.ByteArrayOutputStream
import java.io.InputStream

const val SERVER_BUF_SIZE = 1024

fun read(res: ContentResolver, uri: Uri): ByteArray {
    val stream: InputStream? = res.openInputStream(uri)
    if (stream != null) {
        return getBytes(stream)
    }
    return ByteArray(0)
}

private fun getBytes(inputStream: InputStream): ByteArray {
    val byteBuffer = ByteArrayOutputStream()
    val buffer = ByteArray(SERVER_BUF_SIZE)
    var len = 0
    while (inputStream.read(buffer)
            .also { len = it } != -1) {
        byteBuffer.write(buffer, 0, len)
    }
    return byteBuffer.toByteArray()
}
