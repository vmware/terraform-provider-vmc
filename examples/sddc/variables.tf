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
  description = "SDDC management network CIDR. Only prefix of 16, 20 and 23 are supported."
  default     = ""
}

variable "vxlan_subnet" {
  description = "A logical network segment that will be created with the SDDC under the compute gateway."
  default     = ""
}

variable host_instance_type {
  description = "The instance type for the ESX hosts in the primary cluster of the SDDC. Possible values: I3_METAL, R5_METAL."
  default     = ""
}

variable storage_capacity {
  description = "The storage capacity value to be requested for the SDDC primary cluster. This variable is only for R5.METAL. Possible values are 15TB, 20TB, 25TB, 30TB, 35TB per host."
  default     = ""
}

variable sddc_num_hosts {
  description = "The number of hosts in SDDC."
  default     = 1
}

variable provider_type {
  description = "Determines what additional properties are available based on cloud provider. Default value : AWS"
  default     = "AWS"
}

variable sddc_type {
  description = "Denotes the sddc type, if the value is null or empty, the type is considered as default. Possible values : '1NODE', 'DEFAULT'. "
  default = "1NODE"
}