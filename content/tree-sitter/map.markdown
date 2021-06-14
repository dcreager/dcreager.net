---
kind: article
created_at: 2021-06-14
layout: /post.html
title: A map of the tree-sitter ecosystem
description: terra incognita
---

The tree-sitter ecosystem is divided up across a large number of components,
each in different repositories, which can be quite overwhelming at first.  This
post tries to provide a map of sorts.

## Overview

Say you're interested in the tree-sitter project, so you decide to check out the
[`tree-sitter` organization][org] on GitHub, browsing through its repositories to
determine how the ecosystem is structured.  The list of repositories spills over
onto a second page, and you see entries that seem redundant.  Why is there both
[`tree-sitter-python`][] and [`py-tree-sitter`][]?  Are they competing with each
other?  Is one deprecated?

[org]: https://github.com/orgs/tree-sitter/
[`tree-sitter-python`]: https://github.com/tree-sitter/tree-sitter-python/
[`py-tree-sitter`]: https://github.com/tree-sitter/py-tree-sitter/

You might instead decide to check out the [project homepage][homepage].  The
landing page lists (as of June 2021) over 40 different programming language
parsers that various folks have implemented, as well as a handful of language
bindings.

[homepage]: https://tree-sitter.github.io/

This, at least, points to an answer.  The tree-sitter ecosystem is complicated
because when we write a code analysis tool, we want to support different
programming languages in two separate, orthogonal ways:

* First, we want to be able to _parse_ source code implemented in different
  programming languages.

* Second, and possibly less obviously, we want to _use_ tree-sitter in several
  different programming languages.  You specifically are going to write your
  analysis tool in one language, but we (the tree-sitter developers) don't know
  which one that is!  We've tried to implement tree-sitter so that we don't
  place any restrictions on which language you use.

That at least explains why “Python support” in tree-sitter might mean two
different things.  But why have we separated everything out into distinct
repositories?  The main reason is to make it as clear as possible that all of
these pieces are truly independent of each other.  There shouldn't be any way
for the Python language bindings to influence the design or release process of
the Haskell bindings, for instance, nor of _any_ of the language grammars.

True, it adds complexity to the ecosystem, but we've tried to get around this
with careful naming conventions, and tree-sitter-specific tooling to make it
easy to find and work with whatever pieces you need.

So, given the above, you will encounter all of the following on your journey:

### Language parsers

You must have a tree-sitter grammar for each language that you want to parse.
Each language grammar is typically implemented in a its own repository, named
`tree-sitter-$LANGUAGE`.

- [JavaScript grammar](https://github.com/tree-sitter/tree-sitter-javascript)
- [TypeScript grammar](https://github.com/tree-sitter/tree-sitter-typescript)
- [Python grammar](https://github.com/tree-sitter/tree-sitter-python)

There are some exceptions.  For instance, the `tree-sitter-javascript`
repository lets you parse JavaScript _and_ JSX — although in this case, this is
handled with a single grammar that treats “plain JavaScript” as a file that
happens to not have any JSX expressions in it.  Similarly, the
`tree-sitter-typescript` repository lets you parse TypeScript and TSX, though in
this case, they're handled with distinct grammars.  All of these grammars share
enough structure, and are a coherent enough family of languages, that it would
be overkill to separate them out further.

### The tree-sitter runtime library

The generated parsers only contain some state tables describing the language
being parsed. The “meat” of the parsing logic is implemented in the
`tree-sitter` runtime library, which each parser depends on.  This runtime
library is also where tree-sitter's [query language][] is implemented.

[query language]: https://tree-sitter.github.io/tree-sitter/using-parsers#pattern-matching-with-queries

The runtime library is implemented in the
[`tree-sitter/tree-sitter`][`tree-sitter`] repository on GitHub, under the
`lib/include` and `lib/src` directories.

[`tree-sitter`]: https://github.com/tree-sitter/tree-sitter

### Language bindings

The runtime library and each generated parser are implemented in C.  Assuming
that you aren't writing your analysis tool in C, you will need _bindings_ for
the language that you are using.  This will use your language's FFI mechanism to
link in the tree-sitter C code and make it available using more idiomatic
constructs.

The [Rust][rust binding] and [WASM][wasm binding] bindings are considered “tier
1”, and are implemented directly in the `tree-sitter/tree-sitter` repository.

[Rust binding]: https://github.com/tree-sitter/tree-sitter/tree/master/lib/binding_rust
[wasm binding]: https://github.com/tree-sitter/tree-sitter/tree/master/lib/binding_web

Other bindings are implemented in separate repositories, typically named
`$LANGUAGE-tree-sitter`:

- [Haskell binding](https://github.com/tree-sitter/haskell-tree-sitter/)
  ([documentation](https://hackage.haskell.org/package/tree-sitter Haskell binding documentation))
- [Python binding](https://github.com/tree-sitter/py-tree-sitter/)

### Language parser bindings

Complicating things even more, you need both the runtime library and the
generated parser for each language that you want to parse — and in particular,
you need _bindings_ for both!  The language bindings described above only
include the runtime library, since they can't know in advance which languages
you will want to parse.  The bindings should include instructions for how to
build and include your desired parsers.

For some language bindings, we can lean on the language's package manager for
this.  For instance, for the Rust bindings, we publish packages to crates.io
both for the language binding itself (the [`tree-sitter`][tree-sitter crate]
crate) and for most of the supported grammars (e.g. the
[`tree-sitter-python`][tree-sitter-python crate] crate).  So if you are writing
a tool, which is implemented in Rust, and which analyzes Python code, you would
add both `tree-sitter` and `tree-sitter-python` to your `Cargo.toml` file.
Wherever possible, we follow this approach for other language bindings, too.

[tree-sitter crate]: https://crates.io/crates/tree-sitter
[tree-sitter-python crate]: https://crates.io/crates/tree-sitter-python

You can also read this post via [Gemini][gemini].
{:.thanks}

[gemini]: gemini://dcreager.net/tree-sitter/map.gmi
