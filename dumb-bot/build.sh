#!/bin/sh

export PROJECT_ID=cloudbowl-356114

pack build --builder=gcr.io/buildpacks/builder gcr.io/$PROJECT_ID/cloudbowl-samples-go-dumb

docker push gcr.io/$PROJECT_ID/cloudbowl-samples-go-dumb

gcloud run deploy --image=gcr.io/$PROJECT_ID/cloudbowl-samples-go-dumb --platform=managed --project=$PROJECT_ID --region=us-central1 --allow-unauthenticated --memory=256Mi cloudbowl-samples-go-dumb
