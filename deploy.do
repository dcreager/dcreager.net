exec >&2
redo-ifchange html
echo Deploying HTML content...
rsync -pgrlz --info=COPY,DEL --delete-after --chmod=ug=rwX,o=rX --groupmap=*:http .html/ zeta.dcreager.net:/srv/http/notes.dcreager.net
