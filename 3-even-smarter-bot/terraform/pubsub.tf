resource "google_pubsub_topic" "main" {
  name                       = "arena-updates"
  project                    = google_project.main.project_id
  message_retention_duration = "600s"
}
