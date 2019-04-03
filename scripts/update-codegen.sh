#!/bin/bash

PROJECT_DIR=$(dirname ${BASH_SOURCE})/..
PROJECT_ROOT=$(readlink -f ${PROJECT_DIR})
CODEGEN_PKG=${PROJECT_DIR}/vendor/k8s.io/code-generator

PKG=${PROJECT_ROOT#"${GOPATH}/src/"}

${CODEGEN_PKG}/generate-groups.sh \
    "all" \
    ${PKG}/pkg/client \
    ${PKG}/pkg/apis \
    owner:v1alpha1
