resource "google_redis_instance" "main" {
  name               = "leaderboard"
  memory_size_gb     = 1
  project            = google_project.main.project_id
  region             = var.region
  authorized_network = google_compute_network.main.id

  depends_on = [
    google_project_service.main // need to wait for the redis api to be enabled on the project
  ]
}
