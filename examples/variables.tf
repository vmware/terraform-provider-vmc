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
  description = "The instance type for the ESX hosts in the primary cluster of the SDDC. Possible values: I3_METAL, R5_METAL."
  default     = ""
}

variable storage_capacity {
  description = "The storage capacity value to be requested for the sddc primary cluster, in GiBs. If provided, instead of using the direct-attached storage, a capacity value amount of seperable storage will be used. Possible values for R5 metal are 15TB, 20TB, 25TB, 30TB, 35TB."
  default     = ""
}

variable num_hosts {
  description = "The number of hosts."
  default     = 1
}

variable provider_type {
  description = "Determines what additional properties are available based on cloud provider. Acceptable values are ZEROCLOUD and AWS with AWS as the default value."
  default     = "AWS"
}

variable sddc_type {
description = "Denotes the sddc type, if the value is null or empty, the type is considered as default. Possible values : '1NODE', 'DEFAULT'. "
default = "1NODE"
}

variable srm_extension_key_suffix {
  description = "Customization, for example can specify custom extension key suffix for SRM."
  default     = ""
}
