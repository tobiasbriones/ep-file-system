<!-- Copyright (c) 2022 Tobias Briones. All rights reserved. -->
<!-- SPDX-License-Identifier: BSD-3-Clause -->
<!-- This file is part of https://github.com/tobiasbriones/ep-tcp-file-system -->

# A Development Cycle Can Be Shorter

Analysis this project as case study, I had to release an initial alpha `0.1.0`
version and then a stable `0.1.0` version for an MVP.

The development was quite agile, and I was able to constantly deliver pull
requests for docs and dev, and also making releases in short periods of time.

Can these releases be actually used in production environment according to
Agile? Some of them just can't as they're not stable.

I usually think about the iterative model where I can create a working bicycle
and add more each iteration, unlike the incremental model that builds the parts
instead of the whole.

The initial development went a bit more monolith because I decided from the
beginning to use chunks of bytes to transfer the files, otherwise the system
would be badly designed. That big feature made the initial development more
monolithic and far from a stable-deployable MVP.

If I had to create this project again I would defer that feature for later and
send the whole files at once for the initial releases.

Adding too many features at once also increases the uncertainty, cyclomatic
complexity, and early testing. That makes refactorization more painful and hard
to spot.
