build:

	go build github.com/242617/torture

run: build

	./torture -config config/torture.yaml

docker-build:

	docker build -t 242617/torture src/github.com/242617/torture

docker-run: docker-build

	docker run -it --rm --name torture -p 8080:8080 242617/torture

docker-push: docker-build

	docker push 242617/torture