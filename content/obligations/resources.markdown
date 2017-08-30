---
kind: article
created_at: 2017-08-30
layout: /post.html
series: "Obligations"
title: "Resources"
description: >-
  ...
tags: [obligations]
---

We want to talk about **requirements** and **obligations** in our programming
languages, and how they differ from **permissions** or **possibilities**.  To
start the discussion, we should think about what kinds of requirements we might
want to support.  That will help ground the discussion on how we can implement
those requirements.

The most common example of a requirement in a program involves **resource
management**.  Your program acquires some resource, and we want to ensure that
you clean up that resource when you're done with it.  C is notorious for being
an "unsafe" language because it provides basically no mechanisms for making
these kinds of guarantees.  To be clear, this does *not* mean that it's
impossible to write a safe program in C; it just means that it's *entirely* on
you, the programmer, to implement whatever safety features you need, and to do
all of the mental bookkeeping to make sure you didn't miss anything.  There's
nothing in C that prevents you from doing either of the following:

<% highlight :c do %>
void not_very_nice(void) {
    FILE* file = fopen("/usr/share/words", "r");
    return;
}

void also_not_nice(void) {
    int* i = malloc(sizeof(int));
    return;
}
<% end %>

So what do other languages do to handle these kinds of situations?

##### Automatic memory management

The most common resource that our programs deal with is memory, so it's no
surprise that we've put the most effort into solving this problem for memory
management.

- stack allocation

ensures that an object is deallocated; typically doesn't prevent dangling pointer
references

- garbage collection

a lot of engineering effort goes into the GC; typically the largest part of an
implementation of a VM or language runtime; usually still worth it because of
the time and bugs saved for the end user programmer

This example is a bit of a special case; memory management is so important to so
many programs that most modern languages have baked memory safety *into the
language itself*.  But what if a library author is trying to create a new kind
of resource that other people can use in their programs?  We can't update the
language itself (and all of its compilers or interpreters) to support our new
kind of resource.  What options do we have *in our programs* for defining what
the clean-up behavior should be, and enforcing that it's executed?

##### Resources as objects

- destructors / RAII

the language ensures that the destructor is called just before an object is
deallocated; represent the resource with an object, and put the clean-up logic
in its destructor

    std::ifstream

stack allocation; file is closed
new but no delete; file is not closed!

    std::unique_ptr / std::shared_ptr

modern C++ style is to never call `new` directly

[unique_ptr]: http://en.cppreference.com/w/cpp/memory/unique_ptr
[ifstream]: http://en.cppreference.com/w/cpp/io/basic_ifstream/basic_ifstream
[Google style]: https://google.github.io/styleguide/cppguide.html
[smart pointers style]: https://google.github.io/styleguide/cppguide.html#Ownership_and_Smart_Pointers

- finalizers are different

above only works when you know when an object is going to be destroyed.

The equivalent in a garbage-collected language is called a **finalizer**.

In a garbage collected language, you don't have guarantees about *when* an
object will be collected, just that it eventually *will*.  That can be a problem
when you want to make sure that you clean up your resource as soon as you're
done using it!

##### Context managers

[Python with]: http://www.python.org/peps/pep-0343.html

##### Stack-like scoping

context manager only supports stack-like scoping

C++ move semantics let you do other things

all of this only works for very simple "clean up after yourself" requirements.
more complex patterns aren't supported
