resource "google_project" "main" {
  name            = "microbot hackathon"
  project_id      = var.project_id
  org_id          = var.org_id
  billing_account = var.billing_account
}

resource "google_project_service" "main" {
  for_each = toset([
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "orgpolicy.googleapis.com",
    "run.googleapis.com"
  ])

  service                    = each.value
  project                    = google_project.main.project_id
  disable_dependent_services = false
  disable_on_destroy         = false
}
