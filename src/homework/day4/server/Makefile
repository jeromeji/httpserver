export tag=v1.0
release:
	echo "building httpserver container"
	docker build -t jeromeji/httpserver:${tag} -f ./Dockerfile

push:
	echo "pushing jeromeji/httpserver"
	docker push jeromeji/httpserver:v1.0
