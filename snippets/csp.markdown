Communication is synchronous.  That means that an event does not mean "send" or
"receive" or anything else in particular.  It's just some visible evidence of
something occurring in a process.  The same event can have a different
interpretation in different processes.  It *can* mean "send" or "receive", and
if they do, then composing those two together represents a traditional message
passing communcation.

Synchronous also means that a process can be *prevented from making progress*.
If you compose two processes together, 


CSP is *composable*.  That means that it's very easy for you to construct a
complex CSP process in a modular way: you break your system down into parts,
model each part separately, and then combine them back together.  This is useful
from a design or engineering standpoint, because it makes your job as a modeler
easier.  But it's also useful when we start to analyze processes formally, or
with a tool, since this composability provides us with some interesting
optimizations.


There are two main dialects of CSP: a *whiteboard* language, called CSP, and a
*machine-readable* language, called CSP<sub>M</sub>.


<h5>Tool support</h5>

Another common critique of formal methods is that there can be a tradeoff
between how much detail you include in your models and 

State explosion vs precision of your model
