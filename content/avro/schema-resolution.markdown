---
kind: article
created_at: 2020-03-12
layout: /post.html
title: "Avro schema resolution"
description: >-
  Happy the people whose annals are tiresome
tags: [avro]
---

I have always liked [Avro][] better than the other data serialization frameworks
that are out there ([Protobuf][], [Thrift][], [JSON][], [msgpack][], [CBOR][],
the list goes on and on).  There are many reasons why, which probably deserves
another post.  But a big one is Avro's [schema resolution][] rules.  Let's take
a look at how they work and why they're useful!

[Avro]: https://avro.apache.org/
[Protobuf]: https://developers.google.com/protocol-buffers
[Thrift]: https://thrift.apache.org/
[JSON]: https://json.org/
[msgpack]: https://msgpack.org/
[CBOR]: https://cbor.io/
[schema resolution]: https://avro.apache.org/docs/current/spec.html#Schema+Resolution

<hr class="jump">

One of the main problems that data serialization frameworks try to solve is the
**_backwards compatibility_** problem.  Your software will naturally evolve over
time, as you fix bugs and add features.  And as part of that evolution, the data
that your program writes or reads will change, too.

So as an example, let's say we're creating a program that deals with lists of
countries.  It will need a data type to describe what information it needs to
know about each country, and how to store that information internally.  Let's
use Pascal, because...well...why not!

<% highlight :pascal do %>
record
  IsoCode: string[2];
  Name: string[60];
end
<% end %>

The program needs to read in a list of countries and then do...something...to
it.  (The "something" isn't really relevant to us today.)  We're writing this
in, say, 1987, so XML and JSON don't exist, let alone any of the other
frameworks we've mentioned.  We have to roll our own format!  So we go for
simplicity, and say that you provide a text file, with:

- each country on a separate line
- the ISO code is the first two characters of the line
- skip any whitespace immediately following
- use the remainder of the line as the country name

<pre>
US   United States
CH   China
KZ   Kazakhstan
UK   United Kingdom
TV   Tuvalu
MX   Mexico
EG   Egypt
BZ   Belize
BR   Brazil
</pre>

(Note that's it's okay for the country name to include whitespace!  Our file
format definition doesn't make that ambiguous.)

Back in the old days, when
we would come up with completely custom file formats for each program (text or
binary, doesn't matter), you had to take this into account explicitly.  If you
just wrote out a bunch of fields, in order, without anything 


FAST

A properly tuned Avro library can _rip_ through large files, even if it has to
transform the data by dropping or reordering fields as they're being processed.
How?  Because your program figures out _how_ to do that transformation for _all_
of the records in a file _once_, up front, based on the writer and reader
schemas.  Once it's done that, it can just blindly follow the transformation
instructions for each record in the file.

By constrast, most other frameworks "self-describe" their data by including
extra information _each in value_.  This is most obvious in JSON (and its binary
renderings, like msgpack and CBOR), where the field names are included in every
object.  It's also true in protobufs, where each field in a message has a [field
number][protobuf field number] that is included in the message's serialization.
Protobuf message fields can appear in any order (and for `optional` fields, can
be left out completely), which means that you have to figure out _for each
individual value_ which fields are present, and how they should be deserialized
and transformed into your internal representation.

[protobuf field number]: https://developers.google.com/protocol-buffers/docs/proto#assigning-field-numbers

[dphil]: /publications/012-dphil-thesis/
[dcreager avro resolution]: https://github.com/apache/avro/blob/791ec6017b03e3cc79fbf546dd9d79d20d476b4e/lang/c/src/resolved-writer.c#L2676-L2680
