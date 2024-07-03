exec >&2

redo-ifchange html

cat >.lighttpd.conf <<CONF
server.document-root = var.CWD + "/.html"
server.port = 3000
index-file.names = ( "index.html" )

mimetype.assign = (
  ".html" => "text/html; charset=utf-8",
  ".xml" => "text/xml; charset=utf-8",
  ".css" => "text/css; charset=utf-8",
  ".svg" => "image/svg",
  ".jpg" => "image/jpeg",
  ".png" => "image/png"
)
CONF

echo Starting server on http://localhost:3000/
(
  # Start lighttpd in a subshell where we've closed all file descriptors except
  # for std{in,out,err}.  Redo uses the additional file descriptors to wait for
  # child processes to exit, and we don't want it to consider this background
  # job when doing so.
  FDS=$(find /proc/self/fd -type l -printf '%f\n')
  for n in $FDS; do if ((n > 2)); then eval "exec $n>&-"; fi; done
  lighttpd -f $PWD/.lighttpd.conf -i 600
)
