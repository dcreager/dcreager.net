find . -name '*.go' -print0 | xargs -0 redo-ifchange
go build -o $3 ./cmd
