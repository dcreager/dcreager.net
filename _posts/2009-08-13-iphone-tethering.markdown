---
layout: post
title: iPhone tethering
---

While attending the [Mil-OSS conference](http://mil-oss.org/) this
week, I had the opportunity to use one of the coolest features of my
new iPhone 3GS — Bluetooth Internet tethering.  Assuming that your
mobile carrier allows the feature on their network, it provides a very
easy way to have a persistent Internet connection, for those
situations where a free Ethernet drop or WiFi access point aren't
readily available.

I happened to have my Dell Mini 9 (running Ubuntu Jaunty) with me for
the conference, rather than my MacBook, so I thought that it might be
difficult to get the Bluetooth connection working between the phone
and laptop.  This [blog
posting](http://xn--9bi.net/2009/06/17/tethering-iphone-3-0-to-ubuntu-9-04/),
however, provided exactly what I needed.

Following this process, the connection worked without a hitch.  One
caveat is that I didn't need to explicitly pair my phone with my
laptop; the first time I ran the `pand --connect` command, my phone
prompted me with the pairing confirmation dialog.  Later connections
didn't require re-pairing.

As for bandwidth, the connection was perfectly reasonable for email
and basic Web surfing as long as I had decent 3G coverage — 3 of 5
bars or higher and I was good to go.  I also did an `apt-get` package
update as a “beefier” test; I was usually seeing around 20-30 Kb/sec
of download speed, which would be fine for small daily updates, but
would probably be unworkable for something large like a GNOME,
texlive, or GHC update.  All in all, not bad for some surreptitious
email checking during the talks.

## Command-line scripts

One thing that can be cumbersome about the instructions on the blog
post is that you have to run the `pand` and `ifup`/`ifdown` commands
separately each time you want to start or stop the Bluetooth
connection.  Not a huge waste of effort, to be sure, but we can do
better.  So I wrote a quick little Bash script that will start and
stop the connection with a single command.  You can find the script in
[this
commit](http://github.com/dcreager/home/commit/f5c0db363049f7433494924c63d4a2a19e325b6c)
on my Github page.

The script is fairly straightforward; you start your Bluetooth tether
by running

    tether up

and you stop it by running

    tether down

Instead of hard-coding your phone's Bluetooth MAC address into the
script itself, you place it into the *$HOME/etc/bluetooth.conf* file.
This file isn't checked into Git, so that I'm not putting my personal
MAC addresses into the public repository.  Instead, the commit
contains a *$HOME/etc/bluetooth.conf.sample* file, which you copy over
to *bluetooth.conf*, and then edit appropriately.

## Issues

While this script worked great for during the conference, there are
two main issues with it as it stands.

* **dbus integration** — Many applications now listen to the system's
  dbus message bus to determine different facts about the current
  state of the system, including whether there's an active Internet
  connection.  The NetworkManager application knows to publish the
  correct dbus messages when it starts and stops a network connection.
  The `tether` script does not.  This means that each time I open up
  Firefox after turning on the tether, I have to manually uncheck the
  *File › Work Offline* menu option to be able to access any web
  pages.

  Fixing this issue should only require finding the appropriate dbus
  messages to send, and adding them to the `up` and `down` cases in
  the script.

* **GUI access** — While I don't mind running a simple command to
  activate and deactivate the network connection, I realize that a GUI
  control within the NetworkManager applet would be more ideal.  (This
  would also eliminate the first issue, since NetworkManager would
  then send the appropriate dbus messages when the connection is
  started or stopped.)  Luckily, Dan Williams has recently [added
  Bluetooth PAN
  support](http://blogs.gnome.org/dcbw/2009/07/10/unwire-with-networkmanager/)
  to `nm-applet`.  The new code is in the bleeding edge tree, so
  hopefully it will make it into a release in time to be picked up for
  October's Karmic release.
