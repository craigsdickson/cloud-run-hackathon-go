#!/bin/sh

# the user account running this gcloud command needs the roles/viewer or roles/owner role on the project to be able to stream the cloud build logs to the console

gcloud builds submit --project microbot-hackathon --config ./3-even-smarter-bot/go/cloudbuild.yaml --substitutions=COMMIT_SHA=latest ./3-even-smarter-bot/go
