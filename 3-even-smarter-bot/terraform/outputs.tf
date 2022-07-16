output "container_repository_url" {
  value = "${google_artifact_registry_repository.main.location}-docker.pkg.dev/${google_project.main.project_id}/${google_artifact_registry_repository.main.repository_id}"
}
