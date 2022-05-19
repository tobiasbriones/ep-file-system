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

#### Fun Fact

Primitive Android programmers told you not to use `enum`s in Java because of the overhead. That's because Android phones were super slow and limited back in that time.

We don't have to use primitive data types, or non-cohesive algorithms as I mentioned before. Low-level tricks is for bad programmers, that's part of the past.

That was important at that time, but now is just a sad joke.

We have to evaluate tradeoffs pretty well.

### More than Simple Enums

For this statically typed system I need to send the states via a TCP `Message`.

Enums can be easily defined via integers (`iota`) in Go, but they don't have a string representation. That implies that if I change the order of definition then the whole system will immediatelly break and is also a troble for backwards compatibility.

If I use strings to define the `enum`s I lose the advantage of switching on integer IDs. So I had to define a parallel string array to convert the value into string.

This is too bad:

![Commit Fix Enum Strings](commit-fix-enum-strings.png)

The enums are now fragmented, e.g. you have to manually keep the string representation in sync with the bunch of `const`s defined above. That's pretty lame.

Whenever I face these kinds of problems in Go, I have to ask the question: How to solve this in a simple way?. Since that is the way Go is supposed to be.

But simple is not underengineered though.

Messages has to be read into the application memory as programming language constructs or abstractions, but when sending them over the network they're only raw bytes. I heard the new Web 3 standard will fix that and that we'll send objects instead of JSON or bytes over the network. I don't know about that information, but I hope so!.

Transforming DTOs, raw data types, all this is too exhausting, and only shows lack of modern tech. It also adds incorrecness in the way.

What should be sent over the network?. Integers or strings to represent the `enum` values?. Integers use to be physical implementations. I can't add another `enum` because I can't tell wheter a `0` is a state from `FSM1` or `FSM2` if that value comes as raw from the network. What was the client's original intention?.

I initially used strings for the enum values, but then I have to send strings. Then I have to read an action as a string from the client `Message`, and to convert that raw string I need a `switch` to check it's a valid state, or else the not-so-clever string array to use integer indices but get index out of bound panics anytime I update the enum and spend 40 minutes debugging nonsense.

I really want to avoid that fragmentation and `switch`es.

With integer inices I can easily check if the value is valid too.

### What About Underengineering?

Go like many popular languages are for underengineering, for ordinary software written by ordinary programmers.

I can build underengineered systems of course, but the problem is that next move you realize you have to debug nonsense that can be easily avoided.

Robert C. Martin says, and I repeat all the time so others understand my pain: "The only way to go fast is to go well". Those phrases are the only boilerplate I love repeating. Java boilerplate is useless but these sentences are gold for me to defend my position as a professional engineer against annoying "stakeholders" or managers. They can't just fired us by writing working and tested software.

I got a better phrase for this, it looks stupid but you have to tell this to people because people only understands obvious things by recalling them all the time: "the only way to do things right is to do them right", there is no shortcut.

What can I do with a moutain of unmaintainable software that gets more complicated and coupled each time? From my experience, I have write refactorizable code to avoid underengineering and overengineering. Refactorize later as required. That's it.

It would be great to use Rust but companies without footprint and massive scalability issues will prefer to use a simpler language like Go with GC.

Debugging is a skill for bad programmers. It's like comments, the more you debug, the more issues there are in the underlying code. In years of writting software, the only debbuging skill I have is to create breakpoints with `println` to trace program states, and saving logs for production troubles. We don't need messy debugging tools, that's not to be proud about. Rust gives you the information right after compiling, and we can use Rust with easier scripting languages like TS (e.g. Deno is written in Rust, but the consumer language is TS), see my point?.

I just can't build underengineered systems because of my professionalism, but can't have much overthinking when using most languages like Go, TS, etc. Thus the answer to this section is to find a good balance. SWE is all about solving dynamic problems, is all about tradeoffs like how much to design in this part of the system?.

### Final Design

I have found the following design as best for Go enums:

```go
type State string

const (
	Start  State = "START"
	Data   State = "DATA"
	Stream State = "STREAM"
	Eof    State = "EOF"
	Error  State = "ERROR"
	Done   State = "DONE"
)

var stateStrings = map[string]struct{}{
	"start":  valid,
	"data":   valid,
	"stream": valid,
	"eof":    valid,
	"error":  valid,
	"done":   valid,
}

func ToState(value string) (State, error) {
	if _, isValid := stateStrings[value]; !isValid {
		return "", errors.New("invalid state value: " + value)
	}
	return State(value), nil
}
```
