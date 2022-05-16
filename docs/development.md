# Development

## EOF Issues

EOF signals might be an issue in a Go implementation. They might be read as errors, or skipped, so take special care with this signal. It's defined at `io.EOF`.

When passing to state EOF the client sent the EOF signal I had defined in the beggining, that is, an empty chunk: `[]byte{}`. The problem is that the server never heard back from the client and the communication was on hold forever. The server didn't read that EOF signal, so it wasn't able to proceed with the next state.

This had me some time reasoning about some syncronization issue, but the problem is that the signal was not getting read by the server because it is an empty message literally. Go writes and reads keep in sync due to its good concurrency model to write linear code, so syncronization was not the problem. It might have been due to some validation logic I wrote or the way the empty message is handled.

I wanted to use the empty chunk as my EOF signal to keep reading chunks. This is because data is sent as chunks instead of actual high-level `Message`s. So I though it was ok to keep with that inertia and get the EOF signal just as the empty chunk.

In the end, I realized that was a primitive design. I chose to use raw chunks to remove the overhead when streaming files. The EOF is only sent onece, so I must use a proper high-level `Message`.

After implementing the design said above, all my headaches were immediately terminated.

There's no need to implement archaic primitive-obsession systems nowadays. You need less algorithms or tricks to handle that data. Replace stupid algorithms with well-defined domain specific centric systems.
