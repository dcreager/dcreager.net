exec >&2
SRC=fonts/"$(basename $2)".woff2
redo-ifchange "$SRC"
pyftsubset "$SRC" --unicodes=21d2,219d --output-file=$3 --flavor=woff2
