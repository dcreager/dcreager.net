---
layout: post
title: Using callbacks with the subprocess module
---

# {{ page.title }}

In a [previous post](/2009/08/06/subprocess-communicate-drawbacks/),
we pointed out two drawbacks of Python's `subprocess.communicate`
method.  In this post, we look at the first problem in more detail.
To reiterate, the problem is that we collect the subprocess's output
streams into strings.  If the subprocess is going to generate a huge
amount of output, it can be better to process the output data in a
stream-oriented manner — that way we use a constant amount of memory
regardless of how much output is produced.

If we look at the
[implementation](http://svn.python.org/view/python/trunk/Lib/subprocess.py?revision=74029&view=markup)
of the `communicate` method, we see this:

{% highlight python linenos %}
def _communicate_with_select(self, input):
    read_set = []
    write_set = []
    stdout = None # Return
    stderr = None # Return

    if self.stdin and input:
        write_set.append(self.stdin)
    if self.stdout:
        read_set.append(self.stdout)
        stdout = []
    if self.stderr:
        read_set.append(self.stderr)
        stderr = []

    input_offset = 0
    while read_set or write_set:
        try:
            rlist, wlist, xlist = select.select(read_set, write_set, [])
        except select.error, e:
            if e.args[0] == errno.EINTR:
                continue
            raise

        if self.stdin in wlist:
            chunk = input[input_offset : input_offset + _PIPE_BUF]
            bytes_written = os.write(self.stdin.fileno(), chunk)
            input_offset += bytes_written
            if input_offset >= len(input):
                self.stdin.close()
                write_set.remove(self.stdin)

        if self.stdout in rlist:
            data = os.read(self.stdout.fileno(), 1024)
            if data == "":
                self.stdout.close()
                read_set.remove(self.stdout)
            stdout.append(data)

        if self.stderr in rlist:
            data = os.read(self.stderr.fileno(), 1024)
            if data == "":
                self.stderr.close()
                read_set.remove(self.stderr)
            stderr.append(data)

    return (stdout, stderr)
{% endhighlight %}

(There are actually several different `communicate` implementations in
the module: a Windows-specific implementation, an implementation using
the POSIX `poll` function, and one using POSIX `select`.  We're going
to look at the `select` implementation; the modifications we make can
be rolled into the other methods, too.)

## Output callbacks

For collecting stdout, the important part is lines 33-38.  If the
`select` call tells us that the stdout stream is ready for reading, we
try to read up to 1024 bytes from it.  If we get the empty string,
this means we've reached EOF, and can close down the stream.  (We also
no longer have to keep passing it in to further `select` calls, since
we know we're done with this stream.)  If we get a non-empty string,
then we append it into a list.  The function that calls
`_communicate_with_select` will eventually `join` this list of strings
together, yielding a single string.

It's actually a very simple change to make this process the output
using a stream-based callback.  For now, we assume that we've been
given the callback, in a `stdout_callback` variable.  Then, we can
change line 38 to

{% highlight python %}
            stdout_callback(data)
{% endhighlight %}

The callback can be any Python callable object; it should accept a
single argument, which is the next chunk of data from stdout.  We can
make a similar change to line 45 to send the stderr data to its own
`stderr_callback`.

## Line-oriented callbacks

One possible issue with the output callbacks in the previous section
is that the data is sent into the callback in arbitrary chunks.  We
might prefer to guarantee that the callback will be called exactly
once for each *line* of output.  This would allow us, for intsance, to
easily interleave the output lines of a bunch of subprocesses into the
output of the parent process, without having to worry about locking.

To use a line-based callback, we have to wrap it, creating a
arbitrary-chunk callback that buffers data as necessary.  Once it
receives a full line of data, it sends it to the wrapped line
callback.

{% highlight python linenos %}
def wrap_line_callback(line_callback):
    class Callback(object):
        def __init__(self, line_callback):
            self.line_buffer = []
            self.line_callback = line_callback

        def output_buffer(self):
            line = "".join(self.line_buffer)
            self.line_callback(line)
            self.line_buffer = []

        def data_callback(self, data):
            # If we get an empty string, that represents the end of
            # the input.  If there is anything in the buffer, send
            # then out first, then send an empty string on the to line
            # callback.

            if data == "":
                if len(self.line_buffer) > 0:
                    self.output_buffer()

                self.line_callback("")
                return

            # Otherwise, we split the new data into separate lines,
            # each of which we call an “entry”.  We add each entry to
            # the buffer.  If the entry ends with a newline, we output
            # the buffer and then clear it.

            lines = data.splitlines(True)
            for line in lines:
                self.line_buffer.append(line)
                if line.endswith(("\r\n", "\n", "\r")):
                    self.output_buffer()

            return

    return Callback(line_callback).data_callback
{% endhighlight %}

* **line 1** — We start by declaring the wrapping function.  It will
  take in a line-based callback, and return an arbitrary-chunk
  callback.

* **lines 2-5** — We'll need to maintain some state in between
  invocations of the arbitrary-chunk callback — specifically, if a
  line of output data falls across a chunk boundary, we'll need to
  hold onto the part in the first chunk until we receive the part in
  the second chunk.  One relatively easy way to do this is to create a
  new class, and store the state in `self` properties.  Note that the
  `Callback` class is declared inside the `wrap_line_callback`
  function.

  The buffer is a list of strings.  This needs to be a list to support
  arbitrarily long lines, since we don't know how many chunks we'll
  need to store before outputting the line.

* **lines 7-10** — We also declare a method in the class that can
  output the current contents of the saved buffer (if any).  Like the
  original `communicate` method, we join the buffer list together into
  a single string, then pass it onto the wrapped line callback.

* **lines 12-23** — Next we can define the arbitrary-chunk callback.
  First, it checks to see if we've received an EOF indicator (the
  empty string).  If so, we first output whatever is currently in the
  buffer; then we pass on the EOF to the line callback.

* **lines 25-34** — If we get some actual data, we first split the
  current chunk into separate lines.  We pass in a `True` parameter to
  the `splitlines` method so that we get the newlines in the split
  strings.  This will let us tell if the final string in the list
  represents a complete line, or if it's the first part of a line that
  falls across a chunk boundary.

  We then loop through the split strings, adding them to the buffer.
  If any of them end with one of the newline endings (Unix, Windows,
  or otherwise), we output the buffer.  Note that only the last string
  in the `lines` list can end with something other than a newline;
  anything at the beginning of the list is only a separate element
  because `splitlines` found a newline to split on!

* **line 38** — Finally, we instantiate the new class and return its
  arbitrary-chunk callback.

## More to come

This gives us a nice output callback mechanism, to prevent us from
buffering the entire contents of the stdout and stderr streams in
memory.  In later posts, we'll look at doing the same with the *input*
data, and then we'll address the second problem, of dealing with more
than one subprocess simultaneously.
