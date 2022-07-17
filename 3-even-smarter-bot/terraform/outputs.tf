output "container_repository_url" {
  value = "${google_artifact_registry_repository.main.location}-docker.pkg.dev/${google_project.main.project_id}/${google_artifact_registry_repository.main.repository_id}"
}

output "cloudrun_service_account_id" {
  value = google_service_account.cloudrun_service_sa.email
}
