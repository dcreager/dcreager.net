---
kind: article
created_at: 2016-08-19
layout: post
series: "Concurrency models in CSP"
title: "Read Atomic: Internal consistency"
tags: [csp]
draft: true
---

We'll start by looking at the weakest concurrency model covered in the paper,
Read Atomic.  All of the concurrency models are defined as a combination of
several axioms.  Read Atomic consists of only two: internal consistency and
external consistency.  In this post, we'll construct a CSP process that
describes the first of the two.

<!--
 Read Atomic...can be implemented without requiring any coordination among
 replicas...a replica can decide to commit a transaction without consulting
 other replicas.
-->

A transaction is *internally consistent* if it "reads its own writes".  This is
the simplest axiom covered in the paper, since it expresses a property that's
strictly local to each transaction; we don't have to consider the behavior of
any other transaction when defining what it means for a transaction to be
internally consistent.

We start by defining some data types and channels:

<% highlight :csp do %>
datatype Object = {1..2}
datatype Value = {1..5}
<% end %>

An `Object` is the *name* of one of the objects that can be read or written; it
corresponds to the \\(Obj\\) set in the original paper.  The `Value` type is the
set of values that can be stored in each object.  In the original paper, they
use \\(\mathbb{Z}\\) (the set of integers); in our CSP spec, we'll also use
simple integers, but we'll limit ourselves to a finite set.

We can then use \\(\mathbf{channel}\\) statements to define some events:

<% highlight :csp do %>
channel read : Object x Value
channel write : Object x Value
<% end %>

A `read` event signifies that a particular value was read for a particular
object; similarly, a `write` event signifies that a particular value was written
to a particular object.

Based on the structure of these two CSP channels, you might expect that they
correspond to the \\(Op\\) set from the original paper; however, they actually
correspond to the \\(REvent\_x\\) and \\(WEvent\_x\\) family of sets.  In the
original paper's formalism, the authors needed to include an \\(EventId\\) in
each event to describe the ordering between events.  In CSP, on the other hand,
the prefix operator (\\(\rightarrow\\)) lets us define how the events in a
process are ordered without having to add extra fields to the events.

With our datatypes and events defined, we can now create a process that
specifies what "internal consistency" means.  For each transaction, each object
can be in one of two states; we'll use subprocesses to define what internal
consistency means in each of these two states.

Each object starts *undefined*, meaning that *this transaction* has not yet
written a value to the object:

<% highlight :csp do %>
Undefined(obj) =
  read!obj?value -> Undefined(obj)
    []
  write!obj?value -> Defined(obj, value)
<% end %>

(Remember that internal consistency does not take into account what values other
transactions write into an object.)

Note that we have to include a `read` clause, even though internal consistency
doesn't say anything about which value is returned when we read from an object
that we didn't already write to.  This `read` clause tells us that we have no
idea which value will be returned.  We can't leave this clause out; that would
say that internal consistency *actively prevents* transactions from performing
these kinds of reads.  Instead, we'll leave these reads unconstrained here, and
use other CSP processes to define how each consistency model determines which
values are returned based on which other transactions have completed and become
visible.

The `write` clause means that you can write (whatever value you want) to an
undefined object; doing so means us into the "defined" state.

When an object is *defined*, we also need to keep track of the current value:

<% highlight :csp do %>
Defined(obj, currentValue) =
  read!obj!currentValue -> Defined(obj, currentValue)
    []
  write!obj?newValue -> Defined(obj, newValue)
<% end %>

When you `read` a defined value, we enforce that the value that's returned is
the current value of the object — that's the entire point of internal
consistency!  You can also `write` a *new* value to the object; doing means that
the object is still "defined", but overwrites the current value with the new
one.
