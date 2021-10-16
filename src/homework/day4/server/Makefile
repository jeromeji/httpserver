export tag=v1.0
release: build
	echo "building httpserver container"
	docker build -t jeromeji/httpserver:${tag} -f ./Dockerfile

push: release
	echo "pushing jeromeji/httpserver"
	docker push jeromeji/httpserver:v1.0
