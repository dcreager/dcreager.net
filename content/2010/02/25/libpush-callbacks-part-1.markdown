---
kind: article
created_at: 2010-02-25
layout: post
title: Parser callbacks in libpush, Part 1 — Streams
tags: [libpush, c]
---

This post is the first in a series that describes the
`push_callback_t` type in the
[libpush](http://github.com/dcreager/libpush/) library.  In these
posts, we'll walk through a couple of possible ways to implement
callbacks under the covers.  At each stage, we'll encounter problems
with the current design.  Fixing these problems should lead closer us
to the actual implementation in libpush, and along the way, we'll gain
a good understanding of how our design decisions affect the
performance and usability of the library.

The `push_callback_t` type is used to define _parser callbacks_, which
are the basic unit of parsing in libpush.  Callbacks are pretty
simple: they take in an _input value_, read some data from the _input
stream_, and produce an _output value_.  (The fact that callbacks take
in an input value, in addition to reading from the input stream, is
what makes them [_arrows_](http://www.haskell.org/arrows/) instead of
[_monads_](http://en.wikipedia.org/wiki/Monad_%28functional_programming%29)
— but that's a story for a later post).

## First attempt: Callbacks as functions

Now, with this simple structure, we might try to implement callbacks
as regular C functions.  For instance, we could use something like the
following to read in a single 32-bit integer:

<% highlight :c do %>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

bool
parse_uint32(void *input, uint32_t *output, FILE *stream)
{
    size_t  num_read;
    num_read = fread(output, sizeof(uint32_t), 1, stream);
    return (num_read == 1);
}
<% end %>

This callback ignores its input value, reads in four bytes from the
input stream, and uses that to output a `uint32_t` value.  The return
value of the function is a boolean, indicating whether the parse was
successful or not.  This lets us handle _parse errors_ — for instance,
if there are only three bytes left in the stream, we can't read in a
full integer.  We return `false` to indicate this error condition.

We've ignored some details here that aren't important for this example
— for instance, we don't worry about the endianness of the integer,
nor do we worry about how the space for the output result is
allocated.  We just assume that someone will pass in a pointer to a
`uint32_t` variable, and our callback function will store its output
value there.

## Drawbacks

This approach works fine for simple cases, but unfortunately has two
drawbacks.  First, we're limited to parsing from `FILE` streams.  Any
real input source will probably be available as a stream, so this
might not seem like a huge problem — though it does rule out parsing
from a memory buffer, unless you use a non-portable function like
`fmemopen`.

The second, more important, problem is that the parser callback has
full control over when and how much to read from the stream.  In this
example, we try to read in the full four bytes for the `uint32_t`
output value.  However, there might not be four bytes available in the
stream.  If this is because we're at the end of a file, then we should
treat this as a parse error.  If we're reading from a network socket,
though, another chunk of data might arrive if we wait for a bit.

We could add logic to the callback to read from the stream repeatedly
until we got enough data, but then we'll start _blocking_ — so that we
can distinguish between "there's no more data here _yet_" from
"there's no more data coming _at all_".

All of this is bad news.  First of all, this extra I/O logic is
starting to get rather big, and we don't want each and every callback
to have to include it.  And second, we don't want the rest of our
program to be held hostage by the callback — it should be up to our
I/O code to decide whether it's okay to block waiting for more input,
or whether to whip up a nice `select` loop of some kind to read things
more efficiently.

In the next post, we'll describe _iteratees_, which give us this
capability.
