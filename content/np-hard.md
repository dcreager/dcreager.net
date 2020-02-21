+++
date = 2018-05-31
title = "NP-hard isn't the end of the world"
description = "learning to dodge"
[extra]
lastmod = 2018-06-12
+++

<p class=big-def>
  <em class=tldr>tl;dr</em> We've typically considered it a deal-breaker to
  discover that an algorithm we care about is <em>NP</em>-hard.  We'll go to
  great lengths to simplify the problem domain so that we can switch to a
  polynomial-time algorithm.  But if we simplify too much, then we run the risk
  that our solution is no longer useful.  Luckily, we might not have to!
  <em>NP</em>-hard is a <strong><em>worst-case</em></strong> bound.  If we can
  convince ourselves that we won't encounter pathological inputs, the
  <em>NP</em>-hard algorithm might be just fine in practice.
</p>

<!-- more -->

Complexity theory is the study of how difficult certain problems are.  We
typically measure a problem's complexity by estimating the amount of time and
memory that a computer would need to solve the problem for an arbitrary input.
We also like to lump together problems with basically the same difficulty into
**_complexity classes_**.  For instance, there is a complexity class named
**_P_** that contains all of the problems that are "easy enough" for computers
to solve.  (More formally, that can be solved in "polynomial time".)  There are
important differences between problems in _P_ — for instance, a constant-time
algorithm is definitely better than a linear-time one.  But if you zoom out far
enough, those differences stop mattering — an exponential algorithm is so much
worse than either that we can throw them both into _P_ and call it a day.

Another important complexity class is **_NP_**, which is where things first
start to get difficult for our poor little computers.  A problem in _NP_ is hard
to **_solve_** (as in, if you hand it the wrong input, there isn't enough time
left in the _life of the universe_ for your computer to finish churning away!)
But somewhat surprisingly, these problems are easy to **_verify_**.  That is, if
some oracle were able to magically produce a solution out of thin air, our
computer could easily (again, "in polynomial time") check that the solution is
correct.  This means that there's always a simple, but long-running, algorithm
for any _NP_ problem: list out every single possible solution, and check them
one by one.  There will always be exponentially many possible solutions, so
checking each one in turn will take exponential (i.e., really f---ing long)
time.

{% comment() %}
(As an aside, this is why everyone is so keen on quantum computing: a
quantum computer could solve an _NP_ problem easily, since it could check all of
the exponentially-many solutions at the same time.)
{% end %}

The last wrinkle is that even after decades upon decades, we still don't know
whether or not _P_ and _NP_ are the same thing.  We **_believe_** that there are
no easy solutions to any _NP_ problem, but haven't been able to **_prove_** it.
That leads us to the **_NP-hard_** problems.  These problems have some kind of
fundamental difficulty: if we were ever able to find an easy solution to an
_NP_-hard problem, we could use it to construct easy solutions to **_all_** _NP_
problems.  That, in turn, would prove that _P = NP_.  We believe so strongly
that _P ≠ NP_ that we're willing to base entire fields like cryptography on that
hunch.

All of this means that we have an (understandably) innate fear of _NP_-hard
problems.  With an algorithm that's just in _NP_, there's always a chance that
we just haven't thought about things long enough or hard enough — if we keep at
it, we'd see that while _this algorithm_ is hard, the underlying problem is not,
and we'd find a different algorithm that solves the problem more easily.  On the
other hand, if we discover that our algorithm is _NP_-hard, we're hosed.  To
crack this nut, we'd have to do something that none of the greatest minds in the
history of our field have done.

So instead of trying to solve the unsolvable, we look for a simpler problem to
solve.  For example, dependency management (like that provided by build tools
like [cargo][] and [npm][]) is a [provably _NP_-hard problem][dep np-hard] — you
can construct a set of libraries whose dependency relationships with each other
would cause `npm install` to run effectively forever.  This is...not the best
usability story, to say the least!

[cargo]: https://crates.io/
[npm]: https://www.npmjs.com/
[dep np-hard]: http://www.mancoosi.org/edos/manager/

While researching this problem for the Go ecosystem, Russ Cox pointed out that
you can [simplify the problem][cox] by prohibiting certain kinds of negative
dependencies (e.g., "the `kitten` package doesn't work with version 1.2.6 of the
`catnip` package").  That in turn, moves the dependency management problem from
_NP_-hard to _P_.  Perfect!  Now our tool will never blow up on us.

[cox]: https://research.swtch.com/vgo-mvs

Unfortunately, by simplifying the problem, we run the risk that our solution is
no longer useful.  Sam Boyer suggests [this is the case][boyer] for Go
dependency management: that there are real use cases where the community will
need the kinds of negative dependencies that would be outlawed.

[boyer]: https://sdboyer.io/vgo/failure-modes/

That leaves us with a conundrum.  On the one hand, we can't simplify our problem
without sacrificing real use cases.  On the other, we have an algorithm that
might not finish in any reasonable amount of time.  What can we do?

I would argue that there's a third possibilty, which is to just...shrug.
_NP_-hard is a **_worst-case_** bound on your algorithm.  It might only apply to
"pathological" inputs.  For dependency management, for instance, we said that
_"you could construct"_ a dependency graph that causes `npm install` to blow up.
But how many hoops do you have to jump through to construct that graph?  Are you
**_really_** going to encounter a graph like that in practice?  If not, we don't
**_really_** have to worry about the worst-case complexity.

Now, this doesn't necessarily apply to all _NP_-hard problems.  You might be
dealing with something so fundamentally difficult that all interesting inputs
take exponential time to solve.  If that's the case...well, this post isn't
going to help you all that much.  You'll still have to simplify your problem
domain somewhow to make any forward progress.

But before you go down that road, it's worth thinking through some real-world
use cases.  You might be pleasantly surprised!

<p class=thanks>
  If you have a hankering for formal methods, check out <a
  href="/publications/014-csp-algorithm-study/">this paper</a> that I wrote
  during my grad school days, where I use a CSP refinement checker to explore
  the complexity space of an interesting <em>NP</em>-hard problem, trying to
  find what kinds of inputs lead to the worst-case exponential behavior.
</p>
