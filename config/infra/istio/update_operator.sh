#!/bin/sh
set -o errexit

ISTIO_VERSION=$1
ISTIO_REPO_PATH=$2

if [ -z $ISTIO_VERSION ] || [ -z $ISTIO_REPO_PATH ]; then
  echo "$0 (istio_version) (/path/to/istio/repo)"
  echo ""
  exit 1
fi

helm template $ISTIO_REPO_PATH/manifests/charts/istio-operator/ \
  --set hub=docker.io/istio \
  --set tag=$ISTIO_VERSION \
  --set operatorNamespace=istio-operator \
  --set istioNamespace=istio-system > istio_operator.yaml

