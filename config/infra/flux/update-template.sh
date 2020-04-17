#!/bin/bash
set -o errexit

FLUX_PATH=$1

if [ -z $FLUX_PATH ]; then
  echo "Usage: $0 <path-to-flux-repo>"
  echo ""
  exit 1
fi


helm template flux $FLUX_PATH/chart/flux --namespace=flux --values=helm_values.yaml > apply.yaml