---
kind: article
created_at: 2017-01-17
layout: /post.html
series: "HST"
title: "Lazy processes"
description: >-
  in which we reaffirm that laziness is a virtue
tags: [csp]
prev: /2016/11/semantic-methods.markdown
---

{::options parse_block_html="true" /}

<script type="text/x-mathjax-config">
MathJax.Hub.Config({
  TeX: {
    Macros: {
      interleavesym: "|\\mkern-2mu|\\mkern-2mu|",
      interleave: "\\mathrel{\\interleavesym}",
    }
  }
});
</script>

{:start="1" .csp-algo-step}
  1. Load in a description of the \\(Spec\\) and \\(Impl\\) processes,
     transforming them each into a **labeled transition system (LTS)**.

In this post we're going to look at how to represent processes in code.  As
mentioned in the [previous post][semantic methods], we're going to rely on the
*labeled transition system* defined by CSP's *operational semantics*.  An LTS is
just a directed graph, with nodes representing processes and subprocesses, and
edges representing events.  Our question then becomes how to represent this
graph.

[semantic methods]: /2016/11/semantic-methods/

The simplest answer would be to store the LTS graph directly.  You wouldn't even
have to write very much code, since you can find a decent pre-canned graph
library for just about any programming language you can think of.  You'd walk
through the operational semantics rules of the process's description, creating
explicit nodes and edges in the graph for each subprocess and transition.

This approach has some real advantages.  First off, it's conceptually simple:
there's a physical graph edge in memory for each transition rule describe by the
operational semantics.  Second, if you choose a good graph representation, an
explicit LTS can be very dense, memory-efficient, and cache-friendly.

However, there is also one very large drawback: CSP processes can get very big,
especially once you start using the parallel composition operators.  And by
"big", I mean "large enough to exhaust your physical memory".  This is the
"state space explosion" problem that is the bane of any exhaustive formal model
checking technique.

<hr class="jump">

As a somewhat contrived example, consider the following processes:

{::comment}
A(x) = a!x -> x > 1 & A(x-1)
B(x) = b!x -> x > 1 & B(x-1)
{:/comment}

\\[
\textrm{A}(\textit{x}) =
  \texttt{a}\,!\textit{x} \rightarrow \textit{x} > 0 \mathop{\&}
  \textrm{A}(\textit{x} - 1) \\
\\]

These processes just count down from \\(\textit{x}\\) to \\(1\\), producing one
event for each value.  For instance, \\(\textrm{A}(2)\\) is equivalent to:

\\[
  \texttt{a}{.}2 \rightarrow \texttt{a}{.}1 \rightarrow \textrm{STOP}
\\]

By themselves, these processes aren't that big of a deal.  Things become
unwieldy when we start combining them together.  For instance, we can interleave
the processes together — \\(\textrm{A}(2) \interleave \textrm{B}(2)\\) — giving
us the following LTS:

<%= figure "csp/a2-interleave-b2" %>

This interleaving has six possible traces:

\\[
\langle \texttt{a}{.}2, \texttt{a}{.}1, \texttt{b}{.}2, \texttt{b}{.}1 \rangle \quad\quad
\langle \texttt{b}{.}2, \texttt{b}{.}1, \texttt{a}{.}2, \texttt{a}{.}1 \rangle\\\
\langle \texttt{a}{.}2, \texttt{b}{.}2, \texttt{a}{.}1, \texttt{b}{.}1 \rangle \quad\quad
\langle \texttt{b}{.}2, \texttt{a}{.}2, \texttt{b}{.}1, \texttt{a}{.}1 \rangle\\\
\langle \texttt{a}{.}2, \texttt{b}{.}2, \texttt{b}{.}1, \texttt{a}{.}1 \rangle \quad\quad
\langle \texttt{b}{.}2, \texttt{a}{.}2, \texttt{a}{.}1, \texttt{b}{.}1 \rangle
\\]

The number of possible interleavings is \\(O(n!)\\) in the length of the
processes.  Bumping up to \\(\textrm{A}(20) \interleave \textrm{B}(20)\\), we
already have **137 billion** possible interleavings.  Going further to
\\(\textrm{A}(50) \interleave \textrm{B}(50)\\), we have **\\(1 \times
10^{29}\\)** possible interleavings.  Fifty possible internal states is not
outrageously large, but combining two reasonably-sized processes gives us
something that we can't possibly load into memory at once.  We want to avoid
creating a full LTS in advance for a large process like this, *especially* if we
don't need to analyze every subprocess to determine whether a particular
refinement check holds.  This puts us in a bit of a pickle: we want to avoid
explicitly creating the full LTS graph of a process, but we still need walk
through that LTS graph as part of a refinement check.

To get around this problem, we need to implement some kind of **laziness**.
Instead of storing the LTS directly, we store a "recipe" for constructing the
LTS on the fly, as we need it.

We can implement laziness in a number of different ways.  Roscoe spends several
pages describing FDR's approach: "supercompilation".  Briefly, supercompilation
defines an internal language (almost like a bytecode) that can be used to encode
the operational semantics rules of each CSP operator.

In HST, we're going to take a different approach.  Instead of using something
like a bytecode to represent the transition rules, we're going to use code
itself, using two common object-oriented patterns: **interfaces** and
**visitors**.  For each CSP operator, we'll implement the following interface:

    interface Process
      initials(visitor: EventVisitor)
      afters(initial: Event, visitor: EdgeVisitor)

    interface EventVisitor
      visit(event: Event)

    interface EdgeVisitor
      visit(event: Event, after: Process)

I've used pseudocode to describe the interfaces, since every programming
language gives us some construct that we can use to implement this idea.  You'd
use classes in C++ or Java, traits in Rust, interfaces in Go, higher-order
functions in Clojure or Haskell.  I happen to be using C for my first stab at
HST; in later posts, we'll see how to use "interface structs" to implement this
pattern in C.

Regardless of which language you use, this pattern means that you don't have to
instantiate the in-memory representation of any particular sub-process until you
actually need it.



{::comment}
Under the covers, we're using the [Judy][] library to do the heavy lifting.
[Judy]: http://judy.sourceforge.net/
{:/comment}
