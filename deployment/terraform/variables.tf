variable "gcp_project_id" {
  type        = string
  description = "GCP project id"
}

variable "gcp_region" {
  type        = string
  description = "GCP project region"
}

variable "gke_cluster_name" {
  type        = string
  description = "GKE cluster name"
}

variable "gke_location" {
  type        = string
  description = "GKE location"
}

variable "gke_pool_name" {
  type        = string
  description = "GKE node pool name"
}

variable "gke_node_preemptible" {
  type        = bool
  description = "GKE node preemptible"
}

variable "gke_node_machine_type" {
  type        = string
  description = "GKE node machine type"
}

variable "gke_node_disk_size_gb" {
  type        = number
  description = "GKE node disk size in gb"
}

variable "gcr_image_name" {
  type        = string
  description = "GCR image name"
}

variable "gke_deployment_name" {
  type        = string
  description = "GKE deployment bot name"
}

variable "naka_kai_discord_token" {
  type        = string
  description = "Discord token"
}

variable "naka_kai_discord_delete_time" {
  type        = string
  description = "Discord message delete time"
}

variable "naka_kai_discord_queue_limit" {
  type        = string
  description = "Discord queue limit"
}

variable "naka_kai_db_dialect" {
  type        = string
  description = "Database dialect"
}

variable "naka_kai_db_address" {
  type        = string
  description = "Database address"
}

variable "naka_kai_db_name" {
  type        = string
  description = "Database name"
}

variable "naka_kai_db_user" {
  type        = string
  description = "Database username"
}

variable "naka_kai_db_password" {
  type        = string
  description = "Database password"
}

variable "naka_kai_youtube_key" {
  type        = string
  description = "Youtube API key"
}

variable "naka_kai_log_level" {
  type        = number
  description = "Log level"
}

variable "naka_kai_log_json" {
  type        = bool
  description = "Log json"
}

variable "naka_kai_newrelic_license_key" {
  type        = string
  description = "Newrelic license key"
}
