---
kind: article
created_at: 2014-06-08
layout: post
title: "CSP: Events and processes"
tags: [csp]
draft: true
---

Okay, we have [an overview](/2014/02/csp-basics/) that gives us the big picture
about what CSP is, and how you could use it.  But I always find it easier to
learn from concrete examples.  In the grand tradition of basically every other
introduction to CSP, we're going to create a formal model of a vending machine.
We'll start with a very simple CSP model, and then add in more and more detail
as we go, learning new CSP concepts as we need them.


## Processes

When we model a system using CSP, that system is represented by a *process*.  If
you come from a programming background, this is not the same thing as an OS
process; just think of it as "something that can do something".  Very
loosy-goosy, at least to start with.

We use CSP to model a system because we want to describe how that system works
internally, and more importantly, how it interacts with the world around it.
These interactions are represented by *events*, and "the world around it" is
represented by an *environment*.  When we talk about a process in
CSP, basically the only thing that we care about is what patterns of events that
it can perform.

<div class="defn">
<em>Processes</em> interact with its <em>environment</em> via <em>events</em>.
</div>

If all that we care about is what events a process can perform, then the
simplest process is one that can't do anything.  In CSP, this process is called
`Stop`.



<div class="alert alert-info">
<h3>Two kinds of notation</h3>
But first, a digression.
</div>
