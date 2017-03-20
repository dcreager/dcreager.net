---
kind: article
created_at: 2010-01-08
layout: /post.html
title: Default “scons -c” targets
tags: [scons]
---

As I mentioned in a [previous
post](/2009/12/18/make-distclean-in-scons/), the automatic “clean”
target provided by SCons (`scons -c`) is very useful for cleaning out
build files, without requiring much in the way of configuration.
Anything that SCons generates when you run `scons` will be
automatically cleaned when you run `scons -c`.

While useful, I'd like more control over the behavior of `scons -c`.
Specifically, being a good TDD junkie, I have several test cases that
I can run using `scons test`:

<% highlight :python do %>
build_test = env.Program( ... )
env.Alias("build-tests", build_test)

run_test = env.Alias("test", [build_test],
                     ["@%s" % build_test[0].abspath])
env.AlwaysBuild(run_test)
<% end %>

By setting it up this way, the test programs aren't built by default:
you have to explicitly run `scons build-tests` (if you want to build
the tests but not run them) or `scons test` (if you want to build and
run them).  Moreover, because of SCons's dependency tracking, I can
just use `scons test` as my usual build command during my
Edit-Test-Debug loop.  SCons will automatically rebuild any changed
source files before running the tests.

All of this is great.  So what's the problem?  As I mentioned above,
`scons -c` only cleans the build files that are created by `scons` —
and since I've explicitly set things up so that tests aren't _built_
by default, they'll also not be _cleaned_ by default.  This means that
to fully clean out my build targets, I have to run two commands:

    $ scons -c
    $ scons -c build-tests

Not ideal!  I'd prefer if `scons -c` cleaned everything, just like
`make clean` would in the Automake world.

## The solution

So how to fix this?  First we need to understand how SCons decides
what to clean when you run `scons -c`.  The answer is “exactly what's
built by `scons`”.  And how does SCons decide what to build when you
run `scons`?  That's what the `Default` command is for.

For instance, I could add

<% highlight :python do %>
env.Default("build-tests")
<% end %>

to my _SConstruct_ file.  This would cause all of my tests to be built
by default, and by extension, to have them all cleaned by default, as
well.

This is close, since `scons -c` now does what we want, but this means
that `scons` is now building more than we would like.  What we need is
a way to have a different list of default targets depending on whether
we're building or cleaning.  Luckily, the `GetOption` function gives
us exactly that:

<% highlight :python do %>
if GetOption("clean"):
    env.Default("build-tests")
<% end %>

With this in our _SConstruct_ file, the tests will be considered a
default target when we're cleaning, but not when we're building.  So
now we have what we want: `scons` builds just the code, `scons test`
builds and runs the tests, and `scons -c` cleans it all.
