---
kind: article
draft: true
created_at: 2018-04-19
updated_at: 2018-05-02
layout: /post.html
title: "Clean git histories and code review workflows"
description: >-
  in which we get the browser to do our work for us
tags: [git]
---

<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

A couple of recent tweets really resonated me.

[Aditya][] started an interesting conversation:

[Aditya]: https://adityamukerjee.net/

<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">If you&#39;d asked me a few weeks ago whether I preferred squash-merges or merge commits, I&#39;d have been somewhat ambivalent.<br><br>After today, I&#39;m pretty sold on squash-merge as the better strategy. It&#39;s not without drawbacks, but on net, it&#39;s still so much better to work with.</p>&mdash; Aditya Mukerjee (@chimeracoder) <a href="https://twitter.com/chimeracoder/status/986376765567356928?ref_src=twsrc%5Etfw">April 17, 2018</a></blockquote>

I've always really liked having a "clean" history in the git logs of my
projects.  They look nicer!  There's a real joy to running [`git log`][] (or
even better, [`git log --graph --all`][`git log`]!) on a project that's
maintained a beautiful log.  There are a number of tangible benefits, too, which
[others][chris-beams] [have][tower] [described][gamberini] well enough that I
don't need to repeat them here.

[chris-beams]: https://chris.beams.io/posts/git-commit/
[tower]: https://www.git-tower.com/learn/git/ebook/en/command-line/appendix/best-practices
[gamberini]: https://www.slideshare.net/TarinGamberini/commit-messages-goodpractices

More than that, though, for a long time I assumed that if you were using GitHub
pull requests, there was One True Way to produce a clean history, echoed by
[Jesús]:

[Jesús]: https://sevein.com/

<blockquote class="twitter-tweet" data-conversation="none" data-lang="en"><p lang="en" dir="ltr">Best merges to me are the fast-forward merges that bring a number of meaninful, atomic, functional and self-contained commits! Not always possible...</p>&mdash; Jesús García Crespo (@sevein) <a href="https://twitter.com/sevein/status/986412919964368896?ref_src=twsrc%5Etfw">April 18, 2018</a></blockquote>

But recently, just like Aditya, I've had a change of heart!  I realized that I
was thinking about pull requests, and what they represent in terms of the final
commit log I'm aiming for, in the wrong way.  In this post I'll describe the
come-to-Aditya moment that made me appreciate squash-commits.  And along the way
I'll have described two workflows that make it easy to keep a clean history: one
that doesn't use GitHub at all (and mimics how we did git-driven development
before GitHub existed), and one that does.

<hr class="jump">

#### What is a clean history anyway?

Before we can figure out how to achieve a clean history, we have to decide what
exactly we mean by one.  If you have a "clean history", then if I look at the
log of git commits for the project, all of the following should be true:

1. **Each commit must be self-contained**.

   This is sometimes stated as "your commit must make exactly one change", or
   "your commit should be as small as possible, but no smaller".

   What does this mean?  You should be able to describe — very briefly and
   precisely — what **_single_** change the commit makes to the code base.  What
   can the code do now that it couldn't do before?  What bug was there before
   that isn't now?

   If your commit makes two or more changes to the code, each of which you could
   describe very briefly, then you should split that into two or more commits!
   If your commit doesn't contain **_all_** of the modifications needed for that
   change, then it's not "self-contained", and needs to be bigger!

2. **If you have an automated test suite, the code must pass all tests after
   each commit.**

   This is a corollary to rule #1: if the test suite doesn't pass, then you've
   either added a new test without adding the new code that it's meant to cover,
   or you've updated some code without updating all of the effected test cases.
   In both cases, your commit isn't really self-contained.

3. **Each commit must have a clear description, formatted as described in the
   [Gospel According to Tim][tpope-commit-messages]**.

   Tim's post (now >10 years old!) contains the clearest description of the
   Generally Accepted Rules for writing commit messages.  These rules are pretty
   strict!  You can probably loosen some of grammar recommendations if you like.
   (...not really, just train yourself to follow all of them!)

   But you really should keep to a short 50-character subject line + longer
   description.  And that description should actually be descriptive!  Describe
   what the commit is trying to accomplish, and why.  You might end up copying a
   lot of text from your issue tracker, or from comments or documentation that
   is also in the code itself.  And that's okay!  I want to be able to see at
   least a high-level overview of each change without needing anything other
   than [`git log`][].

4. **Code review does not appear in the final history of your project**.

   This rule is probably the least obvious.  My hunch is that if you canvas a
   whole bunch of Clean History Zealots like myself, this is the rule most
   likely to be left out.  But it's important!

   Note that I'm **not** saying you shouldn't have some kind of [code review
   process][].  Quite the opposite, you must!  Even if you're working on a small
   project by yourself, you should have some kind of careful self-review step
   before you merge anything into master.

   As part of that code review process, you're going to make edits to your
   commits: fixing typos, changing names, maybe some refactoring.  Maybe in your
   first draft of the commit, you forgot to update some test cases.  Your
   reviewer (or your CI) points this out, and you duly update those test cases.
   Everything's great, and now you're ready to merge your change!  But: you
   don't want each of those drafts showing up as distinct entries in the final
   git log!  It's still a single logical change, deserving of exactly one
   self-contained commit in the history — it just took you a couple of tries to
   get that commit just right.

   Now, you *might* want a record of all of those drafts, and the back-and-forth
   comments between you and your reviewer, to be recorded *somewhere* for
   posterity.  But the git log is not the right place for it!  Instead, make
   sure that your code review tool is itself long-lived and searchable (archived
   GitHub pull requests, mailing list archives, etc).

[tpope-commit-messages]: https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
[code review process]: https://mtlynch.io/human-code-reviews-1/

Some of these rules might seem draconian, especially if you're used to a "just
press the green button already" workflow.  If you're not worried about keeping a
clean history, that's totally fine!  But if you are, I think this is the minimum
set of rules that you should strive to follow.

#### The "Linux kernel" workflow: mailing lists and personal insults

Since git was [originally developed][git history] to track changes to the Linux
kernel, and the kernel has a well-earned reputation for keeping a very clean
[history][weird kernel commits], it makes sense to look at the kernel's
[development workflow][kernel process].

[git history]: https://git-scm.com/book/en/v2/Getting-Started-A-Short-History-of-Git
[weird kernel commits]: https://www.destroyallsoftware.com/blog/2017/the-biggest-and-weirdest-commits-in-linux-kernel-git-history
[kernel process]: https://www.kernel.org/doc/html/v4.10/process/development-process.html

Linux is a big enough project, with a large enough population of contributors,
that there isn't really a single [Grand Unified Process][kernel trees].  That
said, kernel patches are typically reviewed via [regular emails sent to mailing
lists][kernel email].  This, by the way, is the reason for the [Generally
Accepted Rules][tpope-commit-messages] for commit messages — it's so the end
result looks "right" when formatted as an email message!  There are even several
git subcommands (which you probably never use) specifically designed for this
kind of workflow: [`git format-patch`][], [`git send-email`][], [`git apply`][],
[`git am`][].

[kernel trees]: https://www.kernel.org/doc/html/v4.10/process/2.Process.html#how-patches-get-into-the-kernel
[kernel email]: https://www.kernel.org/doc/html/v4.10/process/5.Posting.html

How does it work in practice?  Let's say you've developed an interesting new
feature, and you've structured using three commits that build on each other:

![basic feature branch](/git/workflows/figure-1.png){:.figure}

To get your new feature reviews, you send each proposed commit as a single
email.  You construct the emails so that later commits in the series look like
replies to earlier commits.  That way, the entire series shows up as a single
thread in your email client.  Other developers review your code by replying to
those same emails, using perfectly standard email clients.

At this point, the author will make changes to their code based on the comments
from the reviewers.  And here is where Rule 4 comes into play — the author
**_doesn't_** add new commits to their local branch titled "Edit typos" or
"Address comments from Sue":

![no typo commits](/git/workflows/figure-2.png){:.figure}

Instead, they'll rewrite the existing commits on their local branch (using [`git
commit --amend`][`git commit`], for instance, or the various flavors of [`git
rebase`][]), producing a **_completely new snapshot_** of commits:

![fix typos by rebasing](/git/workflows/figure-3.png){:.figure}

They'll then email that new snapshot to the same mailing list, just like they
did with the original one.

This process repeats until consensus is reached that the commits are ready to be
merged.  At that point, one of the subsystem maintainers responsible for that
section of code will merge the code, using [`git am`][] or [`git fetch`][]+[`git
merge`][]:

![merged into master](/git/workflows/figure-4.png){:.figure}

This process results in a clean history that satisfies all four of the rules
listed above.  The reviewers themselves ensure that Rules 1-3 hold, and the fact
that reviewer feedback is incorporated by **locally rewriting history** ensures
that Rule 4 is followed.  The resulting git log does **not** contain any record
of the back-and-forth that occurred during review; if you want to see that, you
have to dig through the mailing list archives.

#### The squash-merge model

For better or worse, "mailing lists and the command line" aren't a good fit for
most projects these days, and GitHub as won the resulting Great Code Hosting
Wars.  Is there an easy way that we can maintain a clean history, in the style
of the Linux kernel, following the Four Rules, while using only the core tooling
provided by GitHub?  Now that GitHub supports [squash merges][], the answer is a
resounding "maybe".

[squash merges]: https://blog.github.com/2016-04-01-squash-your-commits/

Before the advent of squash merges, the only option for incorporating the
contents of a PR into master was a [merge commit][github merge].  This actually
didn't seem so bad, since under the covers, it's just a [`git merge`][] of the
PR's feature branch — exactly what the maintainer would do in the mailing list
model!

[github merge]: https://help.github.com/articles/about-pull-request-merges/

But!  There's a disconnect, which is that the mailing list model also relies on
locally rewriting history to address any changes suggested by your reviewers.
You can do this with GitHub PRs — rewrite your commit history on your local copy
of the branch, and then force-push to GitHub.  For a good 10 years of my life,
this is exactly what I would do!

However, GitHub [strongly discourages][no force-pushes] you from force-pushing,
even to a PR feature branch.  Having worked recently with other code review
tools (like [Gerrit][] and [Phabricator][]), I eventually saw the error of my
ways, but only after yet another perfectly good review comment disappeared into
the ether, because it was linked to a commit that had been rewritten.

[no force-pushes]: https://help.github.com/articles/about-pull-requests/
[Gerrit]: https://www.gerritcodereview.com/
[Phabricator]: https://www.phacility.com/phabricator/

The key insight is that you're using a GitHub pull request to track two distinct
things:

  - the snapshot of what you want to appear in the final commit history, and

  - a record of the interactions between you and your reviewers.

In the mailing list model, these two things are tracked separately — one using
git commits, and the other using emails and replies.  In a GitHub PR, however,
**both are tracked using git commits**.  And that makes it very easy to conflate
the two.

So mentally, you have to teach yourself to separate the two.  You must create
**exactly one PR for each single commit** you want to appear in the final commit
history:

![single PR for each final commit](/git/workflows/figure-5.png){:.figure}

And that PR's feature branch should track your interactions with your reviewers:

![PR tracks reviewers](/git/workflows/figure-6.png){:.figure}

Even though as a fellow Clean History Zealot, it feels like fingernails on a
chalkboard to create a commit labeled "Fixed some typos", **that's what you have
to do**, because that's the most accurate description of how the new PR snapshot
differs from the previous one!  And luckily, with a squash-commit, **none of
those commits** will end up in the permanent history:

![squash merge](/git/workflows/figure-7.png){:.figure}

#### Squash-merging a *series* of commits

As I've described them so far, there's one remaining discrepancy between the
email-based workflow and the GitHub squash-merge workflow: **_series_** of
commits.  With the email workflow, it was possible to implement a larger feature
as a series of smaller steps (individual emails for each small commit, threaded
together using email replies), and importantly, **to have the entire series
reviewed as a unit**.  Reviewers would reply to the individual emails to make
comments on the small commits, and the patch series as a whole would only be
merged when **all of the commits in the series were approved**.

For a long time, I was confused about how best to accomplish this on GitHub,
because of the same misunderstanding about what GitHub PRs are meant to
represent.  A PR can contain several commits, and so my natural inclination was
to use a single PR to represent the entire patch series.  You would see each of
the individual commits in the PR, just like you would see in an email thread.
Moreover, the force-pushes that I was using to update the contents of the PR
lined up with the local rewrites that I would've performed in the email
workflow.

But we've already seen the problems with force-pushing to a PR branch, and how
we should use the commits on a PR to describe **_code review_**, not the
**_history of features_**.  So is there a way to author and review a patch
series using the squash-merge workflow?

Unfortunately, the answer is "no" — at least, not without paying careful
attention to what you're doing.  You can use a GitHub issue to describe the
patch series, and have separate PRs for each of the commits that you eventually
want to appear in the history.  The problem is that you can't (easily) have all
of those PRs out for review at the same time.  The problem is twofold:

  - When looking at a PR "later" in the series, you don't want the diff view to
    be cluttered with any of the changes added by the PRs "earlier" in the
    series.  Each PR should be reviewed as if all of the earlier PRs had already
    been merged.

  - You don't want to merge the "later" PR until all of the "earlier" PRs have
    been merged in.  For bonus points, you could somehow contrive for the entire
    series to be approved or not, rather than the individual pieces, but as long
    as you have to ability to block the merging of later PRs that's not strictly
    necessary.

You can use the [base branch][] of the PR to take care of the first part,
marking the "earlier" PR as the base of the "later" PR.  If you do this, GitHub
will correctly ignore any commit on the base branch when determining which
commits belong to the later PR, and when rendering the diff view.  But the base
branch also determines which branch the PR would get merged into once it's
approved, which isn't quite what we want.  [Grayson Koonce][] has described a
way to make the merging work, though it's not *quite* the result that I'm
looking for: you end up with a single commit in the git log, containing all of
the changes from all of the commits in the series.  I created the separate PRs
in the first place because I wanted each one to show up separately in the final
history.

[base branch]: https://help.github.com/articles/creating-a-pull-request/#changing-the-branch-range-and-destination-repository
[Grayson Koonce]: https://graysonkoonce.com/stacked-pull-requests-keeping-github-diffs-small/#4landingthestack

Since Grayson wrote his article, GitHub has added support for [changing the base
branch of a PR][].  This gets us **_so close_**: you'd mark the earlier PR as
the base of the later PR.  After everything is approved, you'd squash-merge the
earlier PR into master.  And at that point, you'd update the later PR so that
it's base is now master.  There are two problems: one is that because we
squash-merge, the commit that got merged into master for the earlier PR isn't
*quite* what's in the history of the later PR, and you'll have to perform a
couple of merges (one of them using `--strategy ours`) to get everything to line
up.  And more importantly, you can't currently change the base **_repo_** of the
later PR, which means that your feature branches and the master branch that you
want to merge into have to be in the same upstream repository.  If you're
following the usual strategy of editing feature branches in a personal fork, and
creating pull requests in the "upstream" repository, branch changing won't work.

[changing the base branch of a PR]: https://help.github.com/articles/changing-the-base-branch-of-a-pull-request/

So unfortunately, the current best practice is to not send out your later PRs
for review until your earlier PRs have been approved and merged; or
alternatively, for your reviewers to just deal with the fact that later PRs will
include changes from earlier PRs in the series.

To solve this for real, GitHub should separate the current [base branch][]
concept into two pieces: The **base** would be the branch this PR should
eventually be merged into (i.e, always "master" unless you're doing something
weird).  The **predecessor**, on the other hand, should be the "earlier" PR that
this one depends on.  The "Files changed" tab would only show the diff between
the predecessor PR's branch and the current PR's branch.  GitHub wouldn't allow
you to merge the PR until the predecessor was merged.  But when everything is
ready, it would merge the PR directly into the base branch.

Or, of course, we could just all go back to email patches around on mailing
lists.

[`git am`]: https://git-scm.com/docs/git-am
[`git apply`]: https://git-scm.com/docs/git-apply
[`git commit`]: https://git-scm.com/docs/git-commit
[`git fetch`]: https://git-scm.com/docs/git-fetch
[`git format-patch`]: https://git-scm.com/docs/git-format-patch
[`git log`]: https://git-scm.com/docs/git-log
[`git merge`]: https://git-scm.com/docs/git-merge
[`git rebase`]: https://git-scm.com/docs/git-rebase
[`git send-email`]: https://git-scm.com/docs/git-send-email
