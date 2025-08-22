#!/usr/bin/env bash

ROOT="$(git rev-parse --show-toplevel)"

if [ ! -d "$ROOT" ]; then
  echo "Unable to find ${ROOT} folder"
  exit 1
fi

command -v controller-gen >/dev/null 2>&1 || { echo >&2 "controller-gen should be available in PATH, consider installing using 'make gen-clients-init'"; exit 1; }

# Fetch code-generator repository
TMP_DIR=$(mktemp -d)
CODEGEN_PKG="${TMP_DIR}/code-generator"
CODEGEN_REPO="https://github.com/kubernetes/code-generator.git"
CODEGEN_VERSION="kubernetes-1.30.3"

git clone --branch "${CODEGEN_VERSION}" --single-branch "${CODEGEN_REPO}" "${CODEGEN_PKG}"

. "${CODEGEN_PKG}/kube_codegen.sh"

# Generate client
BOILERPLATE_FILE="${ROOT}/hack/boilerplate.go.txt"
API_PATH="${ROOT}/api"
CLIENT_PATH="${ROOT}/client"

OUTPUT_PKG="github.com/googlecloudplatform/compute-class-api"
OUTPUT_CLIENT_PKG="${OUTPUT_PKG}/client"

echo "INFO: Cleaning client directory..."
rm -rf "${CLIENT_PATH}"

echo "INFO: Generating helpers..."
kube::codegen::gen_helpers "${API_PATH}" --boilerplate "${BOILERPLATE_FILE}"

echo "INFO: Generating clients..."
kube::codegen::gen_client "${API_PATH}" \
    --output-dir "${CLIENT_PATH}" \
    --output-pkg "${OUTPUT_CLIENT_PKG}" \
    --with-applyconfig \
    --with-watch \
    --boilerplate "${BOILERPLATE_FILE}"

echo "INFO: Generating CRD manifest..."
controller-gen \
  crd:generateEmbeddedObjectMeta=true \
  paths="${API_PATH}/..." \
  output:crd:dir="${ROOT}"

echo "INFO: Running go mod tidy..."
go mod tidy
