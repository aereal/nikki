provider "google" {
  project = "nikki-develop"
  region = "asia-northeast1"
}

resource "google_compute_disk" "nikki-db-disk" {
  name = "nikki-db-disk"
  zone = "asia-northeast1-a"
  image = "debian-9-stretch-v20171213"
  type = "pd-ssd"
}

resource "google_compute_instance" "nikki-db" {
  name = "nikki-db"
  machine_type = "f1-micro"
  zone = "asia-northeast1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"

    access_config {
    }
  }

  attached_disk {
    source = "nikki-db-disk"
  }
}
