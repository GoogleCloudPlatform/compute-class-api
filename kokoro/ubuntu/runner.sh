#!/bin/bash
#!/usr/bin/env bash
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
