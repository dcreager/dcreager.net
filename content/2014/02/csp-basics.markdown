---
kind: article
created_at: 2014-02-28
layout: /post.html
title: "CSP: The basics"
tags: [csp]
---

{: .big-def}
*tl;dr*{: .tldr} CSP is a *formal method* that lets you describe and reason
about the behavior of *concurrent systems*.  CSP is *composable*; you write
simple *processes*, and then use special *operators* to combine them together
into larger, more complex processes.  A process is a summary of some system; it
uses *events* to describe how that system works, and to *synchronously
communicate* with other processes.  You can compare two processes using a
*refinement check*; this lets us check, for instance, whether a real-world
system satisfies some important safety or liveness property.  CSP has good *tool
support*, which lets us perform these refinement checks quickly and
automatically.

Well that was easy, wasn't it?  You can boil just about anything down to a
single paragraph.  Let's look at each of those key points in more detail.


<h5>Formal method</h5>

CSP is a *formal method*, which means that the language has a very precise
mathematical definition.  This is what allows us to *reason* about a system that
has been described in CSP: the language defines several different *semantics*,
which are mathematical rules that describe how a CSP model behaves.  Those rules
let us verify whether the model satisfies some property that you care about.  If
your model accurately describes your real-world system, then any conclusions
that we can draw about your model also apply to your system.

Of course, it's up to you to create a CSP model that accurately describes your
real-world system.  For a large and complex system, it might not be immediately
obvious how to describe the system in CSP.  That said, the language's operators
line up pretty well with the techniques we use to implement real-world systems,
so once you've learned the language, it's usually straightforward to construct a
CSP model for the system you want to analyze.


<h5>Concurrent systems</h5>

CSP is also a *process calculus*, which is a particular kind of formal method
that focuses on *concurrent* systems.  A concurrent system is one that is
divided into several smaller parts or subsystems.  For the most part, each of
these subsystems operate independently; however, they can also collaborate by
sending messages to each other.

This modeling approach is very similar to how we actually design complex
systems: we create small, highly focused subsystems that are responsible for
portions of the overall behavior, and we keep the interaction points small so
that we can understand how the parts work together to implement the whole.

Note that concurrency is a design technique, which you think about while
architecting and implementing your system.  It's not the same as *parallelism*,
which is an optimization technique.  Concurrency lets you simplify your design
by dividing up responsibilities among different subsystems.  Parallelism only
comes into play when your system is up and running: you try to do multiple
things simultaneously (or "in parallel") so that you can finish the overall job
faster.  The two are definitely related — you can't execute two things in
parallel unless they were already concurrent — but they're not the same thing.
Since CSP is a modeling language, it's focused on concurrency, not parallelism.


<h5>Processes and events</h5>

In CSP, each system and subsystem is modeled by a *process*, which describes a
particular pattern of *events* that can occur.  These events provide a summary
of the internal behavior of the process.  Processes are very similar to state
machines; you can think of events as evidence of a process moving from one state
to another.  Events usually correspond to a real-world condition that you want
to detect or action that you want to control.  When modeling a system using CSP,
one of your most important steps will be deciding which events to use, and what
those events represent.


<h5>Composition operators</h5>

CSP is *composable* — once you've defined the CSP processes for your small
subsystems, you use one of CSP's composition operators to join those subsystems
together into a larger, more complex process.  There are several composition
operators to choose from, depending on how each of the subsystems should
interact with each other.

This approach makes it much easier to model and reason about a complex system:
each subsystem is small enough that you can understand how it works, and the
composition operators provide well-defined rules for how those smaller pieces
work together to produce some more complex behavior.


<h5>Synchronous communication</h5>

CSP events not only describe the behavior of individual processes; they also
allow multiple processes to communicate with each other.  If two processes both
include the same event, then when you compose those processes together, that
event becomes a communications channel that links the two processes.

For example, a CSP process that models a vending machine might use an event
called `coin` to detect when someone has inserted a coin into the machine.  A
second CSP process that models the person using the machine would use this same
`coin` event to represent the act of inserting a coin into the machine.  When we
compose these two processes together, we get a new process that represents the
interaction between the person and the machine.

One of the most important features of CSP is that **all communication is
synchronous**.  When two processes communicate via an event, that event will
either occur simultaneously in both processes, or it won't occur in either.  In
our example, this means that the machine won't detect a coin if the person never
inserted one; and vice versa, that the person can't insert a coin without the
machine detecting it.


<h5>Refinement checks</h5>

The main reason for using a formal method to model a real-world system is so
that you can verify some properties of that system.  For instance, with our
vending machine, we might want to make sure that "if you insert a coin, you
either get a drink or you get your coin back."

In CSP, you do this using a *refinement check*.  "Refinement" is the
mathematical term for the particular way that we compare the behavior of two
processes.  If you have a CSP process describing your system, and another
describing the property you want to verify, then a refinement tells you whether
the system satisfies that property.  More importantly, the refinement doesn't
just tell you "yes" or "no"; it also provides a *proof* or a *counterexample*
automatically.

This is true of all formal methods — the entire reason for using a formal method
is so you can perform these rigorous checks using well-defined rules, which give
you proofs about the properties that you care about.  That said, CSP differs
from most other formal methods in that you use the same language to describe
your system and the properties that you want your system to have.  In other
formal methods, you use one language to describe your system, and a completely
different language to describe the properties that you want to check.  In CSP,
you use processes and events to describe both.


<h5>Tool support</h5>

Refinement checks are nice, but they'd be rather pointless if we had to do them
by hand.  Our real-world systems are usually too large and complex for that to
be practical.  If you have to write hundreds of pages of detailed mathematical
proof to convince yourself that your system satsifies a safety property, then
the benefit of having that mathematical proof isn't worth the cost of producing
it.

That means that to be useful, a formal method needs to have good *tool support*.
For CSP, the tool of choice is [FDR](https://www.cs.ox.ac.uk/projects/fdr/).
You provide machine-readable descriptions of your CSP processes, and ask it
perform refinement checks for you.  It also provides a nice user interface for
viewing counterexamples, to help you identify which parts of your system violate
the properties that you're checking.
