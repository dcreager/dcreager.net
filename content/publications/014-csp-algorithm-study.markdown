---
layout: /default.html
title: Publication details
---

# Publication details

<ol class="publications">
<li class="reference">
  <div class="citation_number">[14]</div>
  <div class="citation">
    DA Creager and AC Simpson.

    “Empirical analysis and optimization of an <em>NP</em>-hard
    algorithm using CSP and FDR.”

    In <em>Proceedings of the Brazilian Symposium on Formal Methods
      (SBMF 2007)</em>.

    ENTCS, to appear.
  </div>

  <div class="downloads">
    <a href="csp-algorithm-study.pdf">
    <button type="button" class="btn btn-primary">
      <span class="glyphicon glyphicon-download"></span> PDF
    </button>
    </a>
  </div>
</li>
</ol>

## Notes

This is an extended version of the
[paper](../011-csp-algorithm-study/) that I presented at [SBMF
2007](http://www.sbmf2007.ufop.br/) in Ouro Preto, Brazil.

## Abstract

In many cases where an algorithm is provably *NP*-hard, this
intractability is a worst-case bound that only applies to pathological
inputs.  In these cases, by exploiting knowledge of the specific
structure of “real-world” inputs, the algorithm can be shown to be
much more efficient in the “normal” case.  However, when studying a
new problem, this can be hard to show if it is not obvious which
structural constraints exist, and which ones would lead to increases
in efficiency.  In this paper, we show how one can describe the
underlying problem declaratively as a CSP process and use the FDR
refinement checker to explore the complexity space of the problem.  By
knowing which optimizations FDR uses to find solutions more
efficiently, we can determine under which conditions the worst-case
intractable algorithm executes efficiently, and incorporate analogous
optimizations into the algorithm to exploit these conditions.
