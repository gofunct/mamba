variable "project" {
  type        = "string"
  description = "Project to set up."
}

variable "region" {
  type        = "string"
  description = "GCP region to create database and storage in, for example 'us-central1'. See https://cloud.google.com/compute/docs/regions-zones/ for valid values."
}

variable "zone" {
  type        = "string"
  description = "GCP zone to create the GKE cluster in, for example 'us-central1-a'. See https://cloud.google.com/compute/docs/regions-zones/ for valid values."
}

variable "server_service_account_name" {
  default     = "guestbook"
  description = "The username part of the service account email that will be used for the server running inside the GKE cluster."
}

variable "db_access_service_account_name" {
  default     = "guestbook-db"
  description = "The username part of the service account email that will be used for provisioning the database."
}

variable "cluster_name" {
  default     = "guestbook-cluster"
  description = "The GKE cluster name."
}
