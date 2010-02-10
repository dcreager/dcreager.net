---
layout: post
title: Extracting setuptools version numbers from your git repository
tags: [setuptools, python, git]
---

Just like everyone else, we're using
[setuptools](http://pypi.python.org/pypi/setuptools) as the core of
the build system for our Python-based projects.  For the most part,
this has been a painless, straightforward process.  However, one
lingering annoyance is that we've been specifying the version number
directly in our _setup.py_ files:

{% highlight python %}
from setuptools import setup

setup(
    name = "awesomelib",
    version = "1.2",
    # ...etc
)
{% endhighlight %}

On our maintenance branches, we get a nice _awesomelib-1.2.tar.gz_
file when we run `python setup.py sdist`.  On our development branch,
we've also got the following _setup.cfg_ file:

{% highlight ini %}
[egg_info]
tag_build = dev
tag_date = true
{% endhighlight %}

That gives us tarballs like _awesomelib-1.2dev-20100210.tar.gz_ on our
development branch.  Because we're using the `dev` suffix, which
setuptools considers to be a "prerelease", we have to remember to
increment the version number in development whenever we cut a new
release.  The end result is that we have a longish process for
creating releases.  If we want to create a new 1.3 release, we have to
do the following:

 1. Create a new maintenance branch for 1.3:

        $ git checkout -b maint-1.3 master

 2. Update the _setup.cfg_ file to remove the `tag_build` and
    `tag_date` entries.  Commit this with a "Tagging version 1.3"
    commit message.

 3. Back on the development branch, update _setup.py_ to increment the
    "development version" to 1.4.

Granted, this isn't horribly difficult, but we can do better.


## Calculating the version automatically

Taking a page from the
[_GIT-VERSION-GEN_](http://git.kernel.org/?p=git/git.git;a=blob;f=GIT-VERSION-GEN)
script in git's source code, we're going to use the `git describe`
command to automatically generate the version number.

Our logic is implemented in a new `get_git_version()` Python function,
which you can call directly from your _setup.py_ scripts..  You can
find the source code in a [Github
gist](http://gist.github.com/300803).  Our basic strategy is:

 1. First, try to use `git describe` to create a version number.

 2. If this fails, then we're most likely not in a git working copy.
    Probably, someone downloaded a release tarball and unpacked it,
    and we're running inside of there.  In this case, `git describe`
    can't give us a version number.  Instead, we're going to make sure
    we include a _RELEASE-VERSION_ file in every tarball that we
    create.  So, if `git describe` fails, we fall back on the contents
    of this file as our version number.

### Tag names as version numbers

One thing to notice about this strategy is that we use the output of
`git describe` directly as our version number.  This means that our
tag names should be simple version numbers, without decoration.  To
create the awesomelib 1.3 release from our example, we'd just do:

    $ git tag -s 1.3

(Note that the tag needs to be an annotated or signed tag in order to
be picked up by `git describe`.)

On our development branch, once we've created new commits on top of
the release point, we'll start getting output like this from `git
describe`:

    1.3-4-g6f32

This is a valid setuptools "postrelease" — setuptools will consider
this to be a more recent version than `1.3`, which is exactly what we
want.  This eliminates the need to maintain different _setup.cfg_
files in our development and maintenance branches.

### Getting the version number of a distribution tarball

Another thing to notice is that we need to maintain a
_RELEASE-VERSION_ file, ensuring that it always contains the current
version, and always including it when we create any source packages.
That way, we can still get the current version number, even if we
can't get it from `git describe`.

To keep the _RELEASE-VERSION_ file up-to-date, the `get_git_version()`
function always read in the current contents of the file as its first
step.  If the output of `git describe` differs from what's in the
file, we update the file with the new output before returning the
version.

This ensures that the file has the right contents, but we also have to
make sure we include it in our source packages.  To do this, we simply
add the following line to our _MANIFEST.in_ file (creating it if
necessary):

    include RELEASE-VERSION

(Note that we don't want the _RELEASE-VERSION_ file to be checked into
the git repository, so we also add it to the top-level _.gitignore_
file.)

## The simpler release process

With this script, our release process is now much simpler:

 1. Create a maintenance branch if you want to.

 2. Create a signed or annotated tag, whose name is the new version
    number.

Most importantly, no extra commits are needed, since we don't have to
edit any version numbers or maintain different _setup.cfg_ files.
