.PHONY: gen-client-init
gen-clients-init:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0

.PHONY: gen-client
gen-client:
	bash hack/gen-client.sh

# Variables
GITROOT=$(if $(HOST_INPUT_ROOT),$(HOST_INPUT_ROOT)/src/git/compute-class-api,$(shell pwd))

docker-builder:
	docker build -t compute-class-api-builder ./builder

build-in-docker: docker-builder
	docker run -v $(GITROOT):/tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api compute-class-api-builder:latest bash -c 'git config --global --add safe.directory /tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api && cd /tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api && go build ./...'

test-in-docker: docker-builder
	docker run -v $(GITROOT):/tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api compute-class-api-builder:latest bash -c 'git config --global --add safe.directory /tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api && cd /tmpfs/gopath/src/github.com/googlecloudplatform/compute-class-api && go test -v ./api/cloud.google.com/v1'
