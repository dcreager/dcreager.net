---
layout: post
title: Decentralized datatypes
---

Over the past year or so there have been quite a few blog postings in
the REST world about MIME types, and their role in the REST
architecture.  A lot of the discussion seems to be prompted by WADL,
which is an attempt to define a WSDL-style interface description
language for REST services.  [Joe
Gregorio](http://bitworking.org/news/193/Do-we-need-WADL) argues that
MIME types are more useful for describing the semantics of a service
than a WADL document, since there are parts of the service's semantics
that just can't be encoded in a machine-readable format.  MIME types
acknowledge this, providing a standard way of identifying a data
format and pointing to the human- and machine-readable documents (such
as RFCs and XSDs) that define the syntax and accompanying semantics.

Following this idea, several people have begun debating whether or not
the centralized assignment of MIME types is the right way to handle
the variety of data formats that REST-based systems produce and
consume.  [Mark
Baker](http://www.markbaker.ca/blog/2008/02/media-type-centralization-is-a-feature-not-a-bug/)
comes in on the side of centralized assignment, whereas [Stefan
Tilkov](http://www.innoq.com/blog/st/2008/02/decentralizing_media_types.html),
[Dan
Diephouse](http://netzooid.com/blog/2008/02/07/why-a-restful-idl-is-an-oxymoron-and-what-we-really-need-instead/),
and [James
Strachan](http://macstrac.blogspot.com/2007/11/atompub-services-and-auto-detecting.html)
argue in favor of decentralized types.  [Bill
Burke](http://bill.burkecentral.com/2008/03/05/restful-xml-content-negotitation/)
and [Benjamin
Carlyle](http://soundadvice.id.au/blog/2009/08/16/#mimeLimitation)
have good summaries of the different proposed technical solutions that
would enable decentralized types.

## “Extended” types

One of the arguments in favor of centralized assignment is that
allowing everyone to invent their own MIME types would ruin
interoperability.  And for certain cases, this seems pretty obviously
true.  It's a good thing that we have a standardized `image/png` MIME
type; this allows your browser to correctly display the website logo
you see up in the upper left corner.  If I were daft, I could decide
to serve that logo using a MIME type of `image/x-dcreager-png` (or
similar) to indicate that I've included some particular set of
metadata in an ancillary chunk of the PNG.

Why would I want to do this?  Maybe I'm writing an application that
knows how to process this metadata, and I'd like to easily determine
whether a particular resource I'm accessing has this metadata or not.
The [Open Microscopy Environment](http://www.openmicroscopy.org) does
exactly this; they've defined an XML schema that allows biology
researchers to provide additional scientific metadata about an image
or movie captured from a high-end microscope.  One way to encode an
image and its metadata is as an “OME-TIFF”, a data format that
includes the metadata in an optional TIFF section.  OME-TIFF files are
also perfectly valid as regular TIFF files.  This has the benefit that
an OME-aware application can access the scientific metadata, whereas a
“regular” image processing application can read the image using its
normal TIFF decoder.

Of course, now we have competing goals that we have to reconcile.  On
the one hand, we need to ensure that OME-aware applications can see
that a particular image is an OME-TIFF.  On the other, we need
non-OME-aware applications to see the image as a regular TIFF.  One of
the decentralized proposals — MIME type parameters — tries to address
this.  For instance, a MIME type for an OME-TIFF might be `image/tiff;
ome=xml`.  By using the standard `image/tiff` as the base MIME type,
non-OME-aware applications correctly treat it as a simple TIFF.
OME-aware applications would know that the `ome=xml` parameter
indicates that the OME-specific metadata is present.

## The multitude of XML types

Another example given is that of an XML document.  Most applications
will generate XML documents that conform to a particular schema (for
instance, a company-specific purchase order), which they might encode
as an XSD.  Now, the XSD on its own doesn't give you the full story on
how to process that data, but it does provide some detail on how the
data is structured.  If you're writing an application that consumes
this data, having the XSD available would be helpful.  More
interesting is an application that can consume _any_ XML document —
and which might use an XSD or RelaxNG schema to customize the UI used
to display the document.

In both cases, the schema is necessary to process the document, but
for different reasons.  In the first case, the consuming application
was built with advance knowledge of how the data should be handled,
and the schema is used to direct a particular document to the code
that implements this knowledge.  In the second case, the particular
datatype is unimportant, and the application-specific semantics aren't
used; the data is only consumed as a “generic XML document”, and the
schema is used to describe the specific structure of the elements.

## Data doesn't have a single type

The common theme in both of these examples is that a single datatype
isn't enough to describe the data we're dealing with.  As Roy Fielding
[points out](http://roy.gbiv.com/untangled/2009/wrangling-mimetypes),
“all data formats correspond to multiple media types”.  It's tempting
to think of a datatype as just “the syntax and structure of the data”.
But it must also include some intuition about how the data will be
used.

From this point of view, the generic XML processing application does
_not_ handle a multitude of datatypes.  Instead, it handles exactly
one: “generic XML document with associated schema”.  The application
that knows how to process this particular schema will handle a
different, completely distinct, datatype: “company-specific purchase
order XML document”.  And the particular XML document in question — a
single sequence of bytes that is a single representation of a single
resource — is an instance of both types.

Why shift things around like this?  Doesn't it just move the
complexity from the consumer (who used to consume multiple types) to
the producer (who must now publish the XML document under different
types)?  Not necessarily.  The key idea is that we can use
_transformation graphs_ to encode the relationships between the
datatypes:

<div class="figure">
  <img src="/images/2009/12/21/decentralized-datatypes/xform-graph.png" alt="transformation graph"/>
</div>

In this specific example, the transformation is simple — since the
same sequence of bytes is a valid instance of both types, we don't
have to modify the data itself.  The decentralized MIME types
(especially the MIME parameter proposal) already support these kinds
of “no-op” transformations: the more generic type is the “base” MIME
type, and the more specific extensions are encoded as MIME parameters.
However, by modeling the type relationships as an arbitrary graph, we
open up the possibility of more complex sets of types, which might
require actual code to transform between them, but which can be
defined in a decentralized manner.

Even though the model is more complex, we haven't required the
producer or consumer to be more complex.  A transformation graph is
necessary to translate between the two different (but compatible)
types, but the graph doesn't have to specifically live at the producer
or the consumer.  The producer can publish the data using the only
type it knows about (the “company-specific” type), and a consumer can
request the data using the only type it knows about (such as the
“generic XML” type).  Anywhere along the path from the producer to the
consumer, we can use the transformation graph to automatically
transform the data from one type to the other.

More details on transformation graphs, including more complex
examples, can be found in my [DPhil
thesis](/publications/012-dphil-thesis).
