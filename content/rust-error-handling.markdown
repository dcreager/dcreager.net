---
kind: article
created_at: 2020-12-16
layout: /post.html
title: Rust error handling patterns
tags: [rust]
---

This post summarizes how best to produce and consume errors in Rust code.  It’s short and to the point!  If you want more detail, check out [this great article][groenen] from Nick Groenen:

[groenen]: https://nick.groenen.me/posts/rust-error-handling/

## Applications

Much like how you handle Cargo.lock files, the recommendations are different for applications and libraries.  In an application, you’re consuming error conditions from a wide variety of libraries, each of which will have (hopefully) created custom library-specific error types.  The easiest way to consume all of them is to use the [anyhow][] crate.

[anyhow]: https://docs.rs/anyhow/

In particular, in your main.rs, _don’t_ do the bulk of your work directly in your `main` function.  Instead create a helper function (called `go` here) and do your work there:

<% highlight :rust do %>
use std::fs::File;
use std::path::Path;

use anyhow::Error;

fn go() -> Result<(), Error> {
    let path = Path::new("some-file.txt");
    let mut file = File::open(path)
        .with_context(|| format!("Error opening {}", path.display()))?;
    let mut content = String::new();
    file
        .read_to_string(&mut content)
        .with_context(|| format!("Error reading {}", path.display()))?;

    // Now you can do something with content!
    Ok(())
}
<% end %>

Note how you get to use the `?` operator to automatically propagate errors, and you can use `anyhow::Context::with_context` to annotate error messages with additional information about when and where they occur.  And this works with _any_ error type returned from the libraries they use, as long as they implement `std::error::Error`.

Your `main` function can then be very simple:

<% highlight :rust do %>
fn main() {
    if let Err(err) = go() {
        eprintln!("{:#}", err);
        std::process::exit(1);
    }
}
<% end %>

If an error occurs, we use anyhow’s `{:#}` display format to print out a nicely formatted error message, including any attached context.  (Make sure to print it to stderr!)  If no error occurs, we fall through to a normal exit.

## Libraries

So we know how to consume arbitrary errors in an application.  How about how to produce them in libraries?

The main thing is that you should create a custom error type for each library.  Possibly more than one!  Your error types will typically be enums, with a separate variant for each error condition that might occur:

<% highlight :rust do %>
pub enum FoozleError {
    DuplicateData,
    InvalidData,
    MissingData,
}
<% end %>

On its own, this is enough to let you use your error type inside of a `Result` value.  (`Result` itself doesn’t place any restrictions on its `E` parameter!)  But it’s not yet useful enough to work with the application recommendations in the previous section.  Ideally, all of your error types should implement `std::error::Error`, which in turn means that they need to implement `std::fmt::Display`.

Your `Display` implementation should provide a nice human-readable description of each error condition.  For those descriptions to be “nice” you probably want to include additional data about each error condition:

<% highlight :rust do %>
pub enum FoozleError {
    DuplicateData {
        name: String,
        old_value: String,
        new_value: String,
    },
    InvalidData {
        name: String,
        value: String,
    },
    MissingData {
        name: String,
    },
}
<% end %>

It can be tedious to implement `Display` manually, especially with this extra detail.  The [thiserror][] crate provides a derive macro that makes it much easier:

[thiserror]: https://docs.rs/thiserror/

<% highlight :rust do %>
#[derive(Error)]
pub enum FoozleError {
    #[error(
        "{} already exists (new value is {}, existing value is {})",
        .name,
        .new_value,
        .old_value,
    )]
    DuplicateData {
        name: String,
        old_value: String,
        new_value: String,
    },
    #[error(
        "{} is invalid (value is {})",
        .name,
        .value,
    )]
    InvalidData {
        name: String,
        value: String,
    },
    #[error(
        "{} doesn't exist",
        .name,
    )]
    MissingData {
        name: String,
    },
}
<% end %>

And that’s it!  The derive macro takes care of producing `Display` and `Error` implementations for you!
