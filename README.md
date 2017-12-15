# kmake

build

```
	go build -o out/kmake ./cmd/kmake
```


## local build

requirements:
* docker cli
* minikube
* ksonnet

set minikube as the local docker daemon

```
$ eval $(minikube docker-env)
```

run kmake
```
./out/kmake watch --dockerfile ./examples/Dockerfile --image-name hello-node
```

make some changes to `examples/server.js`

see the changes with
```
$ minikube service hello-node
```

## cloud build

requirements:
* container builder api enabled
* minikube
* ksonnet

follow the prompts to mount GCR credentials in minikube
```
$ minikube addons configure registry-creds
$ minikube start
$ minikube addons enable registry-creds
```

start kmake with the `--remote=true` and `--project-id` flags set
```
./out/kmake watch --dockerfile ./examples/Dockerfile --image-name hello-node --remote=true --project-id=r2d4minikube
```

make some changes to `examples/server.js`

see the changes with
```
$ minikube service hello-node
```
