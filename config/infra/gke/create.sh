# Creates a GKE cluster in us-east1 using recommended settings for demoing Riser
PROJECT=$1

if [ -z $PROJECT ]; then
  echo "Usage: $0 (project)"
  exit 1
fi

gcloud beta container clusters create "riser-demo" \
 --project "$PROJECT" \
 --zone "us-east1-b" \
 --no-enable-basic-auth \
 --cluster-version "1.13.7-gke.24" \
 --machine-type "n1-standard-1" \
 --image-type "COS" \
 --disk-type "pd-standard" \
 --disk-size "10" \
 --metadata disable-legacy-endpoints=true \
 --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" \
 --num-nodes "3" \
 --enable-ip-alias \
 --network "projects/$PROJECT/global/networks/default" \
 --subnetwork "projects/$PROJECT/regions/us-east1/subnetworks/default" \
 --default-max-pods-per-node "110" \
 --addons HttpLoadBalancing \
 --enable-autorepair
