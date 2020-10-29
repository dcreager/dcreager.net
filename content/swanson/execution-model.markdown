---
kind: article
created_at: 2020-10-29
layout: /post.html
series: "Swanson"
title: "Execution model"
description: >-
  Linear continuations will keep everyone in line
tags: [swanson]
prev: /swanson/intro.markdown
draft: true
---

This post summarizes the computation model that the Swanson framework builds on.

### Hosts

A Swanson program has to execute somewhere.  This place is called the
**_host_**.  The Swanson execution model takes very great care to make _no
assumptions whatsoever_ about the host.  The only real assumption is that the
host somehow knows how to execute Swanson code!  That might mean that it
interprets it directly, or that there's some compiled representation that can be
executed, and which brings about the same effect as the Swanson program that was
compiled.  As long as the host can faithfully implement the execution model
described here, we're golden!

What kinds of hosts are we talking about?  A "native" host would compile a
Swanson program down to machine code, just like a C, Rust, Go, or Haskell
compiler would, and then execute it directly.  A "JVM" host would compile a
Swanson program down to JVM bytecode, so that it could more easily interact with
code written in Java.  An "embedded interpreter" would implement the execution
rules below as simply as possible, with no compilation or optimization.  That
would make the interpreter very simple, and make it easy to include it in
another application written in some other language.

<div class="aside" markdown=1>
A particular Swanson program is free to assume more about the host than just
this.  It would typically do so by depending on code that isn't available on all
hosts.  (We talk about how you depend on other code below.) And if some of the
code your program depends on isn't available on some host, you can't excute your
program on that host!
</div>

### Values and environments

A Swanson program operates on an _environment_, which is a collection of
_values_, each with a separate name.  (A good intuition is that it's the set of
parameters and variables that are in scope at the current point of execution.)

There are only three kinds of values: **_atom_**, **_literals_**, and
**_invokables_**.

#### Atoms

Atoms are basic entities that exist and can be compared for equality.  They are
primarily used at compile-time, and let us determine or assert that two things
(types, values, anything really) are "the same".

#### Literals

Literals represent immutable binary content.  These are the only constants
available to a Swanson program.

<div class="aside" markdown=1>
Swanson programs use literals to instantiate "real" constants of other types.
To create a numeric constant, for instance, you'd typically create a literal
whose (binary) content is the ASCII rendering of the desired numeric value, and
then use an invokable to parse that binary literal into the internal
representation of the numeric value.  Because we have the ability to explicitly
mark that code should execute at compile time, there's typically no runtime cost
incurred.
</div>

#### Invokables

Everything else — runtime values, types, classes, functions — is represented by
an invokable.  An invokable consists of one or more **_branches_**, each with a
**_name_** and a **_body_**, which is some executable Swanson code.  An
invokable will typically also contain some hidden state.  (Note that does not
necessarily mean _mutable_ state!)

The main thing that you do to an invokable (as you might guess from its name) is
that you **_invoke_** one of its branches.  The current environment is "passed
into" the invokable as its input, and execution proceeds with the body of the
selected branch.

<div class="aside" markdown=1>
Note that control **_never returns_** from an invokable!  Invokables are like
continuations.  Invoking is one-way.  Everything is a tail call.

That doesn't mean that you can't model something like "returning from a
function"!  You'll have an invokable representing the function, and its single
branch contains the body of the function.  To "call" the function, you invoke
its invokable.  But instead of implicitly keeping track of a stack of function
calls, you must explicitly pass in a **_continuation_** as one of the inputs.
This continuation represents the call stack that gets constructed before calling
the function.  Inside the function's invokable, it will invoke the continuation
anywhere that control is `return`ed to the caller.

This also means that invokables don't really have "output values".  In Swanson,
an "output" is really just the input that's passed into the continuation when
it's invoked.  Control flow always moves forward!
</div>

For typical values that you'd encounter in other languages (integers, strings,
booleans, instances of classes, instances of algebraic data types), the
invokable encapsulates together the value itself, along with all of the
operations that you can perform on that value (each represented by a branch).

In languages that have a concept of pointers or references, the "thing pointed
to" and the "pointer" or "reference" are both represented by (different)
invokables.  The "thing pointed to" will have a "create reference" operation as
one of its branches (corresponding with the `&` operator from C or Rust).  The
"pointer" will have a "dereference" operation (corresponding with the `*`
operator).

Invokables are also used to represent the basic blocks in your program.  The
internal state represents any values that are closed over — local variables on
the stack in a typical imperative language, or values captured in a closure in a
functional language.  This kind of invokable would typically have a single
branch, representing the normal control flow of your program.  But in a language
that supports exceptions or JavaScript-style promises, there would be two — one
branch to handle the "success" case, and one to handle the "error" case.

### Primitives

Most invokables are implemented as Swanson code.  (Or more accurately, they're
implemented in S₀, which is the Swanson "assembly language" that we'll learn
about later, or in a language that can be compiled or translated into S₀.)

However, you'll probably have noticed that, unlike other language runtimes,
there aren't any "real" types mentioned above.  No specific integer types, no
strings, no booleans, nothing like that.  In Swanson, none of that is
predefined.  Instead, those are provided by **_primitives_**.  A primitive is a
special type of invokable that is provided by the host, instead of being
implemented in Swanson itself.

<div class="aside" markdown=1>
The Swanson execution model does not predefine any of these primitives.  It
doesn't say which ones are available, what their names are, or how they work.
That might seem to be a hard programming framework to program against!  However,
there _is_ a "standard library" of primitives, which we'll describe later.  Most
hosts provide most of the primitives in the standard library.  So for all
practical purposes, the set of primitives _is_ predefined — just not here at the
lowest level of the framework.
</div>

### Linearity

Environments and values are **_linear_**, which means that they must be used
**_exactly once_**.  That means two things:  First, you cannot use a value more
than once.  Invoking an invokable, or using an atom or literal, consumes the
value.  Second, you must use every value.  When a Swanson program finishes, the
environment must be empty.

<div class="aside" markdown=1>
Primitives are "magical" in the sense that they let you side-step linearity.
For example, a host could provide a primitive for duplicating values that are
safe to duplicate, and wouldn't provide them for other values.  Since there's no
way to duplicate values purely in Swanson, that means we can be sure that all
Swanson code follows the rules provided by the host environment.

This will be a pattern that we encounter _all the time_: linearity exists to
provide a base level of safety, and lets us reason about when things _can_ and
_must_ occur, but there will be "escape hatches" provided by the host that let
us break out of those shackles in carefully controlled, well-scoped ways.
</div>

If an invokable is consumed when you invoke it, that works really well to model
something like a destructor: calling a value's destructor really should "remove"
it from the caller's context, so that the caller can't accidentually try to use
the now-freed value.  In this case, linearity provides a really nice answer to
the "use after free" problem.

But how then do we model functions or methods that _don't_ consume their
receiver?  In Swanson, the invokable must "return back" the value to the caller,
making it available for additional, later method calls.  (You will usually have
to construct a new invokable to represent the returned-back value, since the
original was just consumed; we'll see examples in later sections of how to do
that.)

<div class="aside" markdown=1>
This pattern — where invokables explicitly return back values that can be used
again in other ways — means that the set of operations you can perform on a
value can change over time!  There's no obligation that the value that's
returned back has the same set of available operations (that is, the same set of
branches in the invokable) as the one that was passed in.

This lets us model things like Rust's lifetimes, which say that you can't drop a
value when there are still open references to it somewhere.  In Swanson, the
Rust version of the "create reference" operation would _both_ return the
reference as a new invokable, and also **_change_** the set of available
operations on the original value to prevent you from calling its "drop"
operation.  Dropping the reference would then reenable the "drop" operation on
the original value.
</div>

### Modular programming

The last piece of our execution model lets you "program in the large", breaking
up a large program into smaller pieces that you can write, test, and compile
separately.  Most programming languages provide some kind of facility for this:
packages in Java, crates in Rust, translation units in C, packages in Go.

In Swanson, the primary unit of modularity is the **_unit_**.  Each unit has a
**_name_**, and can be **_loaded_**, producing a Swanson value.  The host is
responsible for determining how to locate a unit with a particular name, and for
loading it.

<div class="aside" markdown=1>
Loading a unit typically produces an invokable, where each branch defines one of
the entities "exported" by the unit.  But this isn't a hard requirement — it's
perfectly fine for a unit to produce an atom or literal, though that wouldn't
really be very useful in practice.
</div>

Units can have **_dependencies_**, which are the other named units that this
unit needs as part of its loading or execution.  The host will take care of
loading those dependencies first (ensuring that there are no circular
dependencies), and providing them to the unit.

The **_entry point_** of a Swanson program is a unit that defines the main body
of the program.  (This corresponds with the `main` function of a C, Rust, Go, or
Java program.)  An entry point unit, when loaded, must produce an invokable with
a single branch.  The host will invoke this branch to execute the program.  Like
any other unit, this entry point can (and almost certainly will) have
dependencies on other units, which it will use to do useful work.

<div class="aside" markdown=1>
A unit can be loaded multiple times — once for each time it's used as a
dependency.  Unlike other frameworks, we don't "cache" the unit's value and
return it multiple times, since that would violate linearity!  When a unit is
loaded, the value that's produced can only be used once.  If the unit is used as
a dependency multiple times, each of those dependents must end up with distinct
values to obey linearity.  That, in turn, means that we have to load the unit
multiple times, producing multiple values, each of which can be used (exactly
once) by the dependents.

In later posts we'll see some common patterns that let a unit "share state"
across all of its uses.  As you might guess, that requires taking advantages of
some primitives provided by the host, since basic linear Swanson values don't
give you this ability.
</div>
