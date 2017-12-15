.PHONY: build
build:
	go build -o out/kmake ./cmd/kmake

.PHONY: example
example:
	./out/kmake watch --dockerfile ./examples/Dockerfile --image-name hello-node

.PHONY: remote
remote:
	./out/kmake watch --dockerfile ./examples/Dockerfile --image-name hello-node --project-id=r2d4minikube --remote=true
