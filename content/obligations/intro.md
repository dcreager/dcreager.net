+++
date = 2017-08-10
title = "Introduction"
description = "in which our options are limited"
#tags: [obligations]
[extra]
series = "Obligations"
+++

Programming languages have a lot of interesting ways to think about what you're
allowed to do, and not allowed to do.  In a statically typed language, your
compiler will yell at you if you try to do something that's not allowed, and a
lot of interesting PL research involves teaching your language and compiler to
disallow things in increasingly sophisticated ways.

This is not limited to statically typed languages.  In a dynamically typed
language, your interpreter is just as invested in disallowing certain behaviors.
It's just that the enforcement happens while your program is running!

Someone with a theoretical background might quibble with my choice of the word
"disallow".  They might prefer that I say that certain operations in your
programming language are "invalid".  In the world of mathematics, it's usually
not the case that some little gremlin is actively preventing these operations
from occurring; instead, it's that the invalid operations don't even exist as
possibilities!

But regardless what language you're programming in, and regardless of whether
you think in terms of **permissions** or **possibilities**, the act of
programming is to consider all of these possible operations, and choose which
ones you want to use to accomplish your goal.  At any given time, there are many
possible operations, and you have the control to decide which ones to use!

In this series of posts I want to look at the flip side of this, and talk about
**requirements** or **obligations**.  What mechanisms do programming languages
have to *force* you to perform some operation, or to ensure that some operation
is performed even if you don't actively choose to do it yourself?

At the end of this series, I hope to convince you that obligations are just as
important as possibilities, and that there are interesting complex behaviors
that are easier to think about and work with if we use programming languages
that can talk about obligations directly.

<!-- more -->

{% comment() %}
not permissions "you can only close a file after you opened it"

obligations: "once you've opened a file you MUST close it"

simple initial example: guaranteeing that you close a file exactly once after
you open it

sad face (C)

GC objects + finalizers (Python)

not ideal, since you aren't guaranteed when that will run


Scope-based solutions (Python context managers, RAII/destructors in C++)

enforcing that something will happen by piggy-backing on another language
feature


Works even for things not technically "scope-shaped"

You can pass that file around in C++.  What if they save it somewhere?  You'll
have a dangling memory reference!

You can pass that file object around in Python.  What if they save it somewhere?
The file will be closed at the right time, but you'll have a dangling reference.


Move semantics in C++ / Rust ownership types to the rescue

You can't save the reference, you can only borrow it
{% end %}
