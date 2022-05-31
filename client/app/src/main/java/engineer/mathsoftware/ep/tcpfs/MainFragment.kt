// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.LinearLayoutManager
import com.google.android.material.snackbar.Snackbar
import engineer.mathsoftware.ep.tcpfs.databinding.FragmentMainBinding
import kotlinx.coroutines.launch
import org.json.JSONException

const val PICKFILE_REQUEST_CODE = 1

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
            }

            try {
                val channels = client.readChannels()
                loadChannels(channels)
            }
            catch (e: JSONException) {
                println(e.message)
            }
        }
    }

    private fun subscribe(channel: String) {
        disconnect()
        findNavController().navigate(R.id.action_FirstFragment_to_SecondFragment)
    }

    private fun disconnect() {
        lifecycleScope.launch {
            client.disconnect()
        }
    }

    private fun loadChannels(values: List<String>) {
        channels.clear()
        channels.addAll(values)
        channelsAdapter.notifyDataSetChanged()
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
