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
problem factory soon. A "senior" like me just knows what to do in each situation.

## I Found Something Superb

Cohesiveness here is an **abstract computer science measure** unlike buzzwords
you hear from clowns like YouTubers, marketers, etc.

The Wikipedia article defines it as:

> In computer programming, cohesion refers to the degree to which the
> elements inside a module belong together. In one sense, it is a measure of the
> strength of relationship between the methods and data of a class and some
> unifying purpose or concept served by that class. In another sense, it is a
> measure of the strength of relationship between the class's methods and data
> themselves.
>
> Source: [Wikipedia](https://en.wikipedia.org/wiki/Cohesion_(computer_science))

That definition is too object-oriented, but it gives the idea.

Notice how is says "and *some* unifying purpose or concept served by that class"
so that the "*some*" is defined by your objectives or requirements. That's
why I said cohesiveness is an abstract measure unlike electrical
engineer formulas you find at universities.

That **powerful observation** I made proves that a software engineer (unlike
traditional engineers) needs to go through a **case by case basis** to design
the correct solution for each problem.

That's why we need autonomous engineers and not managers. Manager and its
bureaucratic derivatives are a smell, isn't it more efficient that an
autonomous engineer knows how to "manage" himself/herself for his/her agenda
rather than someone else do it for them? That leads to fragment the
knowledge from "subordinates" and "managers". Sounds like OOP and their
patterns, isn't it?

The main different between a traditional engineer and an autonomous software
engineer regarding correctness[^1] is that they have stable problems already
solved like the gravity value that doesn't significantly change on earth, but it
does change. Software engineers have to deep into the domain problem to come
out with a stable domain solution to that particular domain (anything you can
imagine as everything runs on software nowadays).

[^1]: Many don't consider SWE as real engineering because of lack of
    correctness but think about the prototype cases I mentioned in the section 
    "Make the Right Thing Right"

Building ordinary software doesn't make you an engineer at all but a software
developer who uses the tools made by engineers, engineering is about
correctness not about buzzwords or playing at the computer. It's cool, we're 
also developers, but they are just tools, not engineering as is.

I've been aware of all this, and then in sections like this I have the 
chance to develop it so others learn from me. Can you see? I'm an autonomous 
engineer, I don't need a professor with 10 degrees or a manager to tell what 
I have to tell!. Schools, capitalism, etc., are stereotypes with generic 
rules. We don't need them.

### Take Away

This is like a family in set theory: a set of sets. Or an undefined
integral: a family of functions. SWE principles tell you accurate
definitions of how to solve a family of problems so that for a particular
problem we come out with an accurate way to do the right thing right.

Unlike buzzwords that tell *generic* rules that are not backed by science
and math.

I always contrast math/FP/CS/formalism vs dev/OOP/buzzwords and the result 
is the same:

- accurate and domain specific vs generic and general purpose.
- real engineering vs clowns who get called engineers.

Most software is cheap as I say below, so don't confuse real engineering 
with daily toy development.

#### Taking About Clowns

Those with no skills, even marketers who can edit video and upload it to the
internet get called engineers by capitalist because that's how they
brainwash morons to make them believe they're "important" so they have happy
(non-important) employees.

Even if a job is complete trash like selling cable and internet under the
sun from door to door, they now call them "selling executives". Capitalists
are currently the root of most evil, is like Java: everything is money,
everything (but the 8 primitives) is an "object" (monolith). And guess what,
money is another *generic* stereotype to brainwash peons.

We can have many forms of energy or value instead of money.

Some say over there: "the best slaves are those who don't know they're
slaves", or "who are happy being slaves". So have happy slaves calling them
engineers when it does not even follow any formal definition of what an
engineer is supposed to be.

I know many juniors can grow to become engineers, but a junior engineer !=
engineer (that's for another day). You have to take into account the domain
of engineering too.

The definition for engineering is clear, briefly and simple put is:
building goods for society with correctness, economical, social, legal, etc.
constraints. The key is that for **correctness** *engineering* must be derived
directly from sciences and math (even more SWE), or in my words: engineers
must be good at direct knowledge.

Thus, not everyone can be called an "engineer".

Many look down software engineers like clowns who learned to code in three
months by watching videos. We don't have anything to do with clowns.
Software engineering is powerful, is beautiful, is the engineering that
requires the most math and creativity. Others like electrical engineers
require the most physics for example. Engineering is about direct knowledge.

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
