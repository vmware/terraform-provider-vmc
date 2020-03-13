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

variable "public_ip_displayname" {
  description = "Display name for public IP."
  default     = "public-ip-test"
}


variable host_instance_type {
  description = "The instance type for the ESX hosts in the primary cluster of the SDDC"
  default     = "I3_METAL"
}

variable storage_capacity {
  description = "Storage capacity"
  default     = ""
}

variable num_hosts {
  description = "The number of hosts."
  default     = "1"
}

variable provider_type {
  description = "Determines what additional properties are available based on cloud provider. Acceptable values are ZEROCLOUD and AWS with AWS as the default value."
  default     = ""
}
