---
layout: post
title: Updating graffle-export to work with OmniGraffle 5
tags: [omnigraffle, papers]
---

I recently upgraded to OmniGraffle 5, which caused my
[graffle-export](http://github.com/dcreager/graffle-export/) script to
break:

    $ graffle.sh ~/git/cwa/figures/analyst.graffle foo.pdf 
    OmniGraffle Professional 5
    /Users/dcreager/git/cwa/figures/analyst.graffle
    ./graffle.scpt: execution error: OmniGraffle Professional 5 got an error: The document cannot be exported to the "pdf" format. (-50)

(This was first reported to me by Nima Talebi as [a
ticket](http://github.com/dcreager/graffle-export/issues/issue/1) on
graffle-export's Github page.)

Before we can understand what error we're seeing, a little explanation
is in order.  The core logic of the OmniGraffle exporter is an
AppleScript.  Unfortunately, AppleScripts are stored in a binary
format, so if you go to the Github page, you can't easily view the
contents of the file.  The important line of the script is:

{% highlight applescript %}
save doc as format in file output
{% endhighlight %}

This saves an OmniGraffle document (which an earlier part of the
script makes sure is open) into a new output file.  The `output`
variable is the name of the desired output file, and is taken directly
from the _graffle.sh_ command line.

The `format` variable is what's causing us problems.  This parameter
to the `save` command tells OmniGraffle what format to use for the
file it's about to save.  This is how we get our export functionality;
we just give it the name of one of the export formats that it
supports.  The value of our `format` variable comes from the optional
first parameter to the _graffle.sh_ script.  Previously, if no value
was specified, I used "`pdf`" as a default.

Now, there's no real documentation that I've been able to find out
what values are allowed for this parameter.  I came across `pdf`
simply by trial and error.  "`PDF`" also seems to work, as does "`PDF
vector image`", which is the text that appears in the Format entry of
OmniGraffle's Export dialog box.

Or, to be more accurate, I should say that these values all work **in
OmniGraffle 4**.  Once you upgrade to version 5, these values no
longer seem to be valid choices for that parameter of the `save`
command — hence the error message.  A quick, non-exhaustive test shows
that none of these variations work for EPS or SVG, either.  The only
one that seems to still work is PNG.

So, what are we to do?  After looking at several other related
AppleScripts on the web, it seems that the `as` parameter of the
`save` command is optional.  After some experimentation, it turns out
that if you leave out this parameter, OmniGraffle tries to deduce the
correct output format based on the extension of your output filename.
So, we change our `save` command to the following:

{% highlight applescript %}
if format = "" then
  save doc in file output
else
  save doc as format in file output
end if
{% endhighlight %}

(We also have to modify the _graffle.sh_ wrapper script to not use
`pdf` as a default, but you can see that change [on
Github](http://github.com/dcreager/graffle-export/commit/b605b461a29b73ab4c21bd42b48549bd8bad8fcc).)
This lets us export a PDF version of a _.graffle_ file by giving an
output filename ending in _.pdf_, and leaving out the format
parameter.

I still have my old copy of OmniGraffle 4, and it looks like this
trick works with that version as well.  So, this is now the default
behavior, regardless of which version you have installed.

It would be nice if there was an accurate list of what values were
allowed for the `as` parameter, but we do have a working solution, at
least.  The only problem is if you want to export a PDF with a
different extension; with this solution, you'd have to export to a
_.pdf_ file and then rename it to your new extension.  But then again,
why would you want to do that?
