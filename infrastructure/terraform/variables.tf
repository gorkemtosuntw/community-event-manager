variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "development"
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
  default     = "event-platform-cluster"
}
