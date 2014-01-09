---
kind: article
created_at: 2014-01-07
layout: post
title: "CSP: An introduction"
tags: [csp]
---

Communicating Sequential Processes (CSP) has been around for almost four decades
at this point, but for much of its life, it was only well-known among
theoretical computer scientists and formal methods advocates.  More recently,
many more people have at least *heard* of CSP, largely because it inspired the
[concurrency support](http://golang.org/doc/effective_go.html#concurrency) in
[Go](http://golang.org/), a popular mainstream programming language.  However,
if you ask most people what it *means* to be inspired by CSP, the most common
response would probably be "erm, something about message passing"?

That said, CSP isn't just some dusty theory that inspired part of Go; it can
also help us understand the distributed systems that we create.  We've developed
a [plethora](http://zookeeper.apache.org/) [of](http://www.mongodb.org/)
[tools](http://couchdb.apache.org/) [that](http://redis.io/)
[help](http://basho.com/riak/) [us](http://cassandra.apache.org/)
[build](http://hbase.apache.org/) [distributed](http://hadoop.apache.org/)
[systems](http://storm-project.net/).  But unfortunately, we don't always
understand of how those tools work, how they fail, and how they interact when we
piece them together into a larger system.  We can all name-drop the [CAP
theorem](http://dl.acm.org/citation.cfm?id=564601), but do you *really* know
what your system is going to do when the network partitions, or when a host
dies?  How do you convince someone that you're right?

We can't just rely on intuition and hand-wavy arguments; our systems are too
large, and too important, for that.  So how do you address these concerns with
rigor?  There are two main approaches: you can either *test* your assumptions
empirically on a running system, or you can describe your system in detail and
*prove* that your assumptions are correct.  Kyle Kingsbury has great examples of
both: [Jepsen](http://aphyr.com/tags/jepsen) on the testing side,
[Knossos](http://aphyr.com/posts/309-knossos-redis-and-linearizability) on the
proof side.  Both approaches are important; if you want to be convincing, you
have to choose at least one of them.  If you prefer the proof-based approach,
CSP is another option.  If you only think of CSP in terms of Go's concurrency
primitives, or if you haven't thought of it at all, then you overlook the fact
that CSP was *specifically designed* to help answer these kinds of questions.

In this series of articles, I want to describe how CSP fits into this landscape,
for developers with a range of expertise.  For the every-day programmer, I want
to give a basic, high-level introduction to CSP, and to explain what it means
for Go to be inspired by CSP.  For the distributed systems engineer, I want to
add weight to the argument that formal methods are a useful tool for studying
and designing the systems that we create and use.  And for the formal methodist,
I want to show how to use CSP in particular to specify and reason about those
systems.
