---
kind: article
created_at: 2014-11-21
layout: post
title: "Dependency management in C"
tags: [c, packaging, buzzy]
---

At [RedJack](http://www.redjack.com/), all of our [core
products](http://www.redjack.com/solutions/) depend on a network sensor that
collects various bits of information about the raw traffic that we see on the
network.  We're doing some non-trivial analysis on fairly large network links
using commodity hardware, so we've implemented this sensor in C.  At its core is
an extremely fast custom [flow-based
programming](http://en.wikipedia.org/wiki/Flow-based_programming) framework.
It's a damn cool piece of code, but this post isn't about the code itself; it's
about how we deliver that code to our customers.

Just because we've written this component in C, that doesn't mean we want to
turn our back on the kinds of tooling you get to use when working in other, more
modern languages.  In particular, once you've gotten used to modern package
managers like [npm](https://www.npmjs.org/), [leiningen](http://leiningen.org/),
[go](http://golang.org/doc/articles/go_command.html), and
[Cargo](http://doc.crates.io/guide.html), it's hard to go back to things like
[CMake](http://www.cmake.org/) and *[shudder]* the
[autotools](http://en.wikipedia.org/wiki/GNU_build_system).

<hr class="jump">

And the worst part is that those examples aren't even comparable.  Very little
code exists in isolation these days; in addition to the actual source code that
you've written (in whatever language you've chosen), you're going to have
dependencies on other libraries.  All of the npm-like package managers are
full-featured systems that handle **dependency management** as well as the
actual **building** of your code.  CMake and the autotools, on the other hand,
only focus on the build step — they assume that all of your dependencies are
already available.

At this point, C (and C++) developers usually throw up their hands and give up.
You either try to avoid third-party dependencies at all, by embedding
everything you need so that you can fold the "dependency management" step into
your existing build tool.  Or you list your dependencies in a human-readable
`INSTALL` document in the source tree, and rely on the ingenuity and elbow
grease of your users to track down and install all of the prereqs by hand.  If
you're lucky, and your project is popular enough to include in the major Linux
distros' package repositories, then the various package maintainers will
re-encode all of this information in the distro-specific package definitions.
But because there's no single, standardized packaging specification for Linux
(let alone for other POSIX-y operating systems), that means someone has to do
that separately for each distro that you need to support.

So we created [Buzzy](https://github.com/redjack/buzzy/).  It gives us npm-style
dependency management for our C projects, and — almost as a side effect! — makes
it easy to generate native packages of our software for all of the Linux distros
that we need to support.  The end result is that it's dirt simple to build and
install [one of our libraries](https://github.com/redjack/varon-t/), along with
all of its prerequisites:

    $ git clone git://github.com/redjack/varon-t.git
    $ buzzy install -q
    [1] Clone git://github.com/redjack/check.git (master)
    [2] Clone git://github.com/redjack/libcork.git (master)
    [3] Clone git://github.com/redjack/clogger.git (master)
    [4] Clone git://github.com/redjack/buzzy-core.git (master)
    [5] Install native Arch package pkg-config 0.28+rev.2
    [6] Install native Arch package check 0.9.14
    [7] Install native Arch package cmake 3.0.2
    [8] Build libcork 0.14.1 (cmake)
    [9] Stage libcork 0.14.1 (cmake)
    [10] Package libcork 0.14.1 (pacman)
    [11] Install libcork 0.14.1 (pacman)
    [12] Build clogger 0.3.1 (cmake)
    [13] Stage clogger 0.3.1 (cmake)
    [14] Package clogger 0.3.1 (pacman)
    [15] Install clogger 0.3.1 (pacman)
    [16] Build varon-t 1.0.2 (cmake)
    [17] Stage varon-t 1.0.2 (cmake)
    [18] Package varon-t 1.0.2 (pacman)
    [19] Install varon-t 1.0.2 (pacman)

That single `buzzy` command works unchanged on Arch Linux, RHEL 6 and 7, any
recent Debian or Ubuntu, and Mac OS X (via Homebrew).

<br>

Buzzy has some interesting characteristics that are definitely different than
other bits of C project tooling — we purposefully tried to take into account the
best practices that have emerged in other modern package managers.  Of course,
there are also some idiosyncrasies that come from dealing with C projects.


##### Declarative package descriptions

With most existing C build tools, you provide a **script** that describes how to
build your project.  With Buzzy, you provide a [declarative
description](https://github.com/redjack/libcork/blob/develop/.buzzy/package.yaml)
of the package, and let Buzzy decide what to do with that description for your
particular platform.


##### Decentralized package descriptions

Buzzy does not have a single centralized database of package descriptions.
Instead, you put your package description into the source repository as the code
itself.  Part of your package description is a list of
[links](https://github.com/redjack/libcork/blob/develop/.buzzy/links.yaml) to
other Buzzy-enabled source repositories that you depend on.


##### Just handle dependencies and packaging

Buzzy is not a build tool; it delegates to something else to actually build the
source code for any particular project.  We happen to use CMake for all of our
libraries, and Buzzy also works fine with the autotools.


##### Use the native ecosystem wherever possible

In part, this means that we produce native packages; we haven't invented a new
packaging format or anything like that.  It also means that if a dependency is
available as a native package, we install that, instead of building our own.  We
only fall back on building our own copy for those cases (I'm looking at you,
RHEL!) where the native package is too old to satisfy a versioned dependency.


##### Auto-detect as much as possible

Wherever we can, we try to figure something out for you instead of making you
state it explicitly.  If we see a `CMakeLists.txt` file in your source
repository, we assume that your project is built using CMake.  Moreover, we know
how to build CMake projects without you having to provide an explicit set of
commands to run or options to pass in to those commands.

Buzzy also knows the most common patterns for how the different distros name
their native packages.  So if your package has a build dependency on `jansson`,
we know to look for a native package called `jansson` or `libjansson` on Arch,
`jansson-devel` or `libjansson-devel` on RPM derivatives, and `jansson-dev` or
`libjansson-dev` on Debian derivatives.

And of course, there's a way to override this auto-detection if Buzzy doesn't
know best for a particular situation.  (Because we live in a world where `zlib`
is called `zlib1g` on Debian systems.)


<br>

It's important to note that Buzzy is not a panacea.  It doesn't know, for
instance, that RPM and Debian derivatives want you to make separate runtime and
developer packages for your libraries.  Right now Buzzy just blindly builds a
single monolithic package regardless of what distro you're on.  But, it's served
us well so far, and it's definitely a good start!
