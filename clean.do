for d in */clean.do; do
    echo "${d%.do}"
done | xargs redo

rm -rf .html
