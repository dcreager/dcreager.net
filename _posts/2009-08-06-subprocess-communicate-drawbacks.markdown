---
layout: post
title: Pros and cons of Python's <code>subprocess.communicate</code> method
---

# {{ page.title }}

The `subprocess` module, which was introduced in Python 2.4, provides
you with a convenient interface for spawning *subprocesses*, and for
interacting with these subprocesses in your parent process.  The
module was introduced in [PEP
324](http://www.python.org/dev/peps/pep-0324/), and is a replacement
for the proliferation of other functions and modules that were used
previously for spawning and interacting with processes.  The
`subprocess` module aims to provide a more consistent interface,
regardless of the particulars of how you need to interact with the
subprocesses.

## Overview of the `subprocess` module

Subprocesses are encapsulated in a `Popen` object.  You interact with
a subprocess via its stdin, stdout, and stderr streams.  When you
create a new `Popen` object, you can give a value of `PIPE` for the
`stdin`, `stdout`, and `stderr` keyword parameters.  If you do, then
the `Popen` object that you get back will have `stdin`, `stdout`,
and/or `stderr` attributes.  Each of these is a file-like object,
giving you access to the corresponding stream of the subprocess.

Now, you have to be careful how you use these pipe objects, since it's
easy to fall into a situation where you have deadlock.  For instance,
your parent process might be trying to write some data into the
`stdin` pipe, to send some information into the subprocess.  The
subprocess, on the other hand, is trying to write some data into the
`stdout` pipe, to send some infomration back out to the parent
process.  If the `stdout` pipe's buffer is full, then the subprocess
will block trying write into the pipe; it won't be able to proceed
until the parent process has read some data from the `stdout` pipe,
clearing room in the buffer for the new data.  However, the parent
process is currently trying to write into the `stdin` pipe.  If this
write is also blocked, then we have deadlock — neither process can
proceed.

## The `communicate` method

The usual solution in these cases is to use the `Popen` object's
`communicate` method.  This method takes in an optional string to send
to the subprocess on stdin.  It then collects all of the stdout and
stderr output from the subprocess, and returns these.  The
`communicate` method takes responsibility for avoiding deadlock; it
only sends the next chunk of the stdin string when the subprocess is
ready to read it, and it only tries to read the next chuck of stdout
or stderr when the subprocess is ready to provide it.

Under the covers, the `communicate` method uses a `select` loop to
perform this choreography with the subprocess.  (At least for the Unix
implementation of the `subprocess` module, that is.)  This solution is
nice because it doesn't require introducing threading into the parent
process.  During each iteration of the loop, it calls the OS's
`select` system call, giving it the file descriptors of the stdin,
stdout, and stderr pipes.  The `select` call tells us which of these
file descriptors can perform an I/O operation without blocking.  If
none of them can immediately, it will block until one of them can.
Once the `select` call returns, we read from or write to the pipes
that are ready.  We repeat this process until we see EOF on both
stdout and stderr; this indicates that the subprocess has finished —
or at least, that it's through communicating with us.

## Drawbacks

The `communicate` method provides a nice, simple interface for
interacting with a subprocess, without having to worry about deadlock
situations.  Unfortunately, it has two main drawbacks:

* The subprocess's stdout and stderr are collected into strings.
* You can only interact with one subprocess at a time.

(If neither of these is an issue for you, then the rest of this post
is less interesting to you — `communicate` does exactly what you
want!)

The first item is a problem if your subprocess creates a lot of output
— the worry is the output will be too large to fit into a Python
string.  If it is, then the parent process will (at best) start to
thrash as it eats into virtual memory.

The second item is a problem if you have to spawn multiple
subprocesses, and interact with them simultaneously.  You could argue
that there's no need to fix this problem if you haven't fixed the
first: since the `communicate` method is just going to collec the
stdout and stderr into strings, then you could just loop through each
of your subprocesses, calling `communicate` on each in turn.  The end
result would be what you want — all of the stdout and stderr strings
for all of your subprocesses.

However, do so can make your subprocesses take longer to run, since
you won't be able to exploit parallelism as much.  Since you're firing
off these subprocesses at the same time, you supposedly want them to
execute simultaneously, allowing the OS to schedule them appropriate
so that they finish as quickly as possible.  However, you've
introduced a serialization into this logic, since your parent process
is only able to interact with one subprocess at a time.  For instance,
subprocess #2 might be waiting for some input, while the parent
process is still snarfing up the output from subprocess #1.  In this
case, subprocess #2 is ***not going to be able to start executing***
until subprocess #1 has ***completely finished***.  So your
`communicate` loop has completely eliminated the benefit of starting
the subprocesses simultaneously.

In later posts, I will outline how to solve these two problems.
