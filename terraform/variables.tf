variable "project_id" {
  description = "GCP project that hosts the Cloud Run service."
  type        = string
}

variable "region" {
  description = "Region for the Cloud Run service and domain mapping."
  type        = string
  default     = "europe-west1"
}

variable "service_name" {
  description = "Cloud Run service name."
  type        = string
  default     = "httpbin"
}

variable "image" {
  description = "Container image to run."
  type        = string
  default     = "docker.io/mccutchen/go-httpbin:latest"
}

variable "container_port" {
  description = "Container port the image listens on."
  type        = number
  default     = 8080
}

variable "min_instance_count" {
  description = "Minimum Cloud Run instance count."
  type        = number
  default     = null
}

variable "max_instance_count" {
  description = "Maximum Cloud Run instance count."
  type        = number
  default     = 2
}

variable "custom_domain" {
  description = "Fully qualified custom domain for the service."
  type        = string
}

variable "dns_zone_name" {
  description = "Cloud DNS managed zone resource name (not the DNS name)."
  type        = string
}
