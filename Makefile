.PHONY: all deploy gemini

all:

deploy:
	@echo Deploying Gemini content...
	@rsync -pgrlzDP --delete-after --chmod=ug=rwX,o=rX --groupmap=*:http gemini/ zeta.dcreager.net:/srv/gemini

gemini: .keys/created
	@echo Starting Gemini server at gemini://localhost/
	@AGATE_LOG=info agate 127.0.0.1:1965 gemini .keys/cert.pem .keys/key.rsa

.keys/created:
	@echo Generating TLS keys...
	@openssl req -x509 -newkey rsa:4096 -keyout .keys/key.rsa -out .keys/cert.pem -days 3650 -nodes -subj "/CN=localhost"
	@touch $@
