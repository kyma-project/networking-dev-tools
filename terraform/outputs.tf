output "service_url" {
  description = "Auto-generated Cloud Run *.run.app URL for the service."
  value       = google_cloud_run_v2_service.httpbin.uri
}

output "custom_domain_url" {
  description = "HTTPS URL for the configured custom domain."
  value       = "https://${var.custom_domain}"
}

output "dns_records" {
  description = "DNS records returned by the domain mapping (for debugging first-apply behavior)."
  value       = google_cloud_run_domain_mapping.httpbin.status[0].resource_records
}
