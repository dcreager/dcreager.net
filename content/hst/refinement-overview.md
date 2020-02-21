+++
date = 2016-11-17T08:00:00Z
title = "Refinement overview"
description = "in which we see what's ahead of us and immediately regret this decision"
[extra]
#tags: [csp]
series = "HST"
#prev = "hst/intro.md"
next = "hst/semantic-methods.md"
+++

Our goal is to learn about CSP refinement by implementing a refinement checker.
So a good first step is to make sure we're all on the same page about what
refinement is, and then to step through the refinement algorithm that we mean to
implement.  (If nothing else, that will help make sure I don't go off on too
many tangents while implementing it!)

I've mentioned refinement elsewhere on this blog a few times (for instance,
[here](/csp-concurrency/read-atomic-internal/#refinement)).  The basic idea is
that in CSP, you use the same process language to describe the system you're
designing or investigating, as well as the properties that you would like that
system to have.  (This is unlike most other formal methods, where you have
separate languages for the system and the properties.)  In CSP, the system's
process is typically called \\(Impl\\) (for **implementation**), and the
property description process is typically called \\(Spec\\) (for
**specification**).

CSP then defines several **semantic models** that provide rigorous mathematical
definitions of what a process's behavior is.  You perform a refinement check
within the context of a particular semantic model.  A successful refinement
check tells you that the property defined by \\(Spec\\) "holds" — specifically,
that all of the behaviors of \\(Impl\\) are also allowed behaviors of
\\(Spec\\).  A failed refinement check gives you a **counterexample** — that is,
a specific behavior of \\(Impl\\) that was disallowed by \\(Spec\\).

<!-- more -->

The three most common semantic models are **traces**, **failures**, and
**failures-divergences**.  We'll go into more detail about the mathematics
behind these semantic models in later posts; for now, the 10,000-foot overview
is that:

  - Traces refinements let you check **safety** properties (i.e., that something
    bad is not allowed to occur).
  - Failures refinements let you check **liveness** properties (i.e., that
    something good *must* occur).
  - Failures-divergences refinements (unlike the first two) work in the presence
    of **endless loops**.

In my post about the [Read Atomic concurreny model][RA], I use a traces
refinement check to verify that Read Atomic doesn't allow "unrepeatable reads".
In this example, the \\(Spec\\) process is a description of the Read Atomic
concurrency model, while the \\(Impl\\) process is a "fake" implementation that
immediately tries to perform an unrepeatable read.  Because that unrepeatable
read isn't allowed by the Read Atomic process, the traces refinement check
fails.

[RA]: /csp-concurrency/read-atomic-internal/#testing-it-out

Given all of that, how do we write a program that can perform refinement checks
for us?  The answer is strewn throughout Bill Roscoe's textbook, [*The theory
and practice of concurrency*][textbook].  The bulk of FDR's refinement algorithm
is described in Appendix C (p. 541).  At a high level, we need to:

{% comment() %}
I'll include references to specific sections and pages where appropriate; I'm
specifically looking at the "lightly revised to 2005" edition that's currently
available on Roscoe's Oxford site.
{% end %}

[textbook]: https://www.cs.ox.ac.uk/bill.roscoe/publications/68b.pdf

  1. Load in a description of the \\(Spec\\) and \\(Impl\\) processes,
     transforming them each into a **labeled transition system (LTS)**.

  2. Normalize the \\(Spec\\) process, resulting in a **normalized LTS**.

  3. Perform a **simultaneous breadth-first search** through the \\(Spec\\)'s
     normalized LTS and \\(Impl\\)'s (non-normalized) LTS, looking for a
     counterexample to the refinement.

  4. If we find any counterexample, the refinement check fails.  If we don't,
     the refinement check succeeds.

This is enough to get started for now; in later posts I'll drill down into each
of these steps in (much) more detail, and show how to implement them in HST.
