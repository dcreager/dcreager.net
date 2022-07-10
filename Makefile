.PHONY: all deploy force gemini http

METADATA = \
    .posts.txt \
    .posts-20.txt

FILES = \
    .output/index.html \
    .output/talks/index.html

POSTS = \
    .output/np-hard/index.html

POST_TEMPLATES = \
    templates/foot.html \
    templates/js.html \
    templates/links.html \
    templates/post.html \
    templates/title.html \
    templates/topnav.html

all: $(METADATA) $(FILES) $(POSTS)

deploy:
	@echo Deploying Gemini content...
	@rsync -pgrlvz --delete-after --chmod=ug=rwX,o=rX --groupmap=*:http gemini/ zeta.dcreager.net:/srv/gemini

gemini:
	@echo Starting Gemini server at gemini://localhost/
	@gmnisrv -C config/gmnisrv.ini

http:
	@simple-http-server --index .output

.output/index.html: content/index.html
	@echo $@ $<

.output/%/index.html: content/%.html
	@echo $@ $<

.output/%/index.html: content/%.markdown $(POST_TEMPLATES)
	@echo $@ $<
	pandoc -t html5 -s --template templates/post.html -o $@ $<

# This writes the contents of $(POSTS) into .posts.txt, but only if the value
# has changed.  That allows us to write rules that depend on .posts.txt, and
# they will only be rebuilt when $(POSTS) has a new value.
# https://stackoverflow.com/a/3237349
.posts.txt: force
	@echo '$(POSTS)' | cmp -s - $@ || echo '$(POSTS)' > $@

.posts-20.txt: .posts.txt
	@head -n 20 $< > $@
