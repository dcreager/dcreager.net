exec >&2
redo-ifchange generate/generate .sources
cat .sources | xargs -d '\n' redo-ifchange
echo Generating HTML...
generate/generate
