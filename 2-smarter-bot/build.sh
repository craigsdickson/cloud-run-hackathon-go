#!/bin/sh

export PROJECT_ID=cloudbowl-356114

pack build --builder=gcr.io/buildpacks/builder gcr.io/$PROJECT_ID/cloudbowl-samples-go-smart

docker push gcr.io/$PROJECT_ID/cloudbowl-samples-go-smart

gcloud run deploy --image=gcr.io/$PROJECT_ID/cloudbowl-samples-go-smart --platform=managed --project=$PROJECT_ID --region=us-central1 --allow-unauthenticated --memory=256Mi cloudbowl-samples-go-smart
