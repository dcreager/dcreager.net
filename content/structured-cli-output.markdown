---
kind: article
created_at: 2020-03-03
layout: /post.html
title: "Structured output from command line programs"
description: TBD
---

## Requirements

- Programs primarily consume data by reading from stdin, and produce it by
  writing to stdout (just like now).

- That data is _structured_, with an opinionated choice about which data
  serialization framework programs should use.

- The framework should provide a _schema language_ that lets you describe the
  particular structures that a program will produce and/or consume.

- Given a particular schema, the framework defines how data is _encoded_ into a
  stream of bytes (either in a file or flowing through a pipe).  Individual
  programs do not come up with custom encodings.

- When connecting two programs via a pipe, the schema of the producer and the
  schema of the consumer have to be _compatible_, but not necessarily
  _identical_ — the framework can automatically so some basic "massaging" to
  translate between compatible formats.

- You can interrogate a program to get its _preferred_ input and/or output
  schemas.

## The basics

[RFC 2119]: https://tools.ietf.org/html/rfc2119

### (1) Subcommands are distinct

If a program has "subcommands" that have entirely different uses (like how the
`git` program has `git status` and `git commit` subcommands), those subcommands
SHALL be considered separately.  In this document, "command" refers to all of
the distinct subcommands of a program.

### (2) Machine-focused commands only

If a particular command is only intended for human interaction, then these rules
SHALL NOT apply.

If a command has an argument that toggles between human-focused and
machine-focused use (such as `git status` vs `git status --porcelain`), then
these rules SHALL apply only when the command is operating in machine-focused
mode.

### (3) Use Avro

Any command that produces data that is intended for machine consumption, or
consumes data from a machine-produced source, MUST encode that data using
[Apache Avro][avro].

[avro]: https://avro.apache.org/

### (4) Preferred output schema

Any command that produces structured data MUST have a _preferred output schema_.

If you invoke the command with an additional `--avro-preferred-output-schema`
argument, the command MUST print out its preferred output schema to stdout, in
Avro's [JSON schema format][schema], and do nothing else.

[schema]: https://avro.apache.org/docs/current/spec.html#schemas

The command MAY produce different output depending on its arguments.  That means
that the caller MUST provide the same arguments that it would provide normally,
in addition to the `--avro-preferred-output-schema` argument.  The caller MUST
NOT assume that the preferred output schema for one list of arguments is the
same as the preferred output schema for another list of arguments.

### (5) Backwards-compatible changes

### (6) Requesting a different output schema

### (7) Preferred input schema

### (8) Binary and JSON encodings

### (9) Pretty-printing structured data happens elsewhere

## FAQ

### Why Avro?  Why not {Protobuf, Thrift, CBOR, JSON, etc.}?

- schema resolution is a first-class part of the spec

[schema-resolution]: https://avro.apache.org/docs/current/spec.html#Schema+Resolution
