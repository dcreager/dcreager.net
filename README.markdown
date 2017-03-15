# dcreager.net

This is the source tree for [dcreager.net](http://dcreager.net/).  Most of the
pages on the site are generated from Markdown files or partial HTML files in
this git repository.  I use [nanoc](http://nanoc.ws) to compile the source into
a full static HTML website.


## Initial setup

If you are going to make changes to the website, you need to have all of the
necessary tools installed.  You can use [Bundler](http://bundler.io) to install
everything for you.

First, make sure you have Ruby and RubyGems installed.  The commands for this
will be specific to your distribution:

    # Ubuntu
    $ sudo apt-get install ruby rubygems

    # Arch
    $ sudo pacman -S ruby

Then install Bundler:

    $ gem install bundler

Then use Bundler to install all of the tools needed to work on the website:

    $ cd $LOCAL_CLONE_OF_WEBSITE
    $ bundle

If I ever add additional dependencies, you only need to re-run the last command
to bring everything up to date.


## Directory structure

I encourage you to read through Nanoc's [documentation](http://nanoc.ws/docs/)
for more details; it's quite good.

For the most part, you're going to be working with files in the *content*,
*layouts*, and *media* directories:

- *content* is where the Markdown or HTML source for each individual web page
  lives.  Each file in this directory starts with some [YAML
  front-matter](http://nanoc.ws/docs/basics/#attributes) describing the page,
  and then contains the Markdown or partial HTML for the main content of the
  page.  These files do **not** contain any of the web chrome that makes it a
  full, valid HTML file.

- *layout* is where the [web chrome](http://nanoc.ws/docs/basics/#layouts)
  lives.  There are quite a few layouts, so that we can define navigation bars
  and whatnot once for each of the various sections of the website.  The
  `layout` field in each page's YAML front-matter defines which layout will be
  used to render a particular page.  And some of the top-level layouts use
  [`render` tags](http://nanoc.ws/docs/basics/#as-partials) to include other
  sub-layouts.

- *media* is where you should put any files that Nanoc should copy through,
  unchanged, into the full site.  You'll find my custom CSS in *media/css*, and
  images used for figures in my articles in *media/figures*.


## Compiling the website

If you've made any changes to the source files, you can run the following:

    $ bundle exec nanoc compile

(You can abbreviate `compile` to `co` if you like.)

This will compile the source files and produce a version of the static website
in the *.output* directory.

Running `nanoc compile` repeatedly can be a pain, so you can also start a
watcher process, which will automatically recompile the site whenever it detects
any changed files in any of the various source directories:

    $ bundle exec guard


## Viewing the website locally

When you create and edit files, you're naturally going to want to look at the
beautiful results.  Don't open the files in *.output* directly; the links in the
files are set up with the assumption that they're going to be served from a web
server.  Luckily there's a quick-and-dirty web server that you can use to view
your edits locally:

    $ bundle exec nanoc view

That will let you the currently compiled version of the site (including any
local changes that you've made) at [localhost:3000](http://localhost:3000).


## Deploying the results

If you've made some changes (and made nice git commits out of them), then you
can deploy the compiled site so that it's live at
[dcreager.net](http://dcreager.net/).

    $ bundle exec nanoc compile
    $ bundle exec nanoc deploy

(The `compile` command will usually be superfluous, since you'll have been
compiling or watching as you were editing.  But it doesn't hurt to make sure.)

The `deploy` command will use rsync to copy the current contents of the
*.output* directory up to the documentation server.  It will then be immediately
visible to anyone visiting the site.

**NOTE**:  In order for the `deploy` command to work, you'll need to have an SSH
key that is allowed to log into *dcreager.net*.


## Non-HTML pages

Not all of the pages on the website are generated from partial HTML source.


### Tutorials

There are (or might eventually be) some tutorial pages that are generated
directly from a compilable or runnable source code file, using
[Rocco](http://rtomayko.github.io/rocco/).  All of the setup is already taken
care of; the only thing you need to know is that all of the `*.c` files (and
possibly other extensions; see the *Rules* file for the most up-to-date list) in
the *content* directory will be processed using Rocco to create a tutorial page.
Look at some of the existing tutorials for details.  Short version: use Markdown
inside of `//` comments for the explanatory text, and use the regular `/* */`
notation for comments that should stay in the code in the rendered tutorial
page.  The C source file should be completely valid — you should be able to
compile it and run it as an actual program.  (There's no automated tooling to
verify this, though.)
