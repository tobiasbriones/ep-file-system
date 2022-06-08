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
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.net.SocketException

class ClientFragment : Fragment() {
    companion object {
        private const val PICK_UPLOAD_FILE_REQUEST_CODE = 1
        private const val PICK_DOWNLOAD_DIR_REQUEST_CODE = 2
    }

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
            PICK_UPLOAD_FILE_REQUEST_CODE  -> readFileToUpload(data.data)
            PICK_DOWNLOAD_DIR_REQUEST_CODE -> {
                data?.data?.let { uri ->
                    startDownload(uri)
                }
            }
        }
    }

    private fun initFileList() {
        val recyclerView = binding.fileList
        val layoutManager = LinearLayoutManager(requireContext())
        recyclerView.layoutManager = layoutManager
        recyclerView.adapter = filesAdapter
    }

    private fun connect() {
        val host = Config(requireActivity()).getServerHost() ?: ""
        val input = Input(null)

        lifecycleScope.launch {
            val c = Client.newInstance(host, input)

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
        if (!this::client.isInitialized) return
        val channel = arguments?.getString("channel")
        if (channel != null) {
            client.channel = channel
        }
        readCID()
        readFiles()
    }

    private fun readCID() {
        if (!this::client.isInitialized) return
        val host = Config(requireActivity()).getServerHost()

        lifecycleScope.launch {
            val cid = client.readCID()
            binding.clientText.text = "Client #$cid @$host"
            binding.channelText.text = "Channel: ${client.channel}"
            binding.infoText.text = "Connected"
        }
    }

    private fun disconnect() {
        if (!this::client.isInitialized) return
        lifecycleScope.launch {
            client.disconnect()
        }
    }

    private fun chooseFileToUpload() {
        val intent = Intent(Intent.ACTION_OPEN_DOCUMENT).apply {
            addCategory(Intent.CATEGORY_OPENABLE)
            type = "*/*"
        }
        startActivityForResult(intent, PICK_UPLOAD_FILE_REQUEST_CODE)
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
                var chunksTotal = 0
                client.file = file
                client.upload(bytes) {
                    val percentage = it * 100
                    binding.infoText.text = "Uploading $percentage%"
                    chunksTotal++
                }
                handleFileUploaded(chunksTotal)
            }
            catch (e: SocketException) {
                println("ERROR: ${e.message}")
                handleConnectionFailed()
            }
        }
    }

    private fun handleFileUploaded(chunksTotal: Int) {
        binding.infoText.text = """
            File uploaded: ${client.file} | $chunksTotal chunks sent
        """.trimIndent()
        readFiles()
    }

    private fun download(file: String) {
        if (!this::client.isInitialized) return
        client.file = file
        chooseDownloadFolder()
    }

    private fun startDownload(uri: Uri) {
        lifecycleScope.launch {
            try {
                var chunksTotal = 0
                val array = client.download {
                    val percentage = it * 100
                    binding.infoText.text = "Downloading $percentage%"
                    chunksTotal++
                }
                write(requireContext().contentResolver, uri, array)
                handleFileDownloaded(chunksTotal)
                println("Downloaded ${array.size}")
            }
            catch (e: SocketException) {
                println("ERROR: ${e.message}")
                handleConnectionFailed()
            }
        }
    }

    private fun handleFileDownloaded(chunksTotal: Int) {
        binding.infoText.text = """
            File downloaded: ${client.file} | $chunksTotal chunks received
        """.trimIndent()
        readFiles()
    }

    private fun chooseDownloadFolder() {
        val intent = Intent(Intent.ACTION_CREATE_DOCUMENT).apply {
            addCategory(Intent.CATEGORY_OPENABLE)
            type = "*/*"
            putExtra(Intent.EXTRA_TITLE, client.file)
        }
        startActivityForResult(intent, PICK_DOWNLOAD_DIR_REQUEST_CODE)
    }

    private fun readFiles() {
        if (!this::client.isInitialized) return
        lifecycleScope.launch {
            val res = client.readFiles()
            val size = res.size
            binding.filesText
                .text = "${getString(R.string.files_title)} ($size)"
            files.clear()
            files.addAll(res)
            filesAdapter.notifyDataSetChanged()
        }
    }
}
