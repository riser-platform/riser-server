#!/bin/bash
set -o errexit

ISTIO_PATH=$1

if [ -z $ISTIO_PATH ]; then
  echo "Usage: $0 <path-to-istio-repo>"
  echo ""
  exit 1
fi


helm template $ISTIO_PATH/install/kubernetes/helm/istio-init/ --namespace istio-system --values helm_values.yaml > apply/1_init.yaml
helm template $ISTIO_PATH/install/kubernetes/helm/istio/ --namespace istio-system --values helm_values.yaml > apply/2_helm_template.yaml
