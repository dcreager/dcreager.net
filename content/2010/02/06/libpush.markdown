---
kind: article
created_at: 2010-02-06
layout: post
title: A combinator-based parsing library for C
tags: [libpush, c]
---

Recently I've been working on
[libpush](http://github.com/dcreager/libpush/), which a new parsing
library for C.  It has two main features that I think will be
valuable: it's a _push parser_, which means that instead of parsing a
file, stream, or single memory buffer, you supply the data (or "push"
it) to the parser in chunks, as it becomes available.  I plan to
discuss this aspect of the parser in more detail in a later post.

The other main feature is that you design your parsers using
_combinators_.  Parser combinators are widely used in Haskell, with
[Parsec](http://legacy.cs.uu.nl/daan/parsec.html) being the most
common example.  Combinator-based parsing libraries are especially
nice in Haskell, because Haskell's syntax makes them look very simple.
For instance, a parser that parses matching nested parentheses is:

<% highlight :haskell do %>
parens :: Parser ()
parens = (char '(' >> parens >> char ')' >> parens) <|> return ()
<% end %>

Here, the `<|>` operator represents _choice_: we try parsing the left
operand, and if it fails, then we try the right operand.  In our
example, the right operand is the base case, which matches the empty
string.  The left operand parses an opening parenthesis; then
recursively calls itself to match any parentheses that might be nested
in the current set; then parses the closing parenthesis; and then
finally tries to match a nested set that occurs after the current set.

When we say that this is a combinator-based parser, we mean that it's
implemented by taking _primitive parsers_ — in this case `char '('`
and `return ()` — and combining them into more complex parsers using
generic operators like `>>` and `<|>`.

Now, in order to be able to use combinators like this, parsers have to
be first-class objects in your language.  In the Haskell code, the
parsers are represented by the `Parser ()` type.  In most Haskell
parsing libraries (including Parsec), the parser type is implemented
as a
[_monad_](http://en.wikipedia.org/wiki/Monad_%28functional_programming%29).
Monads have a reputation for being a horribly complex topic, but in
this case, we don't really need to learn about the underlying math.
Instead, we can just view the monad as letting us do two things
concisely:

 1. Parsers can return a value, which could (for instance) be the
    abstract syntax tree that you're building up while parsing your
    language.  The monadic bind operator (`>>=`) gives you a way to
    "pass" these values between parsers, if needed.

 2. Simultaneously, the parser monad maintains the state of the stream
    you're parsing from, keeping track of how many bytes remain,
    whether there's an error condition, and possibly a nice
    human-readable description (line and column) of the current
    location.

This is admittedly a lot of setup; we've been talking a lot about
Haskell in a post that's ostensibly describing a C library.  But
hopefully, this gives you a taste for the kinds of features we want to
support in libpush:

 * Parsers will be represented by a C type.  In libpush, this is the
   `push_callback_t` type.

 * There will be several primitive parsers; these will be functions
   that return a `push_callback_t`.  The functions can take in
   parameters, but none of the parameters will be a `push_callback_t`.
   (See the `char` primitive from above; it needed to take in the
   particular character that is expected.)

 * There will be several combinators; these will be functions that
   return a `push_callback_t`, and take in other `push_callback_t`s as
   parameters.

   You can see several of these primitives and combinators in action
   in the [libpush Github
   repository](http://github.com/dcreager/libpush/).

 * We will use something like a monad to take care of passing values
   between our parsers, and for keeping track of the state of the
   underlying stream.  I say "something like a monad", because, unlike
   the Parsec library, the libpush parser type will _not_ be
   implemented as a monad; in turns out that C is more amenable to
   implementing them as [_arrows_](http://www.haskell.org/arrows/).
   In a later post, I'll explain what this means in terms of writing
   your own parsers, or for building them up from combinators.
