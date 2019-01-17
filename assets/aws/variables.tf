variable "region" {
  type        = "string"
  description = "Region to create resources in. See https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html for valid values."
}

variable "ssh_public_key" {
  type        = "string"
  description = "A public key line in .ssh/authorized_keys format to use to authenticate to your instance. This must be added to your SSH agent for provisioning to succeed."
}

variable "paramstore_var" {
  default     = "/guestbook/motd"
  description = "The location in SSM Parameter Store of the Message of the Day variable."
}
