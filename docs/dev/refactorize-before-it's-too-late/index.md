<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# Refactorize Before It's Too Late

I often have to take significant time to refactor the code base because I know
from experience that the more issues that are carried the more cyclomatic
complexity and expensive further developments will be.

They can't fire you if you do "useless" things like refactoring or writing
better code, an autonomous engineer knows what to do, so don't let managers or
stakeholders tell you what to do, you have communication skills and technical
ones, but managers don't have any other skill than being people friendly if at
all.

A professional developer or engineer will put system constraints clearly, don't
be shy on doing things right. Shame on them if they don't want it so, capitalist
just want to sell, and they might think this is like traditional engineering
where you have a table with formulas, so you can deliver projects "on time".
It's annoying, they kind of tell you "I'm paying you to build this website, so
you have 5 days to finish it", but in reality software engineering is about
being a *domain expert* to understand the domain problem and then being a
*developer* to build the domain solution.

## My Domain

Just check how I perfectly applied this expertise on me: I'm a software engineer
with domain on math, so that is, a math software engineer. That means I can also
be a "mediator" (like OOP programmers would say) between mathematicians and
software engineers to build mathematical software or any other related project.

Computer scientists are bad at engineering and engineers are terrible at math or
anything where you need intellectual conscious skills. If I wouldn't exist then
mathematicians would be doomed to write boring proofs in PDF files (no
indexable, useless nowadays), using proprietary crap software from the 90s, and
nonsense horrible "Alexandria math library" (they even explicitly say their "
math lib" is the hugest monolith and feel proud about it, they don't know what
they do because they're building tools but are not engineers) while engineers
are doomed to write toy scientific software with wrong tools like Python or
Microsoft Excel.

Once again I prove we don't need managers and buzzwords but autonomous engineers
instead.

## Large Refactorization

If you don't refactor constantly the code you find, the project cost will
explode really quick turning into an extremely coupled system.

Large refactorizations are pretty tiring, is something that will make you
exhausted because of the large cognitive load to keep tests working and the
older code working too.

Another tip that I can give is not to underestimate the initial system design
which will avoid making large refactorizations later.

When designing system architectures we need a domain expert (e.g. someone who
has build file systems before), if I would've thought a little more about the
modules of this application the cost of development would've been lower since I
would have scaled those packages or modules since the beginning instead of
making larger refactorizations later.

Some systems are not well known, keeping the balance between under-engineering
and over-engineering is a determinant art in software engineering.
