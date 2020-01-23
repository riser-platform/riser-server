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

# Adds a cluster-local gateway
helm template --namespace=istio-system \
  --set gateways.custom-gateway.autoscaleMin=1 \
  --set gateways.custom-gateway.autoscaleMax=2 \
  --set gateways.custom-gateway.cpu.targetAverageUtilization=60 \
  --set gateways.custom-gateway.labels.app='cluster-local-gateway' \
  --set gateways.custom-gateway.labels.istio='cluster-local-gateway' \
  --set gateways.custom-gateway.type='ClusterIP' \
  --set gateways.istio-ingressgateway.enabled=false \
  --set gateways.istio-egressgateway.enabled=false \
  --set gateways.istio-ilbgateway.enabled=false \
  $ISTIO_PATH/install/kubernetes/helm/istio \
  -f $ISTIO_PATH/install/kubernetes/helm/istio/example-values/values-istio-gateways.yaml \
  | sed -e "s/custom-gateway/cluster-local-gateway/g" -e "s/customgateway/clusterlocalgateway/g" \
  > ./apply/2_local_gateway.yaml

