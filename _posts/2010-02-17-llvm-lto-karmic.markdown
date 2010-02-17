---
layout: post
title: Using LLVM's link-time optimization on Ubuntu Karmic
tags: [ubuntu, llvm]
---

While playing around with
[libpush](http://github.com/dcreager/libpush) on my MacBook, I was
pleasantly surprised to see a huge performance increase when I used
the link-time optimization (LTO) feature of the LLVM GCC front end.
(It's really quite nifty; the new [Homebrew package
manager](http://github.com/mxcl/homebrew) uses it by default when
compiling packages.)  On MacOS, using LTO is as simple as using
`llvm-gcc` as your C compiler (or `llvm-g++` if you're compiling C++),
and passing in `-O4` as your optimization flag.  I use SCons as my
builder, so this turns into:

    $ scons CC=llvm-gcc CCFLAGS=-O4

This will cause GCC to output LLVM bytecode into the _.o_ output
files, and to perform whole-program optimizations during each linking
phase.  I was able to see a big performance win simply from the linker
being able to inline in copies of small functions that live in “other”
compilation units.

## Good news and bad news

Intrigued by the results, I wanted to try the same thing on my Linux
boxes, which are running Ubuntu Karmic.  On the Mac, Apple has made
sure to include support for LLVM in all of the standard Xcode build
tools.  On Linux, you don't get this by default right now — though GCC
is implementing their own LTO project, which is starting to bear
fruit.  Part of this is a new “`gold`” linker, which supports a plugin
architecture.  How is this useful to us?  Well, LLVM already has a
[plugin](http://llvm.org/docs/GoldPlugin.html) for the new linker, so
with everything installed correctly, getting LTO through LLVM on Linux
can be just as simple as it was on the Mac.

Unfortunately, these new tools have only partially made it into the
Ubuntu package tree.  You can get the new `gold` linker by installing
the `binutils-gold` package, and you can get most of the LLVM pieces
by installing the `llvm` and `llvm-gcc-4.2` packages.  Unfortunately,
this doesn't include the LLVM `gold` plugin or the new `clang` C/C++
compiler front-end.  Things look promising for these features being in
the new Lucid packages — which could even lead to a Karmic backport —
but for now, if we want the `gold` plugin, we have to compile
ourselves.


## Getting the prerequisites

As mentioned on the LLVM [linker plugin
page](http://llvm.org/docs/GoldPlugin.html), you need to have the
`binutils` source lying around somewhere if you want to compile the
plugin, since the LLVM source needs to read in `binutils`'s
_plugin-api.h_ file.  The easiest way for us to get the `binutils`
source is using APT:

    $ mkdir -p $HOME/deb
    $ cd $HOME/deb
    $ apt-get source binutils

This will place an unpacked copy of the `binutils` source into
_$HOME/deb/binutils-2.20_ for you.

We can also go ahead and install the `gold` linker:

    $ sudo apt-get install binutils-gold

You'll also need to make sure you've got the basic compilation tools
installed (though if you're at the point where you're trying to play
around with LTO, I've got to assume you've already taken care of
this...):

    $ sudo apt-get install build-essential

Finally, my main Linux box is 64-bit, so I need to install multilib
support before we can compile the LLVM GCC front end:

    $ sudo apt-get install gcc-multilib


## Compiling LLVM

With all of the prerequisites installed, we can download and unpack
LLVM:

    $ mkdir -p $HOME/tmp
    $ cd $HOME/tmp
    $ wget http://llvm.org/releases/2.6/llvm-2.6.tar.gz
    $ wget http://llvm.org/releases/2.6/clang-2.6.tar.gz

    $ tar xzvf llvm-2.6.tar.gz
    $ tar xzvf clang-2.6.tar.gz

`clang` is distributed as a separate download, but we actually want to
place it into the main LLVM directory; the LLVM build scripts will
find it and build it automatically:

    $ mv clang-2.6 llvm-2.6/tools/clang

At this point we can do the usual compilation steps:

    $ cd llvm-2.6
    $ ./configure \
        --with-binutils-include=$HOME/deb/binutils-2.20/include \
        --enable-optimized \
        --prefix=/usr/local
    $ make
    $ sudo make install
    $ sudo ldconfig

Notice how we're going to install everything into _/usr/local_, so as
not to step on the toes of the package manager.  This means we have to
run `ldconfig` so that the system linker knows about the new libraries
we just put in _/usr/local/lib_.


## Compiling LLVM-GCC

At this point, we have the `gold` linker installed, and have a copy of
LLVM that includes its `gold` plugin.  Ideally, we could start
compiling with `clang` and get LTO, but it doesn't seem like there's
currently a way to have `clang` pass in the necessary `--plugin`
option to the linker.  So, all we need now is the GCC front end.

As before, we start by downloading and unpacking:

    $ cd $HOME/tmp
    $ wget http://llvm.org/releases/2.6/llvm-gcc-4.2-2.6.source.tar.gz
    $ tar xzvf llvm-gcc-4.2-2.6.source.tar.gz

The _README.LLVM_ file in the source tree gives more detail on the
options you have available; for me, the following worked:

    $ mkdir -p $HOME/tmp/obj
    $ cd $HOME/tmp/obj
    $ ../llvm-gcc-4.2-2.6.source/configure \
        --prefix=/usr/local \
        --program-prefix=llvm- \
        --enable-llvm=$HOME/tmp/llvm-2.6 \
        --enable-languages=c,c++
    $ make
    $ sudo make install
    $ sudo ldconfig

The only interesting wrinkle is that we have to do an out-of-source
build — the object files will end up in the _$HOME/tmp/obj_ directory,
rather than being created directly in the unpacked source directory.

As this point we're nearly there; we have `llvm-gcc` installed, but
its `-use-gold-plugin` option won't work just yet.  If you look
closely at one sentence on the [LLVM plugin
page](http://llvm.org/docs/GoldPlugin.html), you'll see that the
option “looks for the `gold` plugin in the same directories as it
looks for `cc1`”.  The LLVM GCC package installed the `cc1` program
into the */usr/local/libexec/gcc/x86_64-unknown-linux-gnu/4.2.1*
directory.  (The *x86_64* will be different if you're on a different
architecture.)  However, the LLVM plugin is in _/usr/local/lib_.  If
you try to use the `-use-gold-plugin` parameter, you'll get the
following error message:

    $ llvm-gcc -use-gold-plugin \
        -o foo.o -c -O4 -g -Wall -Werror foo.c
    llvm-gcc: -use-gold-plugin, but libLLVMgold.so not found.

Not good.  The solution (which is admittedly a bit of a hack) is to
copy the plugin into the directory that `llvm-gcc` expects to find it
in:

    $ sudo cp /usr/local/lib/libLLVMgold.so \
        /usr/local/libexec/gcc/x86_64-unknown-linux-gnu/4.2.1


## Using your new toy

Now that we've got all of the pieces installed, you can create
libraries and executables that are optimized at link time.  The
“Quickstart” section at the end of the [LLVM plugin
page](http://llvm.org/docs/GoldPlugin.html) gives you the outline.  I
use SCons as my build tool, so I have to run the following:

    $ scons \
        CC="llvm-gcc -use-gold-plugin" \
        AR="ar --plugin libLLVMgold.so" \
        RANLIB=/bin/true \
        CCFLAGS=-O4

This is slightly more than what's needed on the Mac, but all in all,
not bad.  Enjoy!
