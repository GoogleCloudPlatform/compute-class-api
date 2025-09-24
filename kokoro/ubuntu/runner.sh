#!/bin/bash

set -e
set -x

gcloud auth configure-docker us-central1-docker.pkg.dev --quiet

# KOKORO_ROOT = "/tmpfs"
mkdir -p $KOKORO_ROOT/gopath/src/github.com/googlecloudplatform
export GOPATH=$KOKORO_ROOT/gopath
BUILDPATH=$KOKORO_ROOT/gopath/src/github.com/googlecloudplatform/compute-class-api
sudo cp -r git/compute-class-api $BUILDPATH
cd $BUILDPATH

if [[ -z "${HOST_INPUT_ROOT:-}" ]]; then
  export GITROOT=$BUILDPATH
else
  export GITROOT="$HOST_INPUT_ROOT/src/git/compute-class-api"
fi

export PATH=$GOPATH/bin:$PATH

make build-in-docker
make test-in-docker
