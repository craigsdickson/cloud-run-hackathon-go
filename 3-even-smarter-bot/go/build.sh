#!/bin/sh

export PROJECT_ID=cloudbowl-356114

pack build --builder=gcr.io/buildpacks/builder us-docker.pkg.dev/$PROJECT_ID/bots-repo/even-smarter-bot

docker push us-docker.pkg.dev/$PROJECT_ID/bots-repo/even-smarter-bot

gcloud run deploy --image=us-docker.pkg.dev/$PROJECT_ID/bots-repo/even-smarter-bot --platform=managed --project=$PROJECT_ID --region=us-central1 --allow-unauthenticated --memory=256Mi even-smarter-bot
