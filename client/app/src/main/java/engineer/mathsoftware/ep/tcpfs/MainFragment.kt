// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.app.AlertDialog
import android.os.Bundle
import android.text.InputType
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.core.os.bundleOf
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.LinearLayoutManager
import com.google.android.material.snackbar.Snackbar
import engineer.mathsoftware.ep.tcpfs.databinding.FragmentMainBinding
import kotlinx.coroutines.launch
import org.json.JSONException

class MainFragment : Fragment() {
    private val channels = ArrayList<String>()
    private var _binding: FragmentMainBinding? = null
    private lateinit var channelsAdapter: ChannelsAdapter
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
        channelsAdapter = ChannelsAdapter(channels) { subscribe(it) }
        binding.fab.setOnClickListener { showCreateChannelDialog() }
        initChannelList()
        connect()
    }

    private fun initChannelList() {
        val recyclerView = binding.channelList
        val layoutManager = LinearLayoutManager(requireContext())
        recyclerView.layoutManager = layoutManager
        recyclerView.adapter = channelsAdapter
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
        disconnect()
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
                handleConnected()
            }
        }
    }

    private suspend fun handleConnected() {
        if (!this::client.isInitialized) return
        try {
            val channels = client.readChannels()
            loadChannels(channels)
        }
        catch (e: JSONException) {
            println(e.message)
        }
    }

    private fun subscribe(channel: String) {
        disconnect()
        val bundle = bundleOf("channel" to channel)
        findNavController().navigate(
            R.id.action_FirstFragment_to_SecondFragment,
            bundle
        )
    }

    private fun disconnect() {
        if (!this::client.isInitialized) return
        lifecycleScope.launch {
            client.disconnect()
        }
    }

    private fun loadChannels(values: List<String>) {
        channels.clear()
        channels.addAll(values)
        channelsAdapter.notifyDataSetChanged()
    }

    private fun showCreateChannelDialog() {
        val builder: AlertDialog.Builder = AlertDialog.Builder(requireContext())
        val view = LayoutInflater.from(requireContext())
            .inflate(R.layout.dialog_input_text, null)
        val input = view.findViewById<TextView>(R.id.dialogInputText)
        input.setHint("Enter channel name")
        input.inputType = InputType.TYPE_CLASS_TEXT
        builder.setTitle("Create Channel")
        builder.setView(view)
        builder.setPositiveButton("CREATE") { _, _ ->
            var channelName = input.text.toString()
            createChannel(channelName)
        }
        builder.setNegativeButton("Cancel") { dialog, _ -> dialog.cancel() }
        builder.show()
    }

    private fun createChannel(channelName: String) {
        if (!this::client.isInitialized) return
        lifecycleScope.launch {
            client.createChannel(channelName)
            Snackbar.make(
                requireView(),
                "Channel created",
                Snackbar.LENGTH_LONG
            )
                .show()
            handleConnected() // Reload
        }
    }
}

fun Fragment.handleConnectionFailed() {
    Snackbar.make(
        requireView(),
        "Fail to connect",
        Snackbar.LENGTH_LONG
    )
        .setAction("Action", null)
        .show()
}
