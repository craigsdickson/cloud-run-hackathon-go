#!/bin/sh

export PROJECT_ID=microbot-hackathon

pack build --builder=gcr.io/buildpacks/builder gcr.io/$PROJECT_ID/dumb-bot

docker push gcr.io/$PROJECT_ID/dumb-bot

gcloud run deploy --image=gcr.io/$PROJECT_ID/dumb-bot --platform=managed --project=$PROJECT_ID --region=us-central1 --allow-unauthenticated --memory=256Mi dumb-bot
