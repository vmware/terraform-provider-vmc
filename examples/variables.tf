variable "sddc_region" {
  description = "The AWS region to create things in."
  default     = "US_WEST_2"
}

variable "vpc_cidr" {
  description = "The cidr for VPC."
  default     = "10.2.0.0/16"
}

variable "vxlan_subnet" {
  description = "The subnet of SDDC."
  default     = "192.168.1.0/24"
}

variable "private_ip" {
  description = "The private IP of SDDC."
  default     = "10.2.33.45"
}
