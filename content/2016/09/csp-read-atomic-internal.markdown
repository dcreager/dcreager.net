---
kind: article
created_at: 2016-09-07
layout: post
series: "Concurrency models in CSP"
title: "Read Atomic: Internal consistency"
tags: [csp]
---

{::options parse_block_html="true" /}

<script type="text/x-mathjax-config">
MathJax.Hub.Config({
  TeX: {
    Macros: {
      interleavesym: "|\\mkern-2mu|\\mkern-2mu|",
      interleave: "\\mathrel{\\interleavesym}",
      Interleave: "\\mathop{\\interleavesym}"
    }
  }
});
</script>

We'll start by looking at the weakest concurrency model covered in the paper,
Read Atomic.  And in fact, to keep things really simple to start with, we're
only going to look at *one* of Read Atomic's two axioms: **internal
consistency**.  (We'll look at external consistency in the next post.)

<!--
 Read Atomic...can be implemented without requiring any coordination among
 replicas...a replica can decide to commit a transaction without consulting
 other replicas.
-->

A transaction is internally consistent if it "reads its own writes".  This is
the simplest axiom covered in the paper, since it expresses a property that's
strictly *local* to each transaction; we don't have to consider the behavior of
any other transaction when defining what it means for a transaction to be
internally consistent.

#### Types and channels

Our goal is to construct a CSP process that describes internal consistency.  We
start by defining some data types and channels.

<div class="aside-def">
## Channels
Each \\(\textbf{channel}\\) statement defines one kind of event that your
processes can use to communicate.  Each event will typically carry some data;
you use \\(\textbf{datatype}\\) and \\(\textbf{nametype}\\) statements to define
the different kinds of data that can be carried by those events.
</div>

We'll first need to talk about the objects in the data store, and the values
that those objects can have:

{::comment}
nametype Object = {1..2}
nametype Value = {1..5}
{:/comment}

\\[
\textbf{nametype} ~ \textsf{Object} = \\{1..2\\} \\\
\textbf{nametype} ~ \textsf{Value} = \\{1..5\\}
\\]

An \\(\textsf{Object}\\) is the name of one of the objects that can be read or
written.  The \\(\textsf{Value}\\) type is the set of values that can be stored
in each object.  We keep both of these sets small, to make it easier for FDR to
perform refinement checks without needing too much memory.  (In particular, we
can't use \\(\mathbb{Z}\\), the infinite set of integers, like they do in the
original paper, since that *definitely* won't fit into memory!)

{::comment}
It corresponds to the \\(Obj\\) set in the original paper.  In the original
paper, they use \\(\mathbb{Z}\\) (the set of integers); in our CSP spec, we'll
also use simple integers, but we'll limit ourselves to a finite set.
{:/comment}

Once we have our data types defined, we can define our channels:

{::comment}
channel read : Object.Value
channel write : Object.Value
{:/comment}

\\[
\textbf{channel} ~ \texttt{read} : \textsf{Object} \times \textsf{Value} \\\
\textbf{channel} ~ \texttt{write} : \textsf{Object} \times \textsf{Value}
\\]

A \\(\texttt{read}\\) event tells us that a particular value was read from a
particular object, while a \\(\texttt{write}\\) event signifies that a
particular value was written to a particular object.  For instance,
\\(\texttt{read}{.}1{.}3\\) tells us that a transaction read the value 3 from
object 1, and \\(\texttt{write}{.}2{.}4\\) tells us that a transaction wrote the
value 4 to object 2.

<div class="aside-def">
## Timing
Note that we don't have to include timestamps or anything like that in our
channel descriptions to indicate that one event occurred before another.  In
CSP, we get that for free from the operators that we use to construct a process.
For instance, the "prefix" operator (\\(\rightarrow\\)) lets us say something
like:

\\[
\texttt{read}{.}1{.}3 \rightarrow \texttt{write}{.}2{.}4 \rightarrow
\ldots
\\]

meaning that a transaction read the value 3 from the object 1, *and then* wrote
the value 4 to the object 2 (and then does something else that we don't care
about right now).
</div>

{::comment}
Based on the structure of these two CSP channels, you might expect that they
correspond to the \\(Op\\) set from the original paper; however, they actually
correspond to the \\(REvent\_x\\) and \\(WEvent\_x\\) family of sets.  In the
original paper's formalism, the authors needed to include an \\(EventId\\) in
each event to describe the ordering between events.  In CSP, on the other hand,
the prefix operator (\\(\rightarrow\\)) lets us define how the events in a
process are ordered without having to add extra fields to the events.
{:/comment}

#### Undefined objects

Now we can finally start defining some processes!  Within a transaction, each
object starts **undefined**, meaning that *this transaction* has not yet written
a value to the object.  (Remember that internal consistency does not take into
account what values other transactions write into an object.)  We can construct
a process (parameterized by an \\(\textit{obj}\\) variable) that remembers that
*one particular object* is currently undefined:

{::comment}
Undefined(obj) =
  write!obj?value -> Defined(obj, value)
    []
  read!obj?value -> Undefined(obj)
{:/comment}

\\[
\textrm{Undefined}(\textit{obj}) = {} \\\
  \quad\texttt{write}\,!\textit{obj}\,?\textit{value} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{value}) \\\
  \qquad \Box \\\
  \quad\texttt{read}\,!\textit{obj}\,?\textit{value} \rightarrow
        \textrm{Undefined}(\textit{obj})
\\]

The first clause tells us that you can write to an undefined object.  The
\\(!\textit{obj}\\) part ensures that we only worry about \\(\texttt{write}\\)
events for this object; the \\(?\textit{value}\\) part means that we don't want
to place any restrictions on what value you can write into the object.  If you
write into an undefined object, the object becomes defined (which we'll consider
below), and we remember the value that was written.

<div class="aside-def">
## \\(\rightarrow\\) (prefix), \\(?\\) and \\(!\\)
The prefix operator lets us introduce sequencing; \\(a \rightarrow B\\) means
that the process performs an event \\(a\\), and then it whatever is described by
the process \\(B\\).  Nice and simple!  The only wrinkle is that the left side
of the arrow is an *event*, while the right side is a *process*.

You will usually see the special \\(?\\) and \\(!\\) event constructors used
with the prefix operator.  This is just some syntactic sugar that means that we
don't have to explicitly write out separate clauses for
\\(\texttt{read}{.}1{.}1\\), \\(\texttt{read}{.}1{.}2\\),
\\(\texttt{read}{.}2{.}1\\), \\(\texttt{read}{.}2{.}2\\), and so on.

The \\(?\\) constructor lets you create a new variable which can take on *any*
of the possible values that could appear in that "slot" of the event; you can
refer to the new variable in the rest of the clause if you need to know which
particular value was selected.

The \\(!\\) constructor lets you assert that only one specific value is allowed;
that value is given by the variable that follows the \\(!\\).  (That means that
\\(?\\) should be followed by a *new* variable name, but \\(!\\) should be
followed by an *existing* variable name.)
</div>

The second clause tells us that you can also *read* from an undefined value.
This might be surprising!  Internal consistency doesn't tell you which value
you'll read, so why should we include this case?  This is one of the first
hurdles you'll encounter when learning about formal methods: in the world of
math, there's a very big difference between *disallowing* something and *not
caring*.  Here, we don't care: just because a value hasn't been defined by this
transaction doesn't mean that some previous transaction hasn't committed a
perfectly good value for it; in fact, all of the later concurrency models will
be interesting exactly because of all of the different ways they decide which
value you see!  We need to include the \\(\texttt{read}\\) clause, and use the
\\(?\textit{value}\\) syntax to not place any constraints on what value is read,
to leave open all of those different interesting possibilities.

<div class="aside-def">
## The environment
When describing the events in a process, we always keep in mind that the process
is going to communicate or interact with something else, called its
**environment**.  CSP uses a "rendezvous model", which means that both sides of
an interaction must agree before an event can take place.  Typically, that
"something else" is another process that we've combined together using one of
CSP's composition operators.  But we can also talk about an "observer" — such as
you, the specification author — who wants to explore the behavior of the process
by walking through all of its possible combinations of events.
</div>

<div class="aside-def">
## Alphabets
Each process has an **alphabet**, which is the set of events that it cares
about.  If an event is not in a process's alphabet, the process cannot perform
that event, and it cannot prevent any other processes from performing the event.
In fact, the process isn't even allowed to *mention* the event.

On the other hand, if an event *is* in a process's alphabet, then the process
controls whether and when that event occurs — or doesn't!  If the environment
wants to perform an event in the process's alphabet, but the process isn't ready
to perform that event, then the environment is *prohibited* from performing the
event!

With great power comes great responsibility!  You must be careful not to
**overspecify** your processes; it can be easy to accidentally *prevent* an
event from occurring, when you actually want to express that you *don't care*
whether it occurs.  Ideally you do that by not including the event in your
alphabet.  But you might need the event in the your alphabet so that you can
control it at some other time; you just don't want to *right now*.  In that
case, you have to include a "fallthrough" clause like we did for reading
undefined objects.
</div>

The two clauses are combined using external choice (\\(\Box\\)), which means
that whenever it's time for the transaction for perform another operation, it's
free to perform *either* a read or a write.

<div class="aside-def">
## \\(\Box\\) (external choice)
The external choice operator combines two processes, giving the environment full
control over which "branch" occurs.  For \\(A \mathrel{\Box} B\\), we're willing
to act like process \\(A\\) or process \\(B\\), and don't care (and have no way
to influence) which one occurs; the environment decides.
</div>

#### Defined objects

Once we've written a value to an object (in this transaction), the object
becomes **defined**.  We use a parameterized process for defined objects, just
like we did for undefined objects, but we need an extra
\\(\textit{currentValue}\\) parameter to keep track of the object's current
value.

{::comment}
Defined(obj, currentValue) =
  write!obj?newValue -> Defined(obj, newValue)
    []
  read!obj!currentValue -> Defined(obj, currentValue)
{:/comment}

\\[
\textrm{Defined}(\textit{obj}, \textit{currentValue}) = {} \\\
  \quad\texttt{write}\,!\textit{obj}\,?\textit{newValue} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{newValue}) \\\
  \qquad \Box \\\
  \quad\texttt{read}\,!\textit{obj}\,!\textit{currentValue} \rightarrow
        \textrm{Defined}(\textit{obj}, \textit{currentValue})
\\]

The first clause tells us that you can write to a defined object (just like you
can write to an undefined one).  Doing so should overwrite the current value
with the new value; we handle this by using the new value as the parameter in
our recursive call.

The second clause tells us that you can read from a defined object.  Unlike the
undefined case, internal consistency constrains what value you read when you do
so.  The \\(!\\textit{currentValue}\\) part enforces that the value returned is
the current value of the object, not allowing any other
\\(\texttt{read}{.}\textit{obj}\\) events to occur at this point.  Because of
CSP's rendezvous model, even if another process *tries* to read some other
value, the fact that we've actively prevented it means that they won't be able
to.

#### Combining it together

We now have subprocesses telling us how internal consistency behaves, on a
per-object basis, for both undefined and defined objects.  The last thing we
need to do is combine all of the per-object subprocesses together into a single
process describing internal consistency as a whole.

<div class="aside-def">
## Parallel composition
The most important operators in CSP are the *parallel composition* operators,
which combine two or more processes together and let them run concurrently.
There are several different flavors of parallel composition, which we'll see as
we work through the concurrency models, but they all have the same basic
"shape".

Each of the processes that you're combining are allowed to execute concurrently.
There is a set of events (called the "interface") that the processes can use to
communicate with each other.  For events outside the interface, the processes
are completely independent: they can proceed as quickly or as slowly as they
want, and cannot interfere which each other in any way.  If two processes both
want to perform the same non-interface event, they can.  However, those events
are also completely independent: the environment will see both "instances" of
the event occur, and can't tell which one belonged to which process.

Events in the interface, however, are **synchronous**: *all* of the processes
being combined must be ready to perform the event before *any* of them are
allowed to.  Once all processes are ready, the event can occur.  If the event
occurs, it's shared; the environment will only see one "instance" of the event,
and all of the processes will *simultaneously* "handle" the event and move on to
the their next state.
</div>

{::comment}
InternalConsistency =
    ||| obj : Object @ Undefined(obj)
{:/comment}

\\[
\textrm{InternalConsistency} =
  \textstyle\Interleave_{obj} \textrm{Undefined}(obj)
\\]

The \\(\Interleave_{obj}\\) notation is kind of like the \\(\sum\\) operator
that you might know from plain old arithmetic; it means that we're going to
instantiate the \\(\textrm{Undefined}\\) subprocess once for each possible
\\(\textit{obj}\\), and then combine them all together using the "interleave"
operator (\\(\interleave\\)).  Each object starts undefined, so we only have to
mention the \\(\textrm{Undefined}\\) subprocess; each subprocess will
automatically transition to \\(\textrm{Defined}\\) at the right time, if the
corresponding object is ever written to.

<div class="aside-def">
## \\(\interleave\\) (interleaving)

The simplest parallel composition operator is "interleaving", where the
interface is empty.  All of the processes execute completely independently, and
cannot communicate with each other, even if their alphabets have any events in
common.
</div>

We can use interleave here because the alphabet of each of our subprocesses is
disjoint — that is, they all work with completely distinct sets of events.
(This is what we want, since whether or not object 1 is defined should have no
bearing on what we can read or write for object 2.) No matter how we try to
combine these subprocesses together, there is no way that they can interact with
each other — so we might as well use the simplest composition operator.

#### Testing it out

We now (presumably!) have a working CSP process describing the internal
consistency axiom from the paper.  If my running commentary isn't enough to
convince you that the process is correct, we can also use [FDR][] to perform
some tests!  This is very preliminary; once we have the full consistency model
specified, along with the reference implementation from the paper, we'll do some
more thorough (and interesting!) checks.

[FDR]: https://www.cs.ox.ac.uk/projects/fdr/

Our goal is to have FDR check whether the following transactions
are valid:


\\[
\textrm{UnrepeatableRead} =
\texttt{write}{.}1{.}1 \rightarrow
\texttt{read}{.}1{.}2 \rightarrow
\textrm{STOP}
\\]

\\[
\textrm{RepeatableRead} =
\texttt{write}{.}1{.}1 \rightarrow
\texttt{read}{.}1{.}1 \rightarrow
\textrm{STOP}
\\]

The first CSP process represents an "unrepeatable read": the transaction writes
the value 1 into object 1, but when it then tries to read from object 1, it
somehow gets the value 2.  This should be outlawed by internal consistency.  The
second CSP process should be valid; when we read the object after writing to it,
we get the value 1, as expected.

We can ask FDR to verify the following refinements:

\\[
\textrm{InternalConsistency}
\mathrel{\sqsubseteq_{T}}
\textrm{UnrepeatableRead}
\\]

\\[
\textrm{InternalConsistency}
\mathrel{\sqsubseteq_{T}}
\textrm{RepeatableRead}
\\]

<div class="aside-def">
## Refinement

Model checking in CSP is all about **refinement**.  \\(Spec
\mathrel{\sqsubseteq} Impl\\) says that \\(Impl\\) *refines* \\(Spec\\): all of
the things that \\(Impl\\) tries to do are "allowed" by \\(Spec\\).

CSP also provides several **semantic models**, which define what kinds of
behavior your \\(Spec\\) process can allow or deny.  Each semantic model has its
own (incredibly mathematical) definition of what "refinement" means.  The
simplest semantic model is **traces refinement** (\\(\sqsubseteq_{T}\\)), which
is normally used to verify a **safety** property — that is, that something bad
is not allowed to occur.  Later on we'll learn about more complex semantic
models that also let us check **liveness** properties — that is, that something
good *must* occur.
</div>

I've wrapped up FDR in a Docker image in my [csp-models][] repository; this
repository also contains a machine-readable version of the internal consistency
specification, and of the two refinements we want to check.  If we run this
example through FDR, we should get:

[csp-models]: https://github.com/dcreager/csp-models

~~~
$ ./refines concurrency-models/read-atomic.csp
InternalConsistency [T= UnrepeatableRead:
    Result: Failed
        Counterexample (Trace Counterexample)
            Specification Debug:
                Trace: <write.1.1>
                Available Events: {|write, read.2, read.1.1|}
            Implementation Debug:
                UnrepeatableRead (Trace Behaviour):
                    Trace: <write.1.1>
                    Error Event: read.1.2

InternalConsistency [T= RepeatableRead:
    Result: Passed
~~~

(Note that I haven't included the entire output, just the important bits.)

We get the results that we expect: \\(\textrm{UnrepeatableRead}\\) is invalid,
and \\(\textrm{RepeatableRead}\\) is valid.  Even better, for the invalid
transaction, FDR gives us a counterexample telling us what went wrong!  Digging
into the counterexample, we see that after the sequence of events \\(\langle
\texttt{write}{.}1{.}1 \rangle\\), our specification allows the following
events:

- we can \\(\\texttt{write}\\) whatever we want;

- we can \\(\\texttt{read}\\) whatever we want from object 2 (since it's still
  undefined);

- but if we want to \\(\\texttt{read}\\) from object 1, the only value we're
  able to see is 1.

However, at this point our implementation wants to \\(\textrm{read}\\) the value
2 from object 1.  Since that's not allowed by our specification, the refinement
fails, and we have proof that \\(\textrm{Unrepeatableread}\\) does not satisfy
\\(\textrm{InternalConsistency}\\).

#### Next steps

That wraps it up for internal consistency!  In the next couple of posts, we'll
take a look at *external* consistency, and then we'll look at the reference Read
Atomic implementation from the paper, and try to use FDR to show that the
implementation really does implement Read Atomic.
