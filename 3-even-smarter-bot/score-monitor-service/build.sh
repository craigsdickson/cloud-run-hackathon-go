#!/bin/sh

PROJECT_ID=microbot-hackathon

DIR=`pwd`

cd ../..

# the user account running this gcloud command needs the roles/viewer or roles/owner role on the project to be able to stream the cloud build logs to the console

gcloud builds submit --project $PROJECT_ID --config $DIR/cloudbuild.yaml --substitutions=COMMIT_SHA=latest --ignore-file=$DIR/.gcloudignore .

cd $DIR
