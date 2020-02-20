---
kind: article
created_at: 2009-08-07
layout: /post.html
title: Adding Disqus comments
tags: [dcreager.net, disqus, jekyll]
---

I've just enabled comments on the posts on my website.  On its own,
that's not a particularly unique or exciting feature.  However, I'm
using [Disqus](http://disqus.com) as the comment engine, and the way
in which I've integrated Disqus into my Jekyll-powered website might
be of interest to others.  (Thanks to [Jack
Moffitt](http://metajack.im/) for the idea!)

## The “generic code” installation target

At the core of Disqus is a snippet of HTML and JavaScript that's
embedded into each of the webpages that you want to contain a comments
section.  If you're using one of the standard blog engines for your
website, Disqus can automatically “install” itself, adding the Disqus
snippet to your pages for you.  If you're not using one of these blog
engines, however, you have to install the Disqus snippet yourself.

Luckily, this is very easy.  If you choose the “generic code”
installation target when setting up the Disqus account for your
website, you see the snippet of code to include.  For
[dcreager.net](http://dcreager.net/), it looks like this:

<% highlight :html, :line_numbers => :table do %>
<div id="disqus_thread"></div>

<script
   type="text/javascript"
   src="http://disqus.com/forums/«account»/embed.js">
</script>
<noscript>
  <a href="http://«account».disqus.com/?url=ref">View the discussion thread.</a>
</noscript>

<a href="http://disqus.com" class="dsq-brlink">
  blog comments powered by <span class="logo-disqus">Disqus</span>
</a>
<% end %>

<div class="alert alert-warning">
  <p>
    <strong>UPDATE</strong>: Please do not copy-paste this JavaScript
    snippet directly into your own website!  There are several
    occurrences of the string “<code>«account»</code>” which should be
    replaced with the name of your own Disqus account.  If you don't
    do this, your comment threads won't be associated with your
    account, and you therefore won't be able to moderate or export
    those comments.
  </p>

  <p>
    To see the specific JavaScript snippet for your own site, please
    go to <a
    href="http://www.disqus.com/comments/install/">http://www.disqus.com/comments/install/</a>.
  </p>
</div>

I have to wrap this in another `<div>` element, to fit into the CSS
layout that I'm using, but otherwise it's a straightforward
copy/paste.

The `disqus_thread` element is a placeholder.  The *embed.js*
JavaScript retrieves the comments that are linked to the current page,
formats them according to the style that you've chosen, and injects
the resulting HTML into the `disqus_thread` element.  The end result
is something like the comment section that you see at the end of this
page.

## Jekyll layouts

At this point, we have the Disqus snippet that we need to include into
each comment-enabled page on the site.  The naïve solution would be to
manually add this snippet to each of the pages that we want to contain
comments.  But of course, as good software engineers, we want to avoid
that kind of repetition — and Jekyll layouts give us good way to do
that.

For my site, I've decided that I want to include a comment section on
each dated “post” — the articles that you see under the “Recent Posts”
on the front page.  I don't want to include comments, for instance, on
my [publication list](/publications/).

Luckily, I already have a system of Jekyll layouts that implements
this distinction.  For instance, dated posts have a “Last updated”
entry at the bottom, which the front page and publications list don't
contain.  The site currently has two layouts defined:

* [*default.html*](http://github.com/dcreager/dcreager.net/blob/master/_layouts/default.html) —
  defines the overall layout of the site: the background logo, the
  navigation bar, the box containing the text of each post, etc.  All
  of the pages on the site use this layout, even if they don't
  reference it directly in their YAML front-matter.

* [*post.html*](http://github.com/dcreager/dcreager.net/blob/master/_layouts/post.html) —
  adds the “last updated” entry at the bottom of a dated post.

So, if I want to include Disqus comments on all of my dated posts, I
can just add the HTML/JavaScript snippet to the *post.html* file.

## Using template variables

This solution is nice, in that we don't have to duplicate the Disqus
snippet on each of the post pages, but we can take this one step
further.  What if I decide that I want to include comments on one of
the pages that isn't a dated blog post?  Now I have to duplicate that
snippet again — maybe not in the page's source, but at least in the
layout that that page uses.

Instead, I'm going to use a new variable in the YAML front-matter to
give me fine-grained control over which pages have comments.  If I
want a page to have comments, I just make sure to include

    comments: true

in that page's YAML header.

Then, instead of including the Disqus snippet directly in the
*post.html* layout, I put it into the *default.html* layout that every
page eventually uses.  I wrap the snippet in a Liquid `if` statement
to only include the Disqus comment section if the `comments` YAML
variable is `true`:

<% highlight :html, :line_numbers => :table do %>
{% if page.comments %}

<div id="disqus_thread"></div>

<script
   type="text/javascript"
   src="http://disqus.com/forums/«account»/embed.js">
</script>
<noscript>
  <a href="http://«account».disqus.com/?url=ref">View the discussion thread.</a>
</noscript>

<a href="http://disqus.com" class="dsq-brlink">
  blog comments powered by <span class="logo-disqus">Disqus</span>
</a>

{% endif %}
<% end %>

If the `comments` YAML variable isn't defined for the current page,
the `if` statement treats it as false, giving us the behavior that we
want — no comments section unless we ask for it.

Finally, to ensure that all dated posts get a comments section,
without me having to explicitly ask for it, I add

    comments: true

to the YAML front-matter of the *post.html* layout.
