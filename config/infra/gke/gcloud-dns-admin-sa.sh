# Creates a service account named "dns-admin" and creates a service JSON to be used for
# a k8s secret. Must be created in the same ns as cert manager
# kubectl create secret generic cert-manager-credentials \
#    --from-file=./gcp-dns-admin.json
# See also: https://medium.com/google-cloud/https-with-cert-manager-on-gke-49a70985d99b
set -o errexit

GCP_PROJECT=$1

if [ -z $GCP_PROJECT ]; then
    echo "Usage $0 <project>"
    exit 1
fi

gcloud iam service-accounts create dns-admin \
    --display-name=dns-admin \
    --project=${GCP_PROJECT}

gcloud iam service-accounts keys create ./gcp-dns-admin.json \
    --iam-account=dns-admin@${GCP_PROJECT}.iam.gserviceaccount.com \
    --project=${GCP_PROJECT}

gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
    --member=serviceAccount:dns-admin@${GCP_PROJECT}.iam.gserviceaccount.com \
    --role=roles/dns.admin