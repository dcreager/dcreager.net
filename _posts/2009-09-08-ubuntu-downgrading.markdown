---
layout: post
title: Downgrading packages in Ubuntu
tags: ubuntu
---

## What kind of trouble did you get yourself into this time?

I have recently been setting up a new machine that I'll be using
jointly as a work machine and a MythTV frontend/backend.  As part of
setting this up, I've had several issues getting the integrated ATI
Radeon HD 3200 display board to work correctly.  I'll save the details
for another post, but the short version is that none of the three
available X drivers (ATI's `fglrx` and the open-source `radeon` and
`radeonhd`) seem to drive the HDMI output connector correctly.

As part of testing this, I used [Tormod Volden's
packages](https://launchpad.net/~tormodvolden/+archive/ppa) for the
`radeon` and `radeonhd` drivers, which are based on newer releases
than are available in the mainline Jaunty package trees.  For some
packages, they are even based on bleeding-edge git checkouts, rather
than released tarballs.  (More notes on using `radeonhd` can be found
[here](https://help.ubuntu.com/community/RadeonHD).)  While neither
package was able to use the HDMI connector properly, the `radeon`
driver was able to give me output on the VGA connector, with full 2D
and video acceleration (which I needed for the MythTV front-end).  My
monitor (a Westinghouse L2410NM) was auto-detected through the
connection, so my *xorg.conf* is trivial.

However, neither open-source driver has 3D acceleration support yet.
To get this, I'll need to use ATI's `fglrx` driver.  Not a problem,
right?  Just install the packages, then change the `radeon` entry in
your *xorg.conf* over to `fglrx`, and you're good to go!  Except that
by using Tormod's package tree to pick up the latest `radeon` and
`radeonhd` drivers, I *also* pick up more recent versions of `xserver`
and friends, and it would seem that the `fglrx` drivers don't play
well with the new version — I get segfaults when the X server tries to
start using `fglrx`, which didn't happen before installing Tormod's
drivers.

So, I need to back out all of the packages installed from Tormod's
tree, and revert back to the versions of these packages that I had
previously installed from the mainline Jaunty trees.

## How you might think it should work

Anticipating that I might want to remove Tormod's tree from my
*sources.list*, I used the nice *sources.list.d* facility.  Instead of
putting all of your package sources into a single file, you put them
in separate files in the */etc/apt/sources.list.d* directory.  That
way, activating and deactivating a particular package tree is as
simple as moving a file and re-running `apt-get update`:

    $ sudo mkdir /etc/apt/sources.list.d.unused
    $ sudo mv /etc/apt/sources.list.d/tormod.list /etc/apt/sources.list.d.unused
    $ sudo apt-get update

Now we've removed Tormod's packages from our list of available
packages.  Ideally, we could run an `apt-get` command to “downgrade”
our packages to those that are mentioned in our package lists.
However, I wasn't able to find such a command.  If you try to run
`apt-get upgrade` or `apt-get dist-upgrade`, APT will see that you
have more recent versions of the X packages installed (the ones from
Tormod's tree), and won't overwrite those with the older packages
mentioned in the sources that we've activated.  Normally, this is the
behavior that we want; but in this case, we're boned.

## Downgrading explicitly

Instead, we'll need to remove the newer packages using `apt-get
remove`, and then reinstall them using `apt-get install`.
Unfortunately, you have to specify which packages you want to remove
and reinstall on the command line.  So, we need a command pipeline
that will tell us “all packages installed from Tormod's tree”, so that
we can call `apt-get remove` and `apt-get install` on them.

First, we can get a list of the packages defined by a particular
source by reading the files in */var/lib/apt/lists*.  This directory
contains the local copies of the package list files that are
downloaded from each source.  Each source has its own files, which
makes it easy to distinguish the packages that came from Tormod's tree
from those that came from mainline Jaunty.  However, there will only
be files for the activated sources — so if you've moved the
*tormod.list* file like I did above, you won't find a file for
Tormod's package tree.  First, we'll have to reactivate his package
tree:

    $ sudo mv /etc/apt/sources.list.d.unused/tormod.list /etc/apt/sources.list.d
    $ sudo apt-get update

Now that we're sure that Tormod's tree has a file in
*/var/lib/apt/lists*.  The filename is based on the URL of the
*sources.list* entry that you want to deal with.  So if you're
following these instructions for a different package tree than
Tormod's X packages, you'll need to find the correct file and `grep`
it instead.  You'll see several files for each source; you want the
file that ends in `Packages` and contains `binary` and your
architecture in the name.

Once we find the file, we can extract out the `Package` lines from it:

    $ cd /var/lib/apt/lists
    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages
    Package: xserver-xorg-video-ati
    Package: xserver-xorg-video-ati-dbg
    Package: xserver-xorg-video-radeon
    Package: xserver-xorg-video-radeon-dbg
    [...]


We only want the package names, though, so we need to use `sed` to
strip out the `Package:` prefix:

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //'
    xserver-xorg-video-ati
    xserver-xorg-video-ati-dbg
    xserver-xorg-video-radeon
    xserver-xorg-video-radeon-dbg
    [...]

This tells us which packages are defined in Tormod's package tree.
However, we can't pass all of them into `apt-get install`, because we
probably haven't installed all of the packages that Tormod made
available.  We can use `dpkg-query` to see which of these packages are
actually installed:

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //' \
      | xargs dpkg-query -W 2>/dev/null
    gnome-screensaver	2.24.0-0ubuntu6
    libgl1-mesa-dev	7.4-0ubuntu3.2
    libgl1-mesa-swx11-dev	
    mesa-common-dev	7.4-0ubuntu3.2
    xdmx	
    xnest	
    [...]

The first four lines of this output describe packages that are
installed, while the last two describe packages that are not
installed.  We can distinguish between the two cases by the presence
or absence of the version number.  Looking closely, we can see that
the package name and version number are separated by a tab character.
Hopefully, that tab isn't printed for uninstalled packages, which
would let us just look for a tab character to filter out the
uninstalled packages.  Let's try it and see:

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //' \
      | xargs dpkg-query -W 2>/dev/null \
      | less -U
    gnome-screensaver^I2.24.0-0ubuntu6
    libgl1-mesa-dev^I7.4-0ubuntu3.2
    libgl1-mesa-swx11-dev^I
    mesa-common-dev^I7.4-0ubuntu3.2
    xdmx^I
    xnest^I
    [...]

Ah well, it was worth a try.  But as a consolation, we can check for a
tab followed by any other character, and we'll get the installed
packages:

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //' \
      | xargs dpkg-query -W 2>/dev/null \
      | grep '\t.'
    gnome-screensaver^I2.24.0-0ubuntu6
    libgl1-mesa-dev^I7.4-0ubuntu3.2
    mesa-common-dev^I7.4-0ubuntu3.2
    xscreensaver-data^I5.08-1~rc2
    [...]

Note that the `\t` in the above is an actual tab character, and not a
backslash followed by a T.  In most shells, you type in the tab
character by pressing Control-V and then the Tab key.

Now we have a list of installed packages that *might* have come from
Tormod's package tree.  I say might have, because it's possible that
the mainline Jaunty has a newer version of a package that Tormod's
tree does.  Earlier, we said we were looking for the packages that
*definitely* came from Tormod's tree, so that we can reinstall them.
If we reinstall all of the packages in this list we've just created,
then we might end up reinstalling a package that we don't need to.
But that's not the worst thing in the world — we're only reinstalling
a dozen or so packages in total, so the extra work of reinstalling a
couple of packages that we don't need to won't really kill us.

So now we need to pass this list into `apt-get` to reinstall the
packages.  `apt-get` only wants the package names, not the versions,
so we'll need to use `sed` again to strip these out: (Like before, the
two `\t` entries are actual tab characters.)

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //' \
      | xargs dpkg-query -W 2>/dev/null \
      | grep '\t.' \
      | sed -e 's/\t.*$//'
    gnome-screensaver
    libgl1-mesa-dev
    mesa-common-dev
    xscreensaver-data
    [...]

Perfect!  Let's save this package list into a file.

    $ grep Package ppa.launchpad.net_tormodvolden_ppa_ubuntu_dists_jaunty_main_binary-amd64_Packages \
      | sed -e 's/^Package: //' \
      | xargs dpkg-query -W 2>/dev/null \
      | grep '\t.' \
      | sed -e 's/\t.*$//' \
      > /tmp/tormod.packages

Then, we can deactivate Tormod's tree, and then forcably reinstall any
of the packages that we might've gotten from it:

    $ sudo mv /etc/apt/sources.list.d/tormod.list /etc/apt/sources.list.d.unused
    $ sudo apt-get update
    $ sudo apt-get remove `cat /tmp/tormod.packages`
    $ sudo apt-get install `cat /tmp/tormod.packages`
