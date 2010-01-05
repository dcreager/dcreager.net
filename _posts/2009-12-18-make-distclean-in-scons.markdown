---
layout: post
title: Simulating “make distclean” in SCons
tags: scons
---

SCons provides an automatic “clean” target out of the box — just run
`scons -c`, and SCons will delete all of the objects that it knows how
to build.  This is a very useful feature; however, there are two main
missing features that I want to add to my build scripts.  First, I
want to be able to delete all of the temporary files SCons uses, such
as its configuration cache and any files I use to store variable
values.  These aren't included in the default list of the files to
clean up.  Second, I want more control over which items are deleted by
default, when you specify `scons -c` without any targets.  I'll
describe my solution to the first problem in this post.  I'll write up
the second problem in another post.

## Deleting SCons's temporary files

This feature is akin to the `make distclean` target that Automake puts
into the Makefiles that it generates.  This differs from `make clean`;
`make clean` is intended to delete all of the build products, but
leave behind the results of the `configure` step, whereas `make
distclean` is supposed delete _everything_, returning the source tree
to the same state as when you'd just unpacked the tarball.

The `scons -c` command is analogous to `make clean`, and requires no
setup.  It will automatically delete any of the build products that
are created by running `scons`.  There are several cache files that
SCons creates, however, and it would be nice to have an equivalent to
`make distclean`.  This is especially useful when developing a new
`Configuration` check, for instance — if you make a change to the
test, you want to be able to (easily) force SCons to ignore any cached
results, and try all of the tests again.

This is actually fairly easy to set up, assuming you know the list of
temporary files that SCons will create.  You can add the following
rule to your top-level _SConstruct_ file:

{% highlight python linenos %}
env.Clean("distclean",
          [
           ".sconsign.dblite",
           ".sconf_temp",
           "config.log",
          ])
{% endhighlight %}

As far as I can tell, these three files are always created by SCons.
To delete these files, you simply run `scons -c distclean`.  Because
we've defined the target using `Clean`, it will only be run when you
pass in the `-c` option to `scons`.

Since we're putting together the file list manually, you should make
sure to add any additional cache files that your SCons scripts use.
For instance, I'm using some `Variable` options, which I store into a
file called _.scons.vars_.  (This means that the user doesn't have to
type them in with every invocation of `scons`.)  By using these
variables, I have to add another entry to the `distclean` target:

{% highlight python linenos %}
vars = Variables('.scons.vars', ARGUMENTS)
# ...define a bunch of variables
vars.Update(env)
vars.Save(".scons.vars", env)

env.Clean("distclean", [".scons.vars"])
{% endhighlight %}

Note how, just like with any SCons target, I can define the
`distclean` target multiple times.  SCons will take care of merging
them into a single action, deleting all of the specified files when
you run `scons -c distclean`.
