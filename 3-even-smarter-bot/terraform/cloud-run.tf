resource "google_service_account" "cloudrun_service_sa" {
  account_id   = "cloudrun-bot"
  display_name = "CloudRun Bot Account"
  project      = google_project.main.project_id
}

resource "google_project_iam_member" "cloud_run_sa" {
  for_each = toset([
    "roles/monitoring.editor",
    "roles/pubsub.publisher"
  ])

  project = google_project.main.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.cloudrun_service_sa.email}"

  #   depends_on = [
  #     google_project_service.main // need to wait for the cloud build api to be enabled on the project otherwise the SA won't exist
  #   ]
}
