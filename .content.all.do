redo-ifchange .sources
cat .sources | grep -E '\.(css|gmi|html)$' | xargs -d '\n' cat >$3
