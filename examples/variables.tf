variable "api_token" {
  description = "API token used to authenticate when calling the VMware Cloud Services API."
  default = ""
}

variable "org_id" {
  description = "Organization Identifier."
  default = ""
}

variable "aws_account_number" {
  description = "The AWS account number."
  default     = ""
}

variable "sddc_name"{
  description = "Name of SDDC."
  default = "sddc-test"
}

variable "sddc_region" {
  description = "The AWS  or VMC specific region."
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "AWS VPC IP range. Only prefix of 16 or 20 is currently supported."
  default     = "10.2.0.0/16"
}

variable "vxlan_subnet" {
  description = "VXLAN IP subnet in CIDR for compute gateway."
  default     = "192.168.1.0/24"
}
