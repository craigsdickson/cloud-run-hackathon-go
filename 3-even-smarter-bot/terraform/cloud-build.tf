// the cloud build SA is created when the API is enabled, but not given any default permissions
resource "google_project_iam_member" "cloud_build_sa_builder" {
  project = google_project.main.project_id
  role    = "roles/cloudbuild.builds.builder"
  member  = "serviceAccount:${google_project.main.number}@cloudbuild.gserviceaccount.com"

  depends_on = [
    google_project_service.main // need to wait for the cloud build api to be enabled on the project otherwise the SA won't exist
  ]
}

// cloud build SA needs Cloud Run Admin role to deploy to Cloud Run
resource "google_project_iam_member" "cloud_build_sa_runadmin" {
  project = google_project.main.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_project.main.number}@cloudbuild.gserviceaccount.com"

  depends_on = [
    google_project_service.main // need to wait for the cloud build api to be enabled on the project otherwise the SA won't exist
  ]
}

resource "google_project_iam_member" "cloud_build_sa_sauser" {
  project = google_project.main.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_project.main.number}@cloudbuild.gserviceaccount.com"

  depends_on = [
    google_project_service.main // need to wait for the cloud build api to be enabled on the project otherwise the SA won't exist
  ]
}

resource "google_cloudbuild_trigger" "main" {
  name        = "deploy-to-cloudrun"
  description = "Builds bot container and deploys to cloudrun"
  filename    = "cloudbuild.yaml"
  included_files = [
    "3-even-smarter-tool/go/**"
  ]
  github {
    name  = "cloud-run-hackathon-go"
    owner = "craigsdickson"
    push {
      branch = "^main$"
    }
  }
  #   trigger_template {
  #     branch_name = "main"
  #     repo_name   = "my-repo"
  #   }
  project = google_project.main.project_id

  #   substitutions = {
  #     _FOO = "bar"
  #     _BAZ = "qux"
  #   }


}
