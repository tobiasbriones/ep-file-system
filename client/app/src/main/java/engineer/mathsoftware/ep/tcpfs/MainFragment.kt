// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.app.Activity
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.documentfile.provider.DocumentFile
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import com.google.android.material.snackbar.Snackbar
import engineer.mathsoftware.ep.tcpfs.databinding.FragmentMainBinding
import kotlinx.coroutines.launch
import org.json.JSONException

const val PICKFILE_REQUEST_CODE = 1

/**
 * A simple [Fragment] subclass as the default destination in the navigation.
 */
class MainFragment : Fragment() {
    private var _binding: FragmentMainBinding? = null
    private lateinit var client: Client

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {

        _binding = FragmentMainBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.buttonUpload.setOnClickListener {
            chooseFileToUpload()
        }
        connect()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
        disconnect()
    }

    override fun onActivityResult(
        requestCode: Int,
        resultCode: Int,
        data: Intent?
    ) {
        super.onActivityResult(requestCode, resultCode, data)
        if (resultCode != Activity.RESULT_OK) {
            return
        }
        if (data == null) {
            return
        }
        when (requestCode) {
            PICKFILE_REQUEST_CODE -> readFileToUpload(data.data)
        }
    }

    private fun connect() {
        lifecycleScope.launch {
            val c = Client.newInstance()

            if (c == null) {
                handleConnectionFailed()
            }
            else {
                println("connected")
                client = c
            }

            try {
                val channels = client.readChannels()
                println(channels.joinToString(", "))
            }
            catch (e: JSONException) {
                println(e.message)
            }
        }
    }

    private fun disconnect() {
        lifecycleScope.launch {
            client.disconnect()
        }
    }

    private fun chooseFileToUpload() {
        val intent = Intent(Intent.ACTION_OPEN_DOCUMENT).apply {
            addCategory(Intent.CATEGORY_OPENABLE)
            type = "*/*"
        }
        startActivityForResult(intent, PICKFILE_REQUEST_CODE)
    }

    private fun readFileToUpload(data: Uri?) {
        var bytes = ByteArray(0)
        val file = data?.let {
            DocumentFile.fromSingleUri(requireContext(), it)
        }?.name.toString()
        if (data != null) {
            bytes = read(requireContext().contentResolver, data)
        }

        lifecycleScope.launch {
            client.file = file
            client.upload(bytes)
        }
    }

    private fun handleConnectionFailed() {
        Snackbar.make(
            requireView(),
            "Fail to connect",
            Snackbar.LENGTH_LONG
        )
            .setAction("Action", null)
            .show()
    }
}
