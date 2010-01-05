---
layout: post
title: Exporting OmniGraffle documents from the command line
tags: [omnigraffle, papers]
---

[OmniGraffle](http://www.omnigroup.com/applications/OmniGraffle/) is
my tool of choice for creating figures for my papers.  It's biggest
drawback is that it's only available for Mac OS, which can make it
cumbersome if I'm working on one of my Linux machines and need to
create or modify a figure.  But it's ease-of-use and the quality of
the figures it creates are hard to beat.

Of course, creating the figure isn't enough — since I write my papers
in LaTeX, I have to export my figures into EPS or PDF (depending on
whether I'm creating a PostScript or PDF version of the paper) before
I can use them in my documents.  It's easy enough to use the Export
dialog to do this (keyboard shortcut: Command-Option-E), but ideally
I'd like the ability to export figures from the command line.  Coupled
with a good Makefile, this would let me run a simple `make paper`
command, and automatically re-export any necessary figures before
rebuilding the paper itself.

Luckily, OmniGraffle has always had rather good support for being
controlled via AppleScript.  The commands can be somewhat
undocumented, requiring a bit of trial and error, but while entrenched
in our PhD studies at Oxford, my colleague David Faitelson and I were
able to whip together a script that suited our needs.  I've recently
extracted the code from our Oxford SVN repository and uploaded it to
[Github](http://github.com/dcreager/graffle-export).

To install the script, just place the _graffle.sh_ and _graffle.scpt_
files into some directory that's on your `$PATH`, such as
_/usr/local/bin_ or _$HOME/bin_.  Then just run

    $ graffle.sh «format» «graffle file» «output file»

This will open the figure in OmniGraffle, and export it into the
format you specify on the command line, saving the result into
`«output file»`.

I'm still using version 4 of OmniGraffle, so I haven't had a chance to
verify that the script still works with version 5.  If you try it, and
it doesn't, feel free to open up a ticket on the [Github
site](http://github.com/dcreager/graffle-export/issues).
