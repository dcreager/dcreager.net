---
kind: article
created_at: 2016-11-17 09:00
layout: post
series: "HST"
title: "Semantic methods"
description: >-
  in which all good things come in threes
tags: [csp]
prev: /2016/11/refinement-overview/
---

<script type="text/x-mathjax-config">
MathJax.Hub.Config({
  TeX: {
    Macros: {
      interleavesym: "|\\mkern-2mu|\\mkern-2mu|",
      interleave: "\\mathrel{\\interleavesym}",
      Interleave: "\\mathop{\\interleavesym}"
    }
  }
});
</script>

Since CSP is a formal method, it's not surprising that Roscoe spends a large
part of his textbook talking about how to use rigorous mathematics to talk about
processes.  He actually goes one step (er, two?) further and defines **three
different** ways to do so:

  - The **denotational semantics** defines (mathematically, using sets and
    sequences) what the *behavior* of a process is.  Each CSP operator comes
    with a rule for how to calculate a process's behavior recursively — that is,
    in terms of the behavior of its operands.  (So for example, the "external
    choice" rule tells you how to define the behavior of \\(P \mathrel{\Box}
    Q\\) in terms of the behavior of \\(P\\) and \\(Q\\).)

  - The **algebraic semantics** tell you to not worry about what a process
    "means" or what it "does".  Instead, it provides a list of *rewrite rules*
    that let you change what the definition of a process looks like without
    changing its behavior.

  - The **operational semantics** say that a process is nothing more than a
    state machine, with nodes representing processes (and subprocesses) and
    edges representing the events that allow you to transition between them.  We
    can learn anything important about a process just by interpreting or
    analyzing this state machine.

One of the important contributions of the book is to not just describe these
three different semantic methods (in detail), but to show that they're
**equivalent**.  This is great, because all three semantics are useful in
different situations:

  - The denotational semantics are where the concept of refinement is actually
    defined, where we have the best intuition about what it means, and give us a
    language for describing counterexamples.

  - The algebraic semantics give us a way to simplify processes by "turning the
    crank"; we make superficial syntactic transformations to a process in a way
    that hopefully leads us to a simpler definition.

  - The operational semantics give us a data structure for representing
    processes that is easy to program with.

And more importantly, we're going to pick and choose the most useful pieces of
each semantic method as part of the refinement algorithm that we describes in
the previous post:

  1. Load in a description of the \\(Spec\\) and \\(Impl\\) processes,
     transforming them each into a **labeled transition system (LTS)** *(using
     the rules from the operational semantics)*.

  2. Normalize the \\(Spec\\) process, resulting in a **normalized LTS** *(using
     the normalization operation from the algebraic semantics)*.

  3. Perform a **simultaneous breadth-first search** through the \\(Spec\\)'s
     normalized LTS and \\(Impl\\)'s (non-normalized) LTS, looking for a
     counterexample to the refinement *(taking advantage of the fact that the
     operational and denotational semantics are equivalent)*.

  4. If we find any counterexample, the refinement check fails.  If we don't,
     the refinement check succeeds.  *(In either case, we describe the result in
     terms of the denotational semantics.)*
