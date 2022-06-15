// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-tcp-file-system

package engineer.mathsoftware.ep.tcpfs

import android.app.AlertDialog
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView

class ChannelsAdapter(
    private val dataSet: List<String>,
    private val l: (channel: String) -> Unit,
    private val del: (channel: String) -> Unit
) :
    RecyclerView.Adapter<ChannelsAdapter.ViewHolder>() {

    class ViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        val textView: TextView

        init {
            textView = view.findViewById(R.id.textView)
        }
    }

    override fun onCreateViewHolder(
        viewGroup: ViewGroup,
        viewType: Int
    ): ViewHolder {
        val view = LayoutInflater.from(viewGroup.context)
            .inflate(R.layout.text_row_item, viewGroup, false)
        return ViewHolder(view)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, position: Int) {
        val channel = dataSet[position]
        viewHolder.textView.text = channel
        viewHolder.itemView.setOnClickListener { l(channel) }
        viewHolder.itemView.setOnLongClickListener {
            val items = arrayOf<CharSequence>("Delete")
            val builder = AlertDialog.Builder(viewHolder.itemView.context)
            builder.setTitle("Channel: ${dataSet[position]}")
            builder.setItems(items) { _, _ -> del(channel) }
            builder.show()
            true
        }
    }

    override fun getItemCount() = dataSet.size
}
