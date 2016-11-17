---
kind: article
created_at: 2016-11-16
layout: post
series: "HST"
title: "Introduction"
description: >-
  in which we learn by implementing
tags: [csp]
next: /2016/11/refinement-overview/
---

It seems like this blog is basically turning into "all things CSP"!  As part of
that trend, I've started implementing a new CSP refinement checker called
[HST][].  Why do this when there's a perfectly good refinement checker in
[FDR][]?  Well, I want to learn more about how FDR's refinement algorithm works.
The algorithm is documented in Bill Roscoe's [textbook][] (and a series of
follow-on papers), and working through those descriptions gives you a good bit
of insight into how refinenment really works.  But I often find it easier to
learn a complex topic by implementing it (or at least, by looking at the code of
an implementation).  Hence HST!  In this new series of blog posts, I'm going to
walk through the CSP refinement algorithm in more detail than is presented in
the academic literature, by implementing it (and describing that implementation)
along the way.

I should emphasize that this is **not** meant to be a replacement for FDR!  FDR
is a very good piece of software, and if you're writing any CSP specs in anger,
you probably want FDR at your disposal.  HST is meant to be more of an
educational exercise.  If people find it more generally useful than that, that's
great!  But it's not what I'm aiming for.

[textbook]: https://www.cs.ox.ac.uk/bill.roscoe/publications/68b.pdf
[FDR]: http://www.cs.ox.ac.uk/projects/fdr/
[HST]: https://github.com/hst/hst/

(And as for the name, "HST" does *not* stand for "Harry S. Truman", just like
"FDR" does *not* stand for "Franklin Delano Roosevelt".)
