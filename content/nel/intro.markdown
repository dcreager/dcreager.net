---
draft: true
kind: article
created_at: 2018-03-26
layout: /post.html
title: "Introducing Network Error Logging"
description: >-
  in which we get the browser to do our work for us
tags: [nel]
---

<style>
.log {
  border-radius: 0;
  font-size: 11px;
  padding: 6px;
}

.actor {
  width: 10em;
  margin-right: 20px;
  float: left;
}
</style>

Let's say you've got a web site or REST service.  You have clients that send
requests to your server, which performs some interesting processing and sends
responses back.  The clients might be people using a browser, or native code in
a mobile app making requests on behalf of the user, or something more exotic.
But the overall picture is the same:

![request service flow](service-flow.png){:.figure}

How do you monitor this service — especially when that cloudy bit in the middle
is completely out of your control?

Server logs are a great source of information, and are typically part of the
answer.  Your server creates a record for every single connection attempt that
it receives, including information about the client, the resource they
requested, and whether or not the request succeeded:

![server](server.png){:.actor}

<pre class="log">
192.0.2.1 - - [26/Mar/2018:03:04:05.123 +0000] "GET /about/" 200 -
192.0.2.1 - - [26/Mar/2018:03:04:17.456 +0000] "GET /login/" 200 -
198.51.100.13 - - [26/Mar/2018:03:04:24.924 +0000] "GET /unknown/" 402 -
198.51.100.13 - - [26/Mar/2018:03:04:38.129 +0000] "GET /" 200 -
198.51.100.13 - - [26/Mar/2018:03:05:08.751 +0000] "GET /about/" 200 -
192.0.2.1 - - [26/Mar/2018:03:05:13.314 +0000] "GET /contact/" 200 -
192.0.2.1 - - [26/Mar/2018:03:05:48.624 +0000] "GET /directions/" 200 -
203.0.113.56 - - [26/Mar/2018:03:08:08.242 +0000] "GET /" 200 -
203.0.113.56 - - [26/Mar/2018:03:08:52.018 +0000] "GET /about/" 200 -
203.0.113.56 - - [26/Mar/2018:03:09:03.483 +0000] "GET /contact/" 200 -
203.0.113.56 - - [26/Mar/2018:03:10:32.851 +0000] "GET /logout/" 500 -
192.0.2.1 - - [26/Mar/2018:03:10:50.823 +0000] "GET /menu/" 200 -
192.0.2.1 - - [26/Mar/2018:03:12:03.802 +0000] "GET /specials/" 200 -
198.51.100.13 - - [26/Mar/2018:03:12:24.516 +0000] "GET /login/" 200 -
</pre>

There are a couple of interesting wrinkles.  First, how do you get at those
server logs?  Your web server will typically write those to disk on the machine
that it's running on, so you need to somehow copy the contents of the logs to
whatever monitoring or observability platform you're using.  (When you're first
starting off, you can get away with manually logging into the servers and
grepping through the files directly, but that doesn't scale as you add machines
or developers.)  This isn't an unsolvable problem, but it does add some
operational complexity to your setup.

More interesting is that your server logs don't have visibility into all of the
problems that your clients might encounter — especially those that involve the
network!  Since the network isn't under your control, you have to consider it a
black box: your client puts requests into one side, and hopefully those requests
come out the other side at your server.

To have an end-to-end view of your request traffic, you need to instrument the
client side of the connection, too.  Have it record every time it attempts a
request:

![client](client.png){:.actor}

<pre class="log">
2018-03-26 03:04:04 https://example.com/about/ ⇒ 200
2018-03-26 03:04:16 https://example.com/login/ ⇒ 200
2018-03-26 03:05:13 https://example.com/contact/ ⇒ 200
2018-03-26 03:05:47 https://example.com/directions/ ⇒ 200
2018-03-26 03:10:50 https://example.com/menu/ ⇒ 200
2018-03-26 03:12:02 https://example.com/specials/ ⇒ 200
</pre>

Then send those records to the same analysis platform that you send the server
logs to.  If there's a record that appears on the client side, but not on the
server side, then you've found an interesting case where the network dropped
your connection!  (The next step is to figure out **why** the connection was
dropped, which is a topic worthy of its own post.)

When you have control over the client code, this is pretty doable.  Most of the
libraries that you'd use to make those HTTP requests provide all of the
necessary hooks to add in this kind of instrumentation.  If you're trying to
monitor a web site that your users access directly from a browser, it's a bit
more complicated; but with the right JavaScript (and maybe a smattering of
service workers), you can hook into and monitor the browser's requests to your
site, too.

Of course, even if it's easy to spackle all of the pieces together, that's one
more bit of operational coding that takes you away from working on your actual
service.  Wouldn't it be nice if the browser or HTTP library would do all of
this client-side monitoring work for you?

Enter [Network Error Logging][NEL] (NEL), a new web platform spec that we're
working on in the W3C's [Web Performance Working Group][WebPerf].  With NEL, you
can instruct user agents to collect the same set of information that would
appear in your server logs.  But since the data is being collected directly by
the user agent, you should have visibility into **all** of your clients'
connection requests, not just those that made it to your serving infrastructure.

[NEL]: https://wicg.github.io/network-error-logging/
[WebPerf]: https://www.w3.org/webperf/

NEL is still in its early stages.  We have what we think is a good initial
draft, and we're wrapping up proof-of-concept implementations (both in the
browser and in a handful of HTTP request libraries).  But we'd love to hear
feedback.  Do you have any interesting outages that NEL would've given you
better visibility into?  Is there anything you think we've missed?  Let me know
through [the usual channels][contact] or by filing an issue on [our Github
repo][NEL github].

[contact]: /about/
[NEL github]: https://github.com/wicg/network-error-logging
