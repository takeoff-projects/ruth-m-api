gcloud config set project roi-takeoff-user47
gcloud auth activate-service-account level-1-service-account@roi-takeoff-user47.iam.gserviceaccount.com --key-file ~/roi-takeoff-key.json

gcloud builds submit --tag=gcr.io/roi-takeoff-user47/events-api:v1.0 .

terraform init
terraform apply -auto-approve