exec >&2
redo-ifchange generate/generate .sources
cat .sources | xargs -d '\n' redo-ifchange
echo Generating HTML...
generate/generate

mkdir -p .html/fonts/woff2
redo-ifchange \
    .html/fonts/woff2/iosevka-ss18-bold.gemlink.woff2 \
    .html/fonts/woff2/iosevka-term-ss18-regular.subset.woff2
