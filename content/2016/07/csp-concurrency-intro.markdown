---
kind: article
created_at: 2016-07-28
layout: post
series: "Concurrency models in CSP"
title: "Introduction"
tags: [csp]
---

Two years ago, I started writing a series of [blog](/2014/01/07/intro-to-csp/)
[posts](/2014/02/csp-basics/) about CSP — Communicating Sequential Processes.
That series has...let's say "stalled".  Part of the reason is that I didn't have
a good non-trivial example to work with.  The [stereotypical running
example](https://en.wikipedia.org/wiki/Communicating_sequential_processes#Examples)
is the vending machine: start with a simple one, which accepts a single coin and
spits out a tea; add more detail as you introduce more of the language.  I
always had a hunch that I needed an example with more meat on it, but could
never find one.

Fast-forward to today.  I was reading through some of [Christopher
Meiklejohn's](https://christophermeiklejohn.com/) work on
[Lasp](http://lasp-lang.org/), I came across a citation to a really nice paper
by [Cerone, Bernardi, and
Gotsman](http://drops.dagstuhl.de/opus/volltexte/2015/5375/), which adds some
formal rigor to the consistency models that we use to describe modern
distributed systems.  Their formalism is a great combination of simple and
expressive.  The core of the paper is about processes accessing a transactional
data store; the authors provide formal definitions of several concurrency
models, and of some reference implementations that supposedly provide those
concurrency models.  They then use a technique called "observational refinement"
to show that the reference implementations really do provide the concurrency
guarantees in question.

This approach lines up very well with how you perform refinement checks in CSP
to show that systems satisfy some specification.  And so I finally found my
meaty running example!  I'm resurrecting this blog series, and plan to work
through each of the consistency models and proofs described in the paper,
translating them into CSP processes and refinement checks.  This isn't an
attempt to replace or outdo anything in the paper!  Far from it — it's my
attempt to use something more familiar to work through the details of something
less familiar.

I'm not going to assume a working knowledge of CSP — or of the consistency
models described in the paper!  If you're familiar with one, my hope is that
you'll be able to follow along and pick up the other.  And if you're not
familiar with either...well, I guess we'll see how good I am at writing an intro
to a difficult topic!
