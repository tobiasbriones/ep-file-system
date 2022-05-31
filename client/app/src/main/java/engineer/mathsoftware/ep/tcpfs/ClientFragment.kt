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
import androidx.recyclerview.widget.LinearLayoutManager
import engineer.mathsoftware.ep.tcpfs.databinding.FragmentClientBinding
import kotlinx.coroutines.launch
import java.net.SocketException

class ClientFragment : Fragment() {
    private val files = ArrayList<String>()
    private var _binding: FragmentClientBinding? = null
    private lateinit var filesAdapter: FilesAdapter

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!
    private lateinit var client: Client

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        _binding = FragmentClientBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        filesAdapter = FilesAdapter(files) { download(it) }
        initFileList()
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

    private fun initFileList() {
        val recyclerView = binding.fileList
        val layoutManager = LinearLayoutManager(requireContext())
        recyclerView.layoutManager = layoutManager
        recyclerView.adapter = filesAdapter
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
                handleConnectionOpened()
            }
        }
    }

    private fun handleConnectionOpened() {
        val channel = arguments?.getString("channel")
        if (channel != null) {
            client.channel = channel
        }
        readCID()
        readFiles()
    }

    private fun readCID() {
        lifecycleScope.launch {
            val cid = client.readCID()
            binding.clientText.text = "Client #$cid"
            binding.channelText.text = "Channel: ${client.channel}"
            binding.infoText.text = "Connected"
        }
    }

    private fun readFiles() {
        lifecycleScope.launch {
            val res = client.readFiles()
            files.clear()
            files.addAll(res)
            filesAdapter.notifyDataSetChanged()
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
            try {
                client.file = file
                client.upload(bytes) {
                    val percentage = it * 100
                    binding.infoText.text = "Uploading $percentage%"
                }
                handleFileUploaded()
            }
            catch (e: SocketException) {
                println("ERROR: ${e.message}")
                handleConnectionFailed()
            }
        }
    }

    private fun handleFileUploaded() {
        binding.infoText.text = "File uploaded: ${client.file}"
    }

    private fun download(file: String) {
        lifecycleScope.launch {
            try {
                client.file = file
                // TODO client.download()
            }
            catch (e: SocketException) {
                println("ERROR: ${e.message}")
                handleConnectionFailed()
            }
        }
    }
}
