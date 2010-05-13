---
layout: post
title: Installing Ubuntu Lucid on a PowerPC QEMU virtual machine
tags: [ubuntu]
---

Part of the software I help develop at
[RedJack](http://www.redjack.com/) needs to be tested on both
little-endian and big-endian machines.  Little-endian machines are
easy, since everyone and their mother is running on a little-endian
Intel or AMD x86 chip.  It used to be that big-endian was pretty easy
to test, too — just break out your trusty Apple Powerbook G4 and
you're good to go.  Since Apple has shifted over to Intel chips,
though, the situation has changed.

Luckily, [QEMU](http://wiki.qemu.org/) has PowerPC as one of the
targets that it can emulate, so in theory, I can still easily test my
code on a big-endian machine by creating a QEMU PowerPC virtual
machine.  There's already a writeup about trying to install Debian
onto a QEMU VM
[here](http://machine-cycle.blogspot.com/2009/05/running-debian-on-qemu-powerpc.html).
[Aurélien Jarno](http://www.aurel32.net/) has graciously put together
downloadable disk images with Debian preinstalled.  If that's good
enough for your purposes, just go download those!  You won't need any
of the rest of the information on this page.

Unfortunately, I didn't want to run stock Debian; my little-endian
build machine is running Ubuntu Lucid, and for consistency, I wanted
my big-endian VM to be running the same.  As it turns out, this also
required a fair dose of masochism on my part.  There are several
issues that you'll encounter if you try to do this by hand.  Here is
my cheat sheet for getting around these issues.

Note that this isn't a full step-by-step account of how to install
Lucid onto a QEMU VM.  For now, I'm just trying to get my notes down
into a more permanent form.


## Getting QEMU

Note that I'm using Ubuntu Lucid as both the host and the guest OS for
this virtual machine; if you're running QEMU on a non-Ubuntu host,
then you can skip this section.

It seems that there's a bug with the current QEMU packages in Lucid.
If you try to run `qemu-system-ppc`, you'll get an error message about
missing the PowerPC BIOS image.  Joy.

Easiest way to get around this is to install QEMU from source.
Download the latest version from
[here](http://download.savannah.gnu.org/releases/qemu/).  Once you've
unpacked it, use the following to build:

    $ sudo apt-get build-dep qemu
    $ ./configure --prefix=/usr/local \
        --enable-sdl --enable-curses --enable-curl \
        --enable-kvm --enable-nptl --enable-uuid \
        --enable-linux-aio --enable-io-thread \
        --audio-drv-list=alsa
    $ make
    $ sudo make install

The first command is just an easy way to ensure that all of the
prerequisite libraries are installed.


## Booting the installation CD

Once you've got a working QEMU installed, you can find the PowerPC
Lucid installation CD
[here](http://cdimage.ubuntu.com/ports/releases/10.04/release/).  I've
decided to use the server installation CD; I don't really need (or
want) X windows running in the VM.

To install this onto a new VM, it should be as simple as:

    $ qemu-img create -f qcow2 ubuntu-ppc.qcow2 10G
    $ qemu-system-ppc -m 1024 -hda ubuntu-ppc.qcow2 \
        -cdrom ubuntu-10.04-server-powerpc.iso -boot d

This links up the Ubuntu installation CD on the VM's CD-ROM drive, and
uses a new disk image for the primary hard disk.  Oh, and we make sure
to give the VM enough RAM to do its business — the default is a paltry
128MB.

Of course, this doesn't work — the Lucid installer suffers from the
same problem described
[here](http://mac.linux.be/content/ubuntu-810-installer-fails-detect-cd-rom)
for the Intrepid installer.  Once you get into the installer, the
installation program can't find the CD-ROM device, and so it can't
read the installation packages.  Unfortunately, the workaround doesn't
work for Lucid, since it uses a newer Linux kernel that has
[eliminated the `ide-scsi`
module](http://www.linux.com/archive/feed/33164).

So, what do we do?  Well, QEMU also allows us to mount a disk image as
a USB removable disk, but it won't let us boot from USB.  We end up
having to mount the disk image twice: Once as a virtual CD, so that we
can boot into the installer, and once as a virtual USB disk, so that
the installer can find the installation packages.  The QEMU command
line becomes:

    $ qemu-system-ppc -m 1024 -hda ubuntu-ppc.qcow2 \
        -cdrom ubuntu-10.04-server-powerpc.iso -boot d \
        -usb -usbdevice disk:ubuntu-10.04-server-powerpc.iso

You won't have to manually load the `usb-storage` module; it gets
loaded automatically, and places the USB disk at `/dev/sda`.

You'll still get the error message about not finding the CD; when this
happens say "no" when it asks whether you need to load a module from a
removable disk.  Say "yes" when it asks if you want to choose a module
and device manually; choose "none" for the module; then type in
`/dev/sda` as the device location.


## Corrupt package files on CD

Right, so now we have to be good, right?  We can start QEMU, we can
boot into the installer, and the installer can find all of the
packages?  Nope!  There were several corrupted package files on the CD
image I downloaded.  If this happens to you, you should certainly try
re-downloading the image, to take care of any spurious transmission
errors.  But if you still end up with some corrupted package files,
there are ways around it.

The installer will try to install its initial set of packages using
`apt-get`.  If you encounter problems with these stages, you'll see
some informative error messages on console 4, which is where the
installer's log output is sent.  You can get there by pressing
_Alt-F4_ in the VM.  (As a warning, don't try to shift to console 4
without ensuring that QEMU is grabbing the input.  In most window
managers, _Alt-F4_ will close the current window, which will just
abruptly stop the VM!)

By the time the installer tries to install packages, the VM's hard
disk will be partitioned and formatted, and so we can drop into a
shell as necessary.  To do so, shift over to console 2 using _Alt-F2_
— again, make sure that QEMU is grabbing all keyboard and mouse input
before switching consoles.

Once you're on console 2, you can `chroot` into the new system as
follows:

    ~ $ mount -o /proc /target/proc
    ~ $ mount -o /sys /target/sys
    ~ $ mount -o /dev /target/dev
    ~ $ chroot /target

At this point, you'll be "inside" the new installation system, and can
run whatever `apt-get` and `dpkg` commands are necessary to fix things
up.

Most likely, you'll see "hash sum mismatch" errors, indicating that a
package file is corrupt.  You need to download the correct version
from the archive at _ports.ubuntu.com_.  To do this, you'll need a
copy of _wget_ installed.

    $ apt-get install wget
    $ wget -nv http://ports.ubuntu.com/pool/main/PATH_TO_DEB

You'll see what to use for the `PATH_TO_DEB` part in the error
message.  Once you've downloaded all of the troublesome package files,
install them using:

    $ dpkg -i *.deb
    $ apt-get -f install

Then you can go back into the installer (on console 1) and try to
repeat the current step.

Note that things might be broken early enough that you can't install
_wget_.  If this is the case, how do you download the non-corrupt
package file?  Luckily, Python was already installed at that point, so
you can use the Python standard library to [emulate
_wget_](http://stackoverflow.com/questions/22676/how-do-i-download-a-file-over-http-using-python):

    $ python
    >>> import urllib2
    >>> pkg = urllib2.urlopen("http://ports.ubuntu.com/BLAH_BLAH")
    >>> output = open("BLAH_BLAH.deb", "wb")
    >>> output.write(pkg.read())
    >>> output.close()
    >>> ^D

You can then install the package as above.


## Installing a bootloader

The installer claims that this architecture doesn't support a
bootloader, so we have to install one by hand.  The usual bootloader
for PowerPC machines is `yaboot`; it's fair
