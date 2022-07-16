variable "org_id" {
  type = string
}

variable "billing_account" {
  type = string
}

variable "project_id" {
  type    = string
  default = "microbot-hackathon"
}

variable "repository_location" {
  type    = string
  default = "us" // default to the "us" multi-region location
}

variable "repository_id" {
  type    = string
  default = "bots-repo"
}
