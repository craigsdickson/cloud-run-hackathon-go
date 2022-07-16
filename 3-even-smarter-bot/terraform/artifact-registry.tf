resource "google_artifact_registry_repository" "main" {
  provider = google-beta

  project       = google_project.main.project_id
  location      = var.repository_location
  repository_id = var.repository_id
  format        = "DOCKER"

  depends_on = [
    google_project_service.main // need to wait for the artifact registry api to be enabled on the project
  ]
}
