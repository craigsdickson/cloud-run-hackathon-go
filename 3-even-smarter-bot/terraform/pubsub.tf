resource "google_pubsub_topic" "main" {
  name                       = "arena-updates"
  project                    = google_project.main.project_id
  message_retention_duration = "600s"
}

resource "google_project_iam_member" "pubsub_sa_account" {
  project = google_project.main.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:service-${google_project.main.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}
