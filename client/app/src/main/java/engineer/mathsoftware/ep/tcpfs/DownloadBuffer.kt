// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class DownloadBuffer(
    private val size: Int,
    private val l: (progress: Float) -> Unit
) {
    var data = ByteArray(0)
    var count = 0
    var chunksTotal = 0

    fun isDone() = count == size

    fun isOverflowed() = count > size

    suspend fun append(chunk: ByteArray) {
        data += chunk
        count += chunk.size
        chunksTotal++
        withContext(Dispatchers.Main) {
            l(getPercentage(count, size))
        }
    }
}
