# Development

## EOF Issues

EOF signals might be an issue in a Go implementation. They might be read as errors, or skipped, so take special care with this signal. It's defined at `io.EOF`.

When passing to state EOF the client sent the EOF signal I had defined in the beggining, that is, an empty chunk: `[]byte{}`. The problem is that the server never heard back from the client and the communication was on hold forever. The server didn't read that EOF signal, so it wasn't able to proceed with the next state.

This had me some time reasoning about some syncronization issue, but the problem is that the signal was not getting read by the server because it is an empty message literally. Go writes and reads keep in sync due to its good concurrency model to write linear code, so syncronization was not the problem. It might have been due to some validation logic I wrote or the way the empty message is handled.

I wanted to use the empty chunk as my EOF signal to keep reading chunks. This is because data is sent as chunks instead of actual high-level `Message`s. So I though it was ok to keep with that inertia and get the EOF signal just as the empty chunk.

In the end, I realized that was a primitive design. I chose to use raw chunks to remove the overhead when streaming files. The EOF is only sent onece, so I must use a proper high-level `Message`.

After implementing the design said above, all my headaches were immediately terminated.

There's no need to implement archaic primitive-obsession systems nowadays. You need less algorithms or tricks to handle that data. Replace stupid algorithms with well-defined domain specific centric systems.

Sending the chunks as raw byte arrays is good because a file is just that, a bunch of bytes without any structure. So there's no problem with that design decision. It also avoids the extra overhead of sending a `Message` as long as the FSM states are valid in the client and server, that is, state `DATA` (upload) or `STREAM` (download).

### Example/Take Away

In my another [file system](https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system) written in Java I had to standardize the path separator because I don't work for M$ Windows or Linux or MacOS. 

I had my own file system, so I don't care if M$ marketers decide to use a cringe symbol that has to be escaped all the time for Windows paths just to brand it as Microsoft Windows only. 

I do direct math, computer science, etc. So my domain is math, computer science, etc. My domain is about doing the right thing right, so most of the time I choose to build the DSL instead of relying on generic/primitive/platform-dependent standards.

Now when coming to performance everything turns obscure and hell (imagine those "algorithms" written in C/C++ for high-performance models). That's because computer hardware is general-purpose (that is good), and imperative. If you want a performant high-level "algorithm" well done, your computer hardware will tell you: we don't do that here.

Fortunately, performance is not the most important tradeoff, and computer hardware has evolved lately to implement better architectures like ARM, or SIMD. That way hardware gets more familiarized with linear algebra and math, and it's more efficient.

## Domain Model

I had some loose file system logic, some weird data structures, and the utility functions. This is a design smell, so I decided to build a basic domain model based on my other [file system](https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system/tree/main/model).

This basic model just adjusts to the project requirements, so I don't have to implement other features like `getParent`, or tree traversal, etc. 

The Android and Web clients need to read these models too, building the base model for the file system (`File`, `Directory`, `Node`, operations like file name/extension, etc.) would be insanely expensive.

To address this problem, the domain model in the server has to be adecuate for this project since it's the most important module. The clients can read just primitive data types to consume the content.

This way, the loose logic is correctly coupled, and testable.

I always test the domain modules as much as I can to avoid propagating errors forward to presentation layers, and avoid system downtimes or failures.

## Enums in Go

Go is a simple language, so it's often underengineered. I think about Go as if Python was a real language, or Python well done.

The states of the FSM need to be well defined constants, sum types like `enum`s.

### Iota a Bad Trick

`iota` is something too implicit, you just define groups of `const` and they get their int value from there on starting from zero.

If you change the order of definition then your enum value is going to change, and that is a huge problem for backwards compatibility.

```go
type State uint

const (
	Start State = iota
	Data
	Stream
	Eof
	Error
	Done
)
```

You may say, `iota` is not a trick but a feature. I don't think is a proper feature. Can you see how fragmented this is going?

Tutorials on the internet tell you to use `iota`, so it looks like it's idiomatic Go.

### Lack of Features

It's like other lame languages like JavaScript or Python, you can't even define a simple `enum` data type.

I know Go is a simple language for concurrent applications and microservices. Easy to adopt and move forward on large teams where programmers are coming and going (Google). But for doing mediocre things I rather use Java. You can't build monoliths on Go by the way, so Go is just a niche language, it is not general-purpose either.

On the other hand, I'd rather write Go than any horrible languages like Python, PHP, or JavaScript, having into account the bloated communities and tools they have. If you don't write static types, then you need a lot of bloated buggy software like Anaconda, Electron, etc. It's a tradeoff, why use those languages then?
