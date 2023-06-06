.PHONY: all clean deploy gemhtml gemini serve

GEMHTML = .html
GENERATE = generate/generate
GENERATE_SRC = \
    generate/gemtext.go \
    generate/gmnitohtml.go \
    generate/main.go

LIGHTTPD ?= lighttpd

all: gemhtml

clean:
	@rm -f $(GENERATE)
	@rm -rf $(GEMHTML)

deploy:
	@echo Deploying Gemini content...
	@rsync -pgrlvz --delete-after --chmod=ug=rwX,o=rX --groupmap=*:http gemini/ zeta.dcreager.net:/srv/gemini/dcreager.net

gemini:
	@echo Starting Gemini server at gemini://localhost/
	@gmnisrv -C config/gmnisrv.ini

gemhtml: $(GENERATE)
	@echo Generating HTML...
	@$(GENERATE)

$(GENERATE): $(GENERATE_SRC)
	@echo Compiling generator...
	@go build -o $@ ./generate

serve: gemhtml
	@echo Starting server on http://localhost:3000/
	@$(LIGHTTPD) -f lighttpd.conf -i 600
