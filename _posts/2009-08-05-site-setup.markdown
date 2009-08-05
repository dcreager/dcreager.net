---
layout: post
title: Site layout
---

# Site layout

This post will probably end up being more useful to me than to anyone
else who might stumble across the page.  Here I'm going to document
how I've set up my homepage, from a technical standpoint.

## Directory layout

The content of the website is stored in a Git repository (found
[here](http://github.com/dcreager/homepage/)).  Most of the pages are
originally written in Markdown.  I use
[Jekyll](http://github.com/mojombo/jekyll/) to process the Markdown
pages into a static website.

The Git repository contains a standard Jekyll layout:

* Dated “posts” (such as blog entries) are placed in the *\_posts*
  directory.

* HTML layouts are placed in the *\_layouts* directory.

* All other content (CSS, images, other pages) lives wherever I
  please; that directory structure is reproduced on the live site.

One difference is that I include the *\_site* directory in the Git
repository; most people seem to include this directory in their
*.gitignore* file so that it's not tracked by Git.  Doing so allows me
to check out the repository and have a working copy of the site,
without having to have Jekyll and its dependencies installed on that
machine.

## Editing and deploying changes

While I edit my pages, I keep a

    jekyll --server --auto

instance running in the background, which allows me to view a local
copy of the new website as I save changes.

For deployment, I have a (non-bare) clone of the Git repository on the
Dreamhost machine that hosts my website.  Once I have a change that
I'm ready to deploy, I make a new Git commit and push it to the
Dreamhost clone.  Since I include the *\_site* directory in my
commits, this places the latest copy of the website onto the Dreamhost
filesystem, ready to go.

Pushing doesn't automatically update the checked-out HEAD on the
remote system, however, so there's an additional step needed.  Once
I've pushed the changes to Dreamhost, I run the following from the
Dreamhost clone:

    git reset --hard master

which updates the working copy on disk to be the same as the latest
commit that I just pushed.  At this point, the Dreamhost clone
contains the latest copy of the site in its *\_site* directory.

Dreamhost is expecting to serve my website out of a particular
directory within my home directory; the final step is having this
served directory be a symlink to the *\_site* directory of the
Dreamhost clone.  Et voila!
