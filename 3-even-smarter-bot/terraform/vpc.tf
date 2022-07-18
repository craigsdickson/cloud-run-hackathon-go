resource "google_compute_network" "main" {
  name                    = "my-network"
  project                 = google_project.main.project_id
  auto_create_subnetworks = true
}

resource "google_vpc_access_connector" "connector" {
  name          = "vpcconn"
  provider      = google-beta
  region        = "us-central1"
  ip_cidr_range = "10.8.0.0/28"
  # max_throughput = 300
  network = google_compute_network.main.name
  project = google_project.main.project_id

  depends_on = [
    google_project_service.main // need to wait for the vpcaccess api to be enabled on the project
  ]
}
