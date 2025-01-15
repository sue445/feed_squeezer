variable "service_name" {
  type        = string
  description = "Cloud Run service name"
  default     = "feed-squeezer"
}

variable "location" {
  type        = string
  description = "Cloud Run app location"
  default     = "us-central1"
}

# c.f. https://console.cloud.google.com/artifacts/docker/feed-squeezer/us/feed-squeezer/app
variable "tag" {
  type        = string
  description = "docker image tag for feed_squeezer"
  default     = "latest"
}

variable "sentry_dsn" {
  type        = string
  description = "Sentry DSN"
  default     = ""
}

variable "min_instance_count" {
  type        = number
  description = "min instances"
  default     = 0
}

variable "max_instance_count" {
  type        = number
  description = "max instances"
  default     = 1
}

variable "service_account_name" {
  type        = string
  description = "ServiceAccount name"
  default     = "feed-squeezer"
}
