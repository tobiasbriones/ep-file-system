<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# Measure Cohesiveness

Cohesive modules do one thing and well done. It's about the same responsibility.
This is one of the most important engineering principles I teach.

Measure this by finding pieces of code that do not match the single
responsibility.

For example, the `Client` `struct` was getting a bit out of hands. I checked and
knew already that it had different type of implementations that are not directly
related to the module.

The key here is to check specifications. The `Client` is an object that works
directly with the `net.Conn` object, so it's a networking detail. If I have file
system, or IO inside there, then is the wrong place for that logic or
implementation.

Just look at this red flag:

```go
switch c.state {
    default:
    c.listenMessage()
    case Data:
    c.listenData()
    case Stream:
    c.listenStream()
    case Eof:
    c.listenEof()
    case Done, Error:
    log.Println("Exiting client")
    return
}
```

I had pushed the `process` domain into this object. That is basically the
definition of the FSM (something pure from the hardcore domain module), while I
must have a network implementation detail into this `Client` module instead.

Projects are never perfect, and they evolve from prototypes or initial
developments. I didn't make a mistake on writing that code that way. We just
need to write code that can be refactorized.

That answers a question I've seen on the internet: Do seniors write great code
from the beginning?

I can't over-engineer to tell that I write the best code since the beginning of
all the projects, I can't under-engineer eiter because that would turn into a
problem factory soon. A senior like me just knows what to do in each situation.

## Make the Right Thing Right

Incompetent programmers or engineers may say sentences like "duplication is 
better than the wrong abstraction". That is not an excuse to make things 
wrong. Something valid is to say "this is a prototype, just get it done" 
because prototypes are not meant to be correct, they are not engineered a 
lot, they're made more by frontend developers than actual engineers.

Some say you should write a prototype in a different language (a toy
scripting language like Ruby, Python, PHP or JS sure) than the final
language you will use to prevent reusing the prototype. This clarifies the
difference I emphasized above:

- to build a real system you have to make it right (as far as requires) without
  ridiculous excuses, and
- to build prototypes (most software out there) you don't have to mess with
  wrong abstractions, so you don't have excuses either.

I hope that insight had given you a better perspective to be a professional 
engineer who acts on behalf computer science rather than excuses and cringe 
marketing buzzwords like "WET", "DRY", ".NET", etc.
