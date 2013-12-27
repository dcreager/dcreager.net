---
layout: default
title: Publication details
---

# Publication details

<ol class="publications">
<li class="reference">
  <div class="citation_number">[12]</div>
  <div class="citation">
    DA Creager.

    <em>A graph-based approach to the automated discovery of data
    transformations</em>.

    D.Phil. thesis, University of Oxford, August 2007.
  </div>

  <div class="downloads">
    <a href="phd.pdf">
    <button type="button" class="btn btn-primary">
      <span class="glyphicon glyphicon-download"></span> PDF
    </button>
    </a>
  </div>
</li>
</ol>

## Abstract

In recent years it has become much more common for software
applications to communicate with each other directly.  Internet
connections have become a standard part of both ofﬁce and home, and as
more business processes and information move into the electronic
realm, direct software communication will become even more prevalent.
One of the largest deterrents to effective communication is the
heterogeneous nature of the data and information involved.  We cannot
guarantee that two software systems that need to communicate will be
running the same software or modeling their data in the same way.
Obviously, the data must be somehow logically similar — otherwise,
there would not be any meaningful communication possible.  A key
element of any modern communication protocol or framework must be a
strategy for resolving any data mismatches that exist between the two
sides.

The data mismatch problem is not new; unsurprisingly, there are many
existing solutions to it.  We would like to judge these solutions by
two criteria: generality and automation.  A generic solution will not
needlessly limit the kinds of applications and data models that are
supported.  An automated solution will limit the amount of tedious,
manual work needed to support a new application or data model.
Unfortunately, none of the existing solutions are both sufficiently
generic and sufficiently automated.

This thesis presents an automated solution to the data mismatch
problem that is also fully generic: it makes absolutely no assumptions
about the underlying data whatsoever.  In order to achieve this
generality, some automation must be sacriﬁced.  Our approach requires
that some atomic transformations be written manually.  However, we can
exploit the fact that transformations are composable — with a
sufficient number of atomic transformations, a compound transformation
can be automatically discovered between arbitrary datatypes.  This
approach is fully generic, since the transformation discovery
algorithms require no knowledge of the structure or semantics of the
datatypes involved; instead, the knowledge of a particular datatype is
encapsulated into the atomic transformations that directly operate on
it.

The contributions of this thesis are threefold.  First, we present a
graph-based model for transformations that has an efficient
polynomial-time discovery algorithm.  While efficient, this model is
limited in that it can only represent unary transformations — those
between one input and output datatype.  This model is still
surprisingly powerful; we present two case studies that show how this
simple model can be used in the context of a real-world application,
and what limitations it has.

Second, we present an extension to this graph-based model that
supports polyadic transformations between multiple input and output
datatypes, and show examples of how this increases the expressive
power of the transformation graphs that we can create.  Unfortunately,
this expressiveness comes at a price: the naïve discovery algorithm
for the new model runs in exponential space and time.

Finally, we show that polyadic transformation discovery is in fact
worst-case *NP*-hard.  Hopefully, the problem is only truly
intractable for pathological inputs, and for real-world transformation
graphs, compound transformations can be discovered with reasonable
time and space requirements.  We use a novel application of CSP to
test this hypothesis, empirically exploring the complexity space of
the problem and highlighting criteria for designing efficient
transformation graphs.
