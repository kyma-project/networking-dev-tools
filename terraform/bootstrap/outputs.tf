output "bucket_name" {
  description = "Name of the GCS bucket that stores OpenTofu state for the main module."
  value       = google_storage_bucket.tfstate.name
}
