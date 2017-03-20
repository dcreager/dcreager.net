---
kind: article
created_at: 2016-08-03
layout: /post.html
series: "Concurrency models in CSP"
title: "Preliminaries"
description: >-
  in which we try to get through all of the annoying boilerplate before getting
  to the good stuff
tags: [csp]
prev: /2016/07/csp-concurrency-intro.markdown
next: /2016/09/csp-read-atomic-internal.markdown
---

Before diving into the details of our first concurrency model, it will be
helpful to give an overview of how we're going to use CSP in these posts.

At the highest level, a CSP specification consists of a set of *processes*,
which use *events* to describe their behavior.  We use a process to represent
any vantage point that lets us talk about which events occur and in which order.
Each entity in our system is a process (which might be defined using
subprocesses to describe each of its components or parts).  We also use
processes to describe how two or more entities interact with each other, and to
describe any global properties of the entire system.  And lastly, because CSP
verification is based on "refinement", we also use processes to define the
properties that we hope our system exhibits.

<hr class="jump">

Each process has an *alphabet*, which is the set of events that it describes
(and constrains!).  The process itself is the description of when each of these
events can occur (and when they cannot!).  There are two ways to define a
process.  The first is as an *automaton* (called an "LTS" or "labeled transition
system" in the CSP literature).  This is just your standard state machine, where
nodes represent the current state of the process, and edges represent which
events are allowed in that state (and which state you move to if the event
occurs).  This kind of process description is useful because you can use
automated tools (like [FDR][]) to analyze them — for instance, proving that a
particular system (described by one CSP process) satisfies a particular property
(described by another).

<!--
It doesn't technically have to be a *finite* state machine, but if you want to
use a tool like FDR to analyze your specification, it will be.
-->

[FDR]: https://www.cs.ox.ac.uk/projects/fdr/

For humans, on the other hand, we define processes using *CSP operators*.  These
operators are more high-level, and more intuitive, than the lower-level LTS
representation, and they line up with common patterns that we use when
describing or implementing a complex system.  In particular, there are several
operators for combining subprocesses into larger processes, giving you full
control over how the events in each subprocess interact with each other.  There
is a well-defined way to translate each operator into an LTS, which is what
allows tools like FDR to analyze our human-readable process specifications.

[Just like the real world][cmeik], CSP does not impose any global ordering on
any of the events in our system.  We can only talk about event ordering in the
context of a single process.  If we want to talk about global event orderings,
we have to construct a process that includes information about how *all* of the
entities in our system participate in those events.

[cmeik]: https://christophermeiklejohn.com/lasp/erlang/2015/10/27/tendency.html

To dig into this paper, then, we're going to define a CSP process for each
concurrency model.  The original paper uses several *axioms* to collectively
define the behavior of each concurrency model; the more constrained models later
in the paper include more axioms.  We'll define each of those axioms as a
subprocess, and use CSP operators to combine them together into a single process
describing the concurrency model as a whole.

We'll also create CSP processes for each of the reference implementations.  The
paper then defines and uses "operational refinement" to show that each reference
implementation really does implement the concurrency model in question.  We'll
do something similar with CSP; we'll make sure to construct the CSP processes
for the concurrency model and for the reference implementation in such a way
that we can use an FDR refinement check to show the same thing.
