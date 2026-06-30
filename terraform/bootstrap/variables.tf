variable "project_id" {
  description = "GCP project that owns the state bucket."
  type        = string
}

variable "state_bucket_name" {
  description = "Globally unique name for the GCS bucket that will hold OpenTofu state."
  type        = string
  default     = "kyma-networking-dev-tools-tfstate"
}

variable "location" {
  description = "Bucket location."
  type        = string
  default     = "EU"
}
