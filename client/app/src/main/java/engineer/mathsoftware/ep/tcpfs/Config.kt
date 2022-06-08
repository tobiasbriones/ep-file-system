// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.app.Activity
import android.content.Context

class Config(activity: Activity) {
    private val res = activity.resources
    private val sharedPref = activity.getPreferences(Context.MODE_PRIVATE)

    fun getServerHost(): String? {
        val def = res.getString(R.string.saved_host_default_value)
        return sharedPref.getString(
            res.getString(R.string.saved_host_key),
            def
        )
    }

    fun saveServerHost(host: String) {
        with(sharedPref.edit()) {
            putString(res.getString(R.string.saved_host_key), host)
            apply()
        }
    }
}
