---
kind: article
created_at: 2014-05-14
layout: /post.html
title: "Tagged releases using git flow"
tags: [git]
---

I've used `git flow` for most of my software projects — specifically for those
that have versioned releases.  The [original
post](http://nvie.com/posts/a-successful-git-branching-model/) describing `git
flow` is still the best overview of how it works.  In short, you have a `master`
branch where **every** commit is a merge commit.  Each of these merge commits
represents a new release, and is tagged with the version number of that release.
The merge brings in all of the subsidiary commits and feature branches that make
up that release.  Ongoing work happens on a separate `develop` branch.  This is
where you merge in completed new features and bug fixes on a day-to-day basis.
`develop` should always be a stable version of the software — you don't merge a
feature branch into `develop` until it passes all of your tests and is
"complete" with regards to the feature you're trying to implement.

My favorite part of this model is how each release is just some tagged commit on
the `master` branch.  You want to see the code for the latest released version?
That's easy — `git checkout master`.  You want version 1.2.5 specifically?  Use
`git checkout 1.2.5` instead.

Unfortunately, the [`git flow` tool](https://github.com/nvie/gitflow) has
implemented a [slightly different
behavior](https://github.com/nvie/gitflow/issues/206) for awhile now.  That
patch makes `git flow` tag the last commit on the release branch, instead of the
merge commit on the `master` branch.  The reasons for this might be perfectly
valid, but it's not what I want, and it's not what the original `git flow` post
described.  That means that I can't use `git flow release finish` as-is.

<hr class="jump">

To recap the usual process:

    $ git flow release start ${VERSION}
    [ bump version numbers and release dates if needed ]
    $ git flow release finish -s -m ${PROJECT}-${VERSION} ${VERSION}

The first command creates a new release branch called `release/${VERSION}`.  If
you need to make any changes that bump version numbers or release dates in the
actual files in your repository, you do that on this new release branch.  You
then finalize the release, which merges the release branch (along with any new
commits on `develop`) into the `master` branch.  It also creates a tag named
`${VERSION}`, and because of the `-s` option, signs it with my GPG key.

The problem is that this tag is placed on the wrong commit.  That's easy to fix,
though:

    $ git tag -sf -m ${PROJECT}-${VERSION} ${VERSION} master

I also like to create a `-dev` version tag on the latest `develop` branch, so
that `git describe` always gives me a nice looking version number, regardless of
which branch I'm on:

    $ git tag -s -m ${PROJECT}-${VERSION}-dev ${VERSION}-dev develop

And that's it!
