resource "google_cloud_run_v2_service" "httpbin" {
  name     = var.service_name
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  scaling {
    manual_instance_count = null
    scaling_mode = "AUTOMATIC"
  }

  template {
    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    containers {
      image = var.image

      ports {
        container_port = var.container_port
      }
    }
  }
}

data "google_dns_managed_zone" "zone" {
  name = var.dns_zone_name
}

resource "google_cloud_run_domain_mapping" "httpbin" {
  location = var.region
  name     = var.custom_domain

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.httpbin.name
  }
}

resource "google_dns_record_set" "httpbin" {
  managed_zone = data.google_dns_managed_zone.zone.name
  name         = "${var.custom_domain}."
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]

  depends_on = [google_cloud_run_domain_mapping.httpbin]
}
