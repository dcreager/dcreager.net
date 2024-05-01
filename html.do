exec >&2
redo-ifchange generate/generate templates/*.html
find $HOME/notes -type f -print0 | xargs -0 redo-ifchange
find overrides -type f -print0 | xargs -0 redo-ifchange
echo Generating HTML...
generate/generate
