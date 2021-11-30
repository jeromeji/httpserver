export tag=v1.1
release:
	echo "building httpserver container"
	docker build -t jeromeji/httpserver:${tag}  .

push:
	echo "pushing jeromeji/httpserver"
	docker push jeromeji/httpserver:${tag}
