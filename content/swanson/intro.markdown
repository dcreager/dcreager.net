---
kind: article
created_at: 2020-10-28
layout: /post.html
series: "Swanson"
title: "Introduction"
description: >-
  Another series of posts kicked off with no end in sight
tags: [swanson]
next: /swanson/execution-model.markdown
draft: true
---

My [Twitter bio](https://twitter.com/dcreager) currently lists me as (among
other things) a "PL dilettante".  Which of course means that I've been hacking
around on my own programming language for a number of years!  (Seriously, this
has been a thing for a **_long_** time.  The earliest [Swanson-related
commit][old commit] I can find is from 2012.)  It's gone through a number of
iterations over that time, but I'm pretty happy with where it's at right now.
Not complete by any stretch of the imagination.  But I've chatted about it in
passing with several people at this point, and I figured I need to write down
the details somewhere.  My friend and colleague Rob can hold forth on these
kinds of topics in [epic][rob 1] [Twitter][rob 2] [threads][rob 3], but I need
something more longform.  So here we are!

[old commit]: https://github.com/swanson-lang/swanson-lang-old/commit/e8aac263d76a63c8667fca5eee079ce8af345c2d
[rob 1]: https://twitter.com/rob_rix/status/1320544724864872453
[rob 2]: https://twitter.com/rob_rix/status/1320488333127110659
[rob 3]: https://twitter.com/rob_rix/status/1320467572161236994

In this series of posts I'm going to describe Swanson, the programming language
framework that I've been noodling on.  It has two main components, one of which
is much more fleshed out than the other.

The first (not as fleshed out) part is an actual programming language, which
doesn't even really have a name yet!  It brings in some interesting notions
about programmable syntax and parsing in a way that (I think) makes it easy to
construct things like DSLs.

The second (more important) part is an execution model that can be used as an IR
of sorts for **_all_** languages.  That makes it similar in spirit to
WebAssembly, in that the goal is to be something that you would compile (or
translate or transpile or whatever) many other languages **_into_**, and not
something to be written directly.  The Great Unnamed Language, like any other
language, would be compiled into Swanson the execution framework.

Both the language and the execution model incorporate two big ideas.  The first
is that they build on Zig's notion of [`comptime`][zig comptime] — being able to
execute functions at compile-time.  Other languages have this capability, but I
find that Zig has done the most to internalize it as a first principle.  Can we
really [build every other useful language feature][zig claim] on top of
`comptime`?

[zig comptime]: https://ziglang.org/documentation/master/#comptime
[zig claim]: https://news.ycombinator.com/item?id=24292760
[zig dcreager]: https://twitter.com/dcreager/status/1299489029881294849

Personally, I think that [should be enough][zig dcreager]!  I think the other
thing that's missing is some form of ownership tracking, along the lines of
Rust's or C++'s move semantics.  For that, Swanson looks to **_linear
continuations_**.  The "continuation" part means that Swanson's basic blocks
don't implicitly pass control to anything else; you must explicitly pass in a
continuation, and the last instruction of a basic block must explicitly invoke a
continuation.  _All_ control flow is handled this way.  The "linear" part means
that you are **_obligated_** to consume every value and invoke every
continuation at some point.  (There are escape hatches that we'll get into
later.)  That lets you handle things like [RAII][] and [`defer` statements][] —
and, possibly less obviously, things like safe memory management, including Rust
lifetimes and borrow checking.

(Part of the purpose of these posts is to show my work on those claims, and see
whether or not they're actually true!)

[RAII]: https://en.wikipedia.org/wiki/Resource_acquisition_is_initialization
[`defer` statements]: https://golang.org/doc/effective_go.html#defer

So, let's jump in and see where this goes!
