build:

	go build github.com/242617/torture

run: build

	./torture -config config/torture.yaml