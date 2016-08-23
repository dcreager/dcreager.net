---
kind: article
created_at: 2016-08-19
layout: post
series: "Concurrency models in CSP"
title: "Read Atomic: Internal consistency"
tags: [csp]
draft: true
---

{::options parse_block_html="true" /}

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

We start by defining some data types and channels.

<div class="aside-def">
## Channels
Each \\(\textbf{channel}\\) statement defines one kind of event that your
processes can produce.  Each event will typically carry some data; you use
\\(\textbf{datatype}\\) and \\(\textbf{nametype}\\) statements to define the
different kinds of data that can be carried by those events.
</div>

{::comment}
nametype Object = {1..2}
nametype Value = {1..5}
{:/comment}

\\[
\textbf{nametype} ~ \textsf{Object} = \\{1..2\\} \\\
\textbf{nametype} ~ \textsf{Value} = \\{1..5\\}
\\]

An \\(\textsf{Object}\\) is the *name* of one of the objects that can be read or
written; it corresponds to the \\(Obj\\) set in the original paper.  The
\\(\textsf{Value}\\) type is the set of values that can be stored in each
object.  In the original paper, they use \\(\mathbb{Z}\\) (the set of integers);
in our CSP spec, we'll also use simple integers, but we'll limit ourselves to a
finite set.

{::comment}
channel read : Object.Value
channel write : Object.Value
{:/comment}

\\[
\textbf{channel} ~ \texttt{read} : \textsf{Object} \times \textsf{Value} \\\
\textbf{channel} ~ \texttt{write} : \textsf{Object} \times \textsf{Value}
\\]

A \\(\texttt{read}\\) event signifies that a particular value was read for a
particular object; similarly, a \\(\texttt{write}\\) event signifies that a
particular value was written to a particular object.

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

{::comment}
Undefined(obj) =
  read!obj?value -> Undefined(obj)
    []
  write!obj?value -> Defined(obj, value)
{:/comment}

\\[
\textrm{Undefined}(\textit{obj}) = {} \\\
  \quad\texttt{read}\,!\textit{obj}\,?\textit{value} \rightarrow
        \textrm{Undefined}(\textit{obj}) \\\
  \qquad \Box \\\
  \quad\texttt{write}\,!\textit{obj}\,?\textit{value} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{value})
\\]

(Remember that internal consistency does not take into account what values other
transactions write into an object.)

Note that we have to include a \\(\texttt{read}\\) clause, even though internal
consistency doesn't say anything about which value is returned when we read from
an object that we didn't already write to.  This \\(\texttt{read}\\) clause
tells us that we have no idea which value will be returned.  We can't leave this
clause out; that would say that internal consistency *actively prevents*
transactions from performing these kinds of reads.  Instead, we'll leave these
reads unconstrained here, and use other CSP processes to define how each
consistency model determines which values are returned based on which other
transactions have completed and become visible.

The \\(\texttt{write}\\) clause means that you can write (whatever value you
want) to an undefined object; doing so means us into the "defined" state.

When an object is *defined*, we also need to keep track of the current value:

{::comment}
Defined(obj, currentValue) =
  read!obj!currentValue -> Defined(obj, currentValue)
    []
  write!obj?newValue -> Defined(obj, newValue)
{:/comment}

\\[
\textrm{Defined}(\textit{obj}, \textit{currentValue}) = {} \\\
  \quad\texttt{read}\,!\textit{obj}\,!\textit{currentValue} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{currentValue}) \\\
  \qquad \Box \\\
  \quad\texttt{write}\,!\textit{obj}\,?\textit{newValue} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{newValue})
\\]

When you \\(\texttt{read}\\) a defined value, we enforce that the value that's
returned is the current value of the object — that's the entire point of
internal consistency!  You can also \\(\texttt{write}\\) a *new* value to the
object; doing means that the object is still "defined", but overwrites the
current value with the new one.
