exec >&2
SRC=fonts/"$(basename $2)".woff2
redo-ifchange "$SRC" .content.all
pyftsubset "$SRC" --text-file=.content.all --output-file=$3 --flavor=woff2
