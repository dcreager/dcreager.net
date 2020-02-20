---
kind: article
created_at: 2009-12-23
layout: /post.html
title: “High-water mark” buffers
tags: [c, libhwm]
---

My coding project for today was to extract out some code for dealing
with “high-water mark buffers”, putting it in a separate library call
`libhwm`.  In this post, I'm going to describe the rationale for using
them, and a brief overview of how to use the library.  (The library is
hosted on [Github](http://github.com/dcreager/libhwm/)).

By the way, this post (and the library) is all in C.

## What's all this then?

A common idiom I'm having to deal with these days is reading a really
large number of records from a data file.  We're talking well into the
millions of records, but we want the code to scale well past that.

### Step 1: Fixed-length records

Let's say that we need to read each record into a simple `struct`.
For now, we're going to use nice, fixed-length fields:

<% highlight :c do %>
typedef struct
{
    uint32_t  id;
    uint32_t  num_bananas;
} rec1_t;
<% end %>

With this datatype, we can actually read data from a file very
quickly; we'll just store each record directly in the file, in binary,
using 8 bytes.  (To simplify things, I'm not worrying about the
endianness of the integers, or whether the `struct` is packed; both
are easily handled with some pretty simple macro-fu.)

<% highlight :c do %>
FILE  *file = /* whatever */;
rec1_t  record;

while (fread(&record, sizeof(rec1_t), 1, file) == 1)
{
    /* process the record */
}
<% end %>

The C library's stream API (`fread` and friends) will buffer the data
from the actual file, so this gives us pretty good performance.

### Step 2: Variable-length records

What if we have a variable-length field in our `struct`, though?

<% highlight :c do %>
typedef struct
{
    uint32_t  id;
    uint32_t  num_bananas;
    char  *name;
} rec2_t;
<% end %>

Often in these cases, you can simplify the problem by deciding not to
let `name` be a variable-length field.  Instead, you decide that
you'll use (say) exactly 20 bytes for the name, padding out short
names and truncating long names as necessary.  We don't want to do
that, however — we want to have a truly variable-length field.

To store this variable-length field in the file, we need some way of
encoding the length of a particular record's `name` field.  If we can
assume that none of the records has a name that's longer than 4
billion characters, we can use a 32-bit length prefix:

<% highlight :c do %>
FILE  *file = /* whatever */;
rec2_t  record;

do
{
    uint32_t  name_length;

    if (fread(&record.id,
              sizeof(uint32_t), 1, file) < 1)
        break;
    if (fread(&record.num_bananas,
              sizeof(uint32_t), 1, file) < 1)
        break;
    if (fread(&name_length,
              sizeof(uint32_t), 1, file) < 1)
        break;
    if (fread(record.name,
              sizeof(char), name_length, file) < name_length)
        break;
    record.name[name_length] = '\0';

    /* process the record */
} while (true);
<% end %>

That's pretty ugly and repetitive, so let's play some macro games:

<% highlight :c do %>
FILE  *file = /* whatever */;
rec2_t  record;

#define READ_FIELD(dest, type, count) \
    if (fread(dest, sizeof(type), count, file) < count) \
        break;

do
{
    uint32_t  name_length;

    READ_FIELD(&record.id, uint32_t, 1);
    READ_FIELD(&record.num_bananas, uint32_t, 1);
    READ_FIELD(&name_length, uint32_t, 1);
    READ_FIELD(record.name, char, name_length);
    record.name[name_length] = '\0';

    /* process the record */
} while (true);
<% end %>

So the basic idea here is pretty sound — we can store a name of any
length without wasted space.  And the code is still rather fast; we'll
have a larger overhead from calling `fread` multiple times, but the
number of low-level I/O reads will still be roughly the same.

But unfortunately, there's a glaring error here.  This code will
segfault, since we haven't actually allocated any memory for the
`record.name` field.

### Step 3: Allocate some memory

So what's the simplest way we can allocate memory for the
`record.name` field?  The naïve approach would be to `malloc` a new
string for every record:

<% highlight :c do %>
FILE  *file = /* whatever */;
rec2_t  record;

#define READ_FIELD(dest, type, count) \
    if (fread(dest, sizeof(type), count, file) < count) \
        break;

do
{
    uint32_t  name_length;

    READ_FIELD(&record.id, uint32_t, 1);
    READ_FIELD(&record.num_bananas, uint32_t, 1);
    READ_FIELD(&name_length, uint32_t, 1);

    /* Remember to include an extra byte for the NUL terminator! */

    record.name = (char *) malloc(sizeof(char) * (name_length+1));
    if (record.name == NULL) break;

    READ_FIELD(record.name, char, name_length);
    record.name[name_length] = '\0';

    /* process the record */

    free(record.name);
} while (true);
<% end %>

This will avoid the segfault, and let you process your data, but it
will perform _horribly_, since we're calling down into the heap
management code for **every single record**!  And remember, we're
talking about millions of records here.

## Step 4: High-water mark buffers

So what's the solution?  A high-water mark buffer.  The idea is that
instead of allocating a new string each time through the loop, you
remember how large your current string is.  As long as the next
record's `name` isn't longer than your buffer, you can reuse it,
saving you a call to `malloc`.  If it is longer, you `realloc` it to
be large enough for the new string.  If you think of the lengths of
the `name` strings as a rising tide of water, you see where the name
of the buffer comes from.

We can do a high-water mark buffer by hand:

<% highlight :c do %>
FILE  *file = /* whatever */;
rec2_t  record;
size_t  allocated_name_size = 0;

#define READ_FIELD(dest, type, count) \
    if (fread(dest, sizeof(type), count, file) < count) \
        break;

record.name = NULL;

do
{
    uint32_t  name_length;
    size_t  name_size;

    READ_FIELD(&record.id, uint32_t, 1);
    READ_FIELD(&record.num_bananas, uint32_t, 1);
    READ_FIELD(&name_length, uint32_t, 1);

    /* Remember to include an extra byte for the NUL terminator! */

    name_size = sizeof(char) * (name_length+1);

    /* Reallocate the buffer if it's not big enough */
    if (name_size > allocated_name_size)
    {
        record.name = (char *) realloc(record.name, name_size);
        if (record.name == NULL) break;
        allocated_name_size = name_size;
    }

    READ_FIELD(record.name, char, name_length);
    record.name[name_length] = '\0';

    /* process the record */
} while (true);

if (record.name != NULL)
    free(record.name);
<% end %>

Note that `realloc` does the “right thing” if `record.name` is `NULL`;
this indicates that we haven't allocated a buffer yet, and so
`realloc` acts like `malloc` in this case.

## High-water mark library

So, we've described why you'd want a high-water mark buffer, and how
to implement one.  But once you write that same code three or four
times, you decide to factor it out into a library.  Hence
[libhwm](http://github.com/dcreager/libhwm/).  Here's the same file
reading code using the library:

<% highlight :c do %>
typedef struct
{
    uint32_t  id;
    uint32_t  num_bananas;
    hwm_buffer_t  name;
} rec3_t;
<% end %>

<% highlight :c do %>
FILE  *file = /* whatever */;
rec3_t  record;

#define READ_FIELD(dest, type, count) \
    if (fread(dest, sizeof(type), count, file) < count) \
        break;

hwm_buffer_init(&record.name);

do
{
    uint32_t  name_length;
    char  *name_ptr;

    READ_FIELD(&record.id, uint32_t, 1);
    READ_FIELD(&record.num_bananas, uint32_t, 1);
    READ_FIELD(&name_length, uint32_t, 1);

    /* Remember to include an extra byte for the NUL terminator! */

    name_size = sizeof(char) * (name_length+1);

    /* Read into the HWM buffer */
    if (!hwm_buffer_ensure_size(&record.name, name_size))
        break;

    name_ptr = hwm_buffer_writable_mem(&record.name, char);
    READ_FIELD(name_ptr, char, name_length);
    name_ptr[name_length] = '\0';

    /* process the record */
} while (true);

hwm_buffer_done(&record.name);
<% end %>

Et voila.  Of course, this last code snippet makes me realize that we
could make things even simpler with an `hwm_buffer_fread` function!
The story never ends...
