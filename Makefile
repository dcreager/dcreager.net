.PHONY: all deploy gemini

all:

deploy:
	@echo Deploying Gemini content...
	@rsync -pgrlvz --delete-after --chmod=ug=rwX,o=rX --groupmap=*:http gemini/ zeta.dcreager.net:/srv/gemini/dcreager.net

gemini:
	@echo Starting Gemini server at gemini://localhost/
	@gmnisrv -C config/gmnisrv.ini
