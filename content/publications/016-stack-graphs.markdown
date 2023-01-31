---
layout: /default.html
title: >-
  Stack graphs: Name resolution at scale
description: >-
  EVCS 2023, Delft, Netherlands
---

# Publication details

<ol class="publications">
<li class="reference">
  <div class="citation_number">[16]</div>
  <div class="citation">
    DA Creager and H van Antwerpen.

    “Stack graphs: Name resolution at scale.”

    In <em>Eelco Visser Commemorative Symposium (EVCS 2023)</em>.

    April 2023.

    Delft, Netherlands.
  </div>

  <div class="downloads">
    <a href="stack-graphs.pdf">
    <button type="button" class="btn btn-primary">
      <span class="glyphicon glyphicon-download"></span> PDF
    </button>
    </a>
  </div>
</li>
</ol>

## Abstract

We present <em>stack graphs</em>, an extension of Visser et al.'s scope graphs
framework. Stack graphs power Precise Code Navigation at GitHub, allowing users
to navigate name binding references both within and across repositories. Like
scope graphs, stack graphs encode the name binding information about a program
in a graph structure, in which paths represent valid name bindings. Resolving a
reference to its definition is then implemented with a simple path-finding
search.

GitHub hosts millions of repositories, containing petabytes of total code,
implemented in hundreds of different programming languages, and receiving
thousands of pushes per minute. To support this scale, we ensure that the graph
construction and path-finding judgments are <em>file-incremental</em>: for each
source file, we create an isolated subgraph without any knowledge of, or
visibility into, any other file in the program. This lets us eliminate the
storage and compute costs of reanalyzing file versions that we have already
seen. Since most commits change a small fraction of the files in a repository,
this greatly amortizes the operational costs of indexing large, frequently
changed repositories over time. To handle type-directed name lookups (which
require “pausing” the current lookup to resolve another name), our name
resolution algorithm maintains a stack of the currently paused (but still
pending) lookups. Stack graphs can be constructed via a purely syntactic
analysis of the program's source code, using a new declarative <em>graph
construction language</em>. This means that we can extract name binding
information for every repository without any per-package configuration, and
without having to invoke an arbitrary, untrusted, package-specific build
process.
