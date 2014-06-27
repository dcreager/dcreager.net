# Blog post ideas

- Linked Data.  Lots of people who are using Linked Data go whole hog, and try
  to use RDF as an internal data model.  You can also use Linked Data only when
  publishing data (esp through a REST API) — either to add LD support to an
  existing application with existing internal data formats and structures, or
  because something else (JSON-style, relational, etc) might be easier to work
  with internally.  How can you do that as easily as possible?  JSON-LD uses the
  right approach — use a context or schema to define the transformation from an
  internal format to LD (serialized however you want).
