---
kind: article
created_at: 2017-10-29
layout: /post.html
title: "Shared library versions"
description: sharing is caring
---

Congratulations, you've written a software library!  You hope that lots of
people will find it useful, and will take it as a dependency when writing their
own software.  You know that at some point you'll have to make changes to your
library, either to add features or to fix bugs.  Being a good maintainer, you
want to be as diligent as possible in telling your users what to expect as you
publish these changes.  Will they need to change their code in response to the
changes that you've made?  Have you retired features that they depend on?  Or
are the changes "safe", presumably requiring no updates on their part?

The traditional approach is to encode all of this information into an
easy-to-digest **version number**.  Of course, nothing in this world is simple,
so there are a number of different systems for encoding compatibility
information into a version number.  And surprisingly, if you're writing a shared
library for a compiled language like C or C++, there are (at least!) two
different versioning systems that you'll need to learn.  In this post, we're
going to look at these different systems, how they relate to each other, and how
to actually apply these version numbers to your library using a couple of common
build tools.

<hr class="jump">

#### Shared libraries: Why?

Before we jump into version numbers, let's talk about shared libraries and why
you might need them.  Shared libraries aren't really a thing in many modern
programming languages.  You don't need them, for instance, in programming
languages which aren't compiled (Python, Ruby), or where the convention is to
always compile static binaries (Go, Rust).  And so knowing how to [manage
them][drepper] has become a bit of a lost art.

[drepper]: https://www.akkadia.org/drepper/dsohowto.pdf

But in the glorious old world of C and C++, shared libraries are still very much
a thing.  Shared libraries provide two main benefits:

- **They reduce the amount of time you have to spend compiling your code.**

  A popular Rust crate, for instance, will be compiled every time that you build
  one of the projects that depends on it.  That's not necessarily a bad thing —
  it means you can optimize the library differently for each project that uses
  it.  But it does come at a cost in [developer time][xkcd compile] and in CPU
  cycles.

  [xkcd compile]: https://xkcd.com/303/

  With a shared library, you compile the library once, and install it into a
  shared location in the filesystem (typically `/usr/lib` on Linux systems).
  Any project that depends on that shared library can use that shared, already
  compiled representation as-is.

  Most Linux distributions further reduce compile times by distributing **binary
  packages** of popular libraries, where the distribution's packaging system has
  compiled the code for you.  By installing the package, you download a
  (hopefully signed) copy of the compiled library, and place it into the shared
  location, all without ever having to invoke the compiler (or any other part of
  the build chain that produced the library).

- **They reduce memory usage at runtime.**

  This part is less well known, I think.  At runtime, the linker and the
  operating system conspire to make sure that any individual shared library file
  is only loaded into memory once, no matter how many running processes are
  using it!  The OS will expose that single copy of the library code to each
  process that uses it, using the magic of virtual memory to ensure that all of
  those virtual "copies" are actually backed by the same chunk of physical
  memory.  Since the code has already been compiled, it's marked read-only,
  making this sharing safe.  (No process can overwrite any of the library code,
  invalidating it for the other processes sharing that copy.)

Both of these benefits are proportional to the number of projects that depend on
the shared library, so they're especially useful for "core" libraries that are
used by a lot of other software packages — e.g., the low-level GUI toolkit
provided by [GTK+][], or the [standard C library][glibc].

[GTK+]: https://www.gtk.org/
[glibc]: https://www.gnu.org/software/libc/

Note that using a compiled language like C doesn't **force** you to use shared
libraries; you're free to compile your dependencies along with your own project,
and link everything together **statically**.  (If you're going to distribute
your project via a container, you might as well link statically, since we
typically don't share code across containers at runtime!)

#### Version numbers

Right, so you're writing some code, which will be compiled into a shared
library, and you want to attach a version number to play nice with your users.
What do you do?

Your job is to think about the **public API** ("application programming
interface") of your library — that is, the promises you make _at the source code
level_ about how other programmers can use your library.  Which functions and
classes are available, what their signatures are, that kind of thing.

All of your changes will fall into one of three categories:

- **backwards incompatible changes**

  These occur when you *change* part of your public API, or *remove* something
  from it.  That means that some user of your library is going to have to
  **change their code** as a result of this release.  You should avoid this if
  you can, since it places more of a burden on your users; when it's
  unavoidable, you need to communicate this clearly so that your users know what
  to expect.

- **backwards compatible changes**

  These occur when you *add* something to your public API.  All of your existing
  users will be able to use the old release or the new release, without having
  to change any of their code.  These kinds of releases are great, since they
  introduce new functionality without adding an extra burden to your users.

- **bug-fix releases**

  These occur when there's no change to your public API at all.  All of your
  existing users should be able to use the old release or the new release,
  without having to change any of their code.  They will presumably want to
  upgrade to this new release at their earliest convenience, what with all of
  those bugs that you just fixed!  (Unless they're depending on the buggy
  behavior?  Never...)

Version numbers can be consumed either by humans ("is this upgrade going to be
annoying and dangerous?") and by computers ("my package manager will upgrade
this for me automatically if it can"), but either way, their goal is to
succinctly describe what kinds of changes you've made to your code from release
to release.

Historically, there have been many competing "patterns" for constructing a
version number from your list of changes.  These days, however, most people have
settled on [Semantic Versioning][semver] (semver) as the best set of rules, with
some languages going so far as to [mandate it][rust semver].

{::comment}
(To be fair, there has also been the [inevitable backlash][semver antipattern]
against semver; though that seems to be less about whether semver is a good
pattern, and more about whether you should use *any* kind of version number.)
{:/comment}

[semver]: http://semver.org/
[rust semver]: http://doc.crates.io/specifying-dependencies.html
[semver antipattern]: https://surfingthe.cloud/semantic-versioning-anti-pattern/

Under semver, a version number consists of three numbers: a **major version**, a
**minor version**, and a **patch level**.  Each of these corresponds to one kind
of change, and each time you cut a new release, you "bump" the portion that
lines up with the "strongest" change that you've made to the public API.  Any
backwards incompatible changes?  Bump the major version, set minor and patch to
0.  Only backwards compatible changes?  Bump the minor version, set patch to 0.
No changes at all?  Bump the patch level.  This intuitively lines up with what
many project maintainers were doing anyway; semver just codifies that behavior
as an explicit set of rules.

{::comment}
What about bug fixes?  Most versioning schemes describe how the library is
**supposed** to work.  If we discover a bug in the library, and we publish a fix
for that bug, then we want to convey that we haven't made any changes at all to
the API of our library.  Consumers of our library should be able to update to
the latest version, which includes the bug fix, and things will be magically
"better".

Of course, if you're responsible for an application that depends on this
library, you probably want to test this instead of blindly adopting the new
version.  

Worse, what if your application depends on the buggy behavior?  Does fixing the
bug mean 
{:/comment}

#### APIs and ABIs

Congratulations!  We've solved library versioning once and for all!

Well, no, not really.  Compiled shared libraries complicate this situation,
because you also have to consider the library's **ABI** (application *binary*
interface).  In a compiled language, you can make a change to your library,
which does *not* require your users to change their code at all, but which still
means they can't use an existing compiled version of your library as-is!

What would this look like?  Let's say you've written a C library for tracking
Aussie football games, and have the following struct in version 1.3.0:

<% highlight :c do %>
struct score {
  int goals;
  int behinds;
};

int total_score(const struct score) {
  return score.goals * 6 + score.behinds;
}
<% end %>

And some other programmer has written some code that uses this struct:

<% highlight :c do %>
struct score adelaide;
adelaide.goals = 11;
adelaide.behinds = 14;
printf("Adelaide %d.%d (%d)\n",
       adelaide.goals, adelaide.behinds,
       total_score(adelaide));
<% end %>

Running this code, your user would get:

```
Adelaide 11.14 (80)
```

Now for some reason, you decide to reorder the fields in your struct:

<% highlight :c do %>
struct score {
  int behinds;  /* behinds come first now! */
  int goals;
};
<% end %>

Your user's code is still perfectly valid!  If you only consider the source API,
since your user doesn't have to make any changes to their code, semver calls
this a bug-fix change.  When it's time to release the new version of the library
including this change, you would bump the version number from **1.3.0** to
**1.3.1**.

But if your user doesn't recompile their code, they'll get the wrong answer when
they run their program using the new version of your library:

```
Adelaide 11.14 (95)
```

Because the two pieces of *compiled* code had different assumptions about how
the fields in `struct score` were laid out in memory, they were incompatible,
even though the original source code was fine!

All of this means that if you're working with a compiled language and shared
libraries, you should consider the **compiled ABI** as well as the **source
API** when deciding what kinds of changes are included in a release.  You should
consider this field-reordering example a backwards-incompatible change, and bump
your version from **1.3.0** to **2.0.0**.

#### Shared library filenames

Now that we've talked about version numbers in the abstract, what do we see in
practice?

On Linux and Mac machines, you will also encode version numbers into the
filenames of your shared libraries.  Each library will end up with **two or
three** different filenames under `/usr/lib`.  For instance, if you have version
1.2.3 of a shared library called _libfoo_ installed, you'll find:

    $ ls /usr/lib/libfoo*
    /usr/lib/libfoo.so
    /usr/lib/libfoo.so.1
    /usr/lib/libfoo.so.1.2.3

Note that there aren't three *copies* of the library; the first two filenames
will be symlinks to the last one.

The `libfoo.so` file is only used at build-time.  (In fact, Debian-based systems
will only include this file in the library's `-dev` package; if you don't have
that package installed, you'll only see the last two versioned filenames.) You
compile some code that uses _libfoo_ by passing in `-lfoo` to your build tools.
But when you do this, the build tools don't know in advance which version of the
library is installed.  Instead of doing some kind of wildcard match, looking for
all filenames that match a pattern, the build tools assume that they can find
the library with a simple `libfoo.so` filename.  It's up to you (or more
realistically, your package manager) to make sure that this points at the
currently installed version.

The `libfoo.so.1` file is used at runtime.  By convention, this base filename,
which only includes the *major version* of the library, is called the library's
**SONAME**.  On a Linux system, you can see the SONAME of a library using the
`objdump` command:

    $ objdump -x /usr/lib/libfoo.so | grep SONAME
      SONAME               libfoo.so.1

When you compile some other code that depends on this library, the build chain
will find the shared library file using the non-versioned `libfoo.so` filename,
extract the library's SONAME from that file, and record that as the dependency.
You can see these dependencies using `objdump`, too:

    $ objdump -x /usr/bin/foo | grep NEEDED
      NEEDED               libfoo.so.1
      NEEDED               libc.so.6

So, at runtime, when `/usr/bin/foo` is loaded, the dynamic linker will see these
`NEEDED` entries, and look for a library file called `libfoo.so.1`.

The last file, with the full version number included in the filename, isn't
technically needed these days.  I guess you could have multiple copies of the
same major version installed, but the symlinks will only point at one of them,
so it's not clear to me how that would be useful beyond satisfying any hoarder
tendencies you might have!

#### libtool versions

Of course, I've glossed over an important detail.  You might assume from the
previous section that the Semantic Version that you've chosen for your project
(Foo Library 2.1.0) will line up with the version number encoded in your shared
library file (`libfoo.so.2.1.0`).  Au contraire, my friend!

While there are many build systems out there these days for the C and C++
ecosystem, the dreaded [autotools][] are still where we get our most cherished
conventions.  The autotool responsible for shared libraries is called
[libtool][], and it has its own [versioning scheme][libtool versions], with
exactly the same goals as semver — but its rules are **just** different enough
to give you different version numbers for the same sequence of API changes.  And
it's the libtool version that determines what ends up in your shared library's
filename.

[autotools]: https://autotools.io/index.html
[libtool]: https://www.gnu.org/software/libtool/
[libtool versions]: https://www.gnu.org/software/libtool/manual/html_node/Versioning.html

Note that I *didn't* say that the libtool version **is** what shows up in the
filename — it **determines** it.  That's right!  Not only is the libtool
versioning scheme *different*, it's also the *input* to a process that is what
*actually* determines what shows up in your filenames and in the linker commands
embedded in your shared libraries.  Fun stuff.

So what are these rules?  The libtool documentation goes into [detail][libtool
rules], but here's a summary:

[libtool rules]: https://www.gnu.org/software/libtool/manual/html_node/Updating-version-info.html

1. You assign your shared library a **version info**.  Like a traditional
   version number, this consists of three numbers, but they're separated by
   colons instead of periods, and they're out of order!

   The three numbers are **current**, **revision**, and **age**.  Current sorta
   lines up with semver's major version; age sorta lines up with the minor
   version; and revision sorta lines up with the patch level.

2. When you release a new version of your library, you bump the version info
   according to libtool's [rules][libtool rules]:

   - Backwards-incompatible change: bump current, set revision and age to 0.
   - Backwards-compatible change: bump current *and* age, set revision to 0.
   - No API change: bump revision.

3. You will endeavour to pass this version info to the `libtool` using the
   `-version-info` flag.  If you're using the autotools, you'll do this using an
   [`LDFLAGS` line][jansson ldflags].

4. libtool will [perform its magic][libtool munging] to transform this version
   info to a shared library filename, which basically involves undoing some of
   the math you performed in step 2.  The result will still not be the same that
   you'd get by following the semver rules instead.

[jansson ldflags]: https://github.com/akheron/jansson/blob/b23201bb1a566d7e4ea84b76b3dcf2efcc025dac/src/Makefile.am#L27
[libtool munging]: http://git.savannah.gnu.org/cgit/libtool.git/tree/build-aux/ltmain.in?id=722b6af0fad19b3d9f21924ae5aa6dfae5957378#n7042

#### CMake

What if you're not using the autotools?  I'm not going to go into every build
system that's out there, but [CMake][] is pretty common, and if you're using
something else, I trust you to wing it.

[CMake]: https://cmake.org/

CMake doesn't implement the [libtool magic][libtool munging] that transforms a
version info into the shared library filename.  Instead, it gives you nice
precise control over exactly what goes into your shared library's filename and
linker commands, using the [`set_target_properties` command][cmake rules] to set
the library's `VERSION` and `SOVERSION` properties.

[cmake rules]: https://cmake.org/cmake/help/v3.0/command/set_target_properties.html

If you want to follow libtool's scheme to the letter, you have two options:

1. Maintain an actual libtool-style version info for your library, like I do
   [here][libcork version-info] for my [libcork][] library, and reimplement
   libtool's magic using some equivalent [CMake trickery][cmake munging].

2. Cry, and do all of the bookkeeping yourself, by hand.

[libcork]: https://github.com/dcreager/libcork/
[libcork version-info]: https://github.com/dcreager/libcork/blob/660a112bdc65c5bea4f3d1f19e4e317a5e04775f/src/CMakeLists.txt#L41
[cmake munging]: https://github.com/dcreager/libcork/blob/660a112bdc65c5bea4f3d1f19e4e317a5e04775f/cmake/FindCTargets.cmake#L58

All of this is confusing and overly complicated, yes, but that's the world we
live in.

{::comment}

#### Multiple libraries in a project

Last wrinkle, I promise!

What if you have a single software *project* that produces more than one
*library*?  [Nettle][nettle] is a good example; it implements some cryptographic
primitives, some of which depend on another library (GMP), and some of which
don't.  Nettle publishes these in separate shared libraries (libnettle and
libhogweed), so that you aren't obligated to take a dependency on GMP (which is
rather large) if you don't need it.

[nettle]: https://www.lysator.liu.se/~nisse/nettle/

If I look at what's currently installed on my Linux laptop, we can see that the
version numbers of the two libraries aren't in sync:

    $ ls /usr/lib/lib{nettle,hogweed}*
    /usr/lib/libhogweed.so
    /usr/lib/libhogweed.so.4
    /usr/lib/libhogweed.so.4.3
    /usr/lib/libnettle.so
    /usr/lib/libnettle.so.6
    /usr/lib/libnettle.so.6.3

Nettle is also not included the patch level in their filenames, but let's
ignore that for now.  The important point is that you should track version
numbers **separately for each library** as well as **for the project as a
whole**.  It's entirely likely that none of these version numbers will line up
with each other, and that's okay!  This does come at the cost of a small amount
of extra bookkeeping, but it gives the build tools and the runtime linker the
most flexibility to only require upgrades when they're actually needed.


[abseil live at head]: https://abseil.io/about/philosophy#we-recommend-that-you-choose-to-live-at-head
[abseil rules]: https://abseil.io/about/compatibility
[libtool rules]: https://www.gnu.org/software/libtool/manual/html_node/Updating-version-info.html#Updating-version-info
[automake libtool]: https://www.gnu.org/software/libtool/manual/html_node/Using-Automake.html

#### Tying it all together

As a concrete example, let's say you're building a project (named "Project
Awesome", natch), and your development timeline looks something like the
following:

- You create a library called libawesome that contains some cool new code, and
  create an initial release to publish your code to the world.

  Project Awesome **1.0.0**,
  libawesome.so.**0.0.0** (0:0:0)

- Oh no!  There were some hideous bugs in your code!  Luckily some kind soul on
  Twitter (wait what?) politely pointed them out.  The fixes were pretty easy,
  and you didn't have to change anything in your API, so you commit them and
  release a new version.

  Project Awesome **1.0.1**,
  libawesome.so.**0.0.1** (0:1:0)

- You now have literally dozens of people using your code, and several of them
  have suggested the same cool new feature.  You find some time to implement it,
  adding it to your public API.

  Project Awesome **1.1.0**,
  libawesome.so.**0.1.0** (1:0:1)

- You notice that several people have implemented some wrapper functions that
  make it slightly easier to use your library.  You like what they've done, and
  you want to incorporate it (and them) into the official project!  Not everyone
  needs this new functionality, though, so you decide to stash it away in a
  separate "extras" shared library.

  Project Awesome **1.2.0**,
  libawesome.so.**0.1.0** (1:0:1),
  libawesome-extras.so.**0.0.0** (0:0:0)

- Crap!  After months of real-world use, you realize that the initial shape of
  your public API is completely wrong!  It takes a few weeks, but you hammer out
  something that's much cleaner and supports several use cases that you couldn't
  handle before.  This has completely changed your main public API, of course,
  but you notice that the API of your helper library didn't change at all!  Each
  of the wrapper functions can be re-written to use the new API.

  Project Awesome **2.0.0**,
  libawesome.so.**2.0.0** (2:0:0),
  libawesome-extras.so.**0.0.1** (0:1:0)
{:/comment}
