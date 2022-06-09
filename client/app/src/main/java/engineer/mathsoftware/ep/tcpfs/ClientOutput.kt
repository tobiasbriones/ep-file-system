// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.widget.TextView

class ClientOutput(private val info: TextView) : Output {
    override fun updateUploadProgress(progress: Float) {
        val percentage = progress * 100
        info.text = "Uploading $percentage%"
    }

    override fun uploadDone(file: String, chunksTotal: Int) {
        info.text = """
            File uploaded: ${file} | $chunksTotal chunks sent
        """.trimIndent()
    }
}
