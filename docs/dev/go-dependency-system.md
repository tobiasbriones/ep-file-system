<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# Go Dependency System

Something I have detected is that Go dependencies are linear to avoid circular
dependencies.

This might be a mess but is a good design decision in the end.

It can be a mess because, for example, in the dungeon game, I have the server
module and the client or game module. I need to extract the game model logic,
but I couldn't. I had to do something gross to move forward:
copy-paste the game model to the game module and to the server module.

This must be because Go is aimed for microservices and each module has to be
small and most important undependable deployable.

That's why I had said that you can't build monoliths with Go.

## Where's the Domain Model Then?

Everything is relative, as I
[recently wrote](https://blog.mathsoftware.engineer/everything-is-relative).

The most important layer or module is the app domain. With this dependency
model, it must go all the way down so the "dependency arrow" goes all the way
down. The model won't be able to see implementation details and keeps pure this
way.

I hate it because the domain must be the less verbose system, and must be in the
top of the project. The problem is that Go should be used for
independently-deployable modules.

With Java, I would create different Java/Gradle modules in the root of the
project and import them as required. I can get to be circular, but it's not too
bad. That way I can build *modular* monoliths that are pretty useful in
development. Check the
[other FS](https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system)
source code for seeing this.
