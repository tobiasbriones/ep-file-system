# Development

## EOF Issues

EOF signals might be an issue in a Go implementation. They might be read as errors, or skipped, so take special care with this signal. It's defined at `io.EOF`.

When passing to state EOF the client sent the EOF signal I had defined in the beggining, that is, an empty chunk: `[]byte{}`. The problem is that the server never heard back from the client and the communication was on hold forever. The server didn't read that EOF signal, so it wasn't able to proceed with the next state.

This had me some time reasoning about some syncronization issue, but the problem is that the signal was not getting read by the server because it is an empty message literally. Go writes and reads keep in sync due to its good concurrency model to write linear code, so syncronization was not the problem. It might have been due to some validation logic I wrote or the way the empty message is handled.

I wanted to use the empty chunk as my EOF signal to keep reading chunks. This is because data is sent as chunks instead of actual high-level `Message`s. So I though it was ok to keep with that inertia and get the EOF signal just as the empty chunk.

In the end, I realized that was a primitive design. I chose to use raw chunks to remove the overhead when streaming files. The EOF is only sent onece, so I must use a proper high-level `Message`.

After implementing the design said above, all my headaches were immediately terminated.

There's no need to implement archaic primitive-obsession systems nowadays. You need less algorithms or tricks to handle that data. Replace stupid algorithms with well-defined domain specific centric systems.

### Example/Take Away

In my another [file system](https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system) written in Java I had to standardize the path separator because I don't work for M$ Windows or Linux or MacOS. 

I had my own file system, so I don't care if M$ marketers decide to use a cringe symbol that has to be escaped all the time for Windows paths just to brand it as Microsoft Windows only. 

I do direct math, computer science, etc. So my domain is math, computer science, etc. My domain is about doing the right thing right, so most of the time I choose to build the DSL instead of relying on generic/primitive/platform-dependent standards.

Now when coming to performance everything turns obscure and hell (imagine those "algorithms" written in C/C++ for high-performance models). That's because computer hardware is general-purpose (that is good), and imperative. If you want a performant high-level "algorithm" well done, your computer hardware will tell you: we don't do that here.

Fortunately, performance is not the most important tradeoff, and computer hardware has evolved lately to implement better architectures like ARM, or SIMD. That way hardware gets more familiarized with linear algebra and math, and it's more efficient.
