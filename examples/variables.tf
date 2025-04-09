variable "client_id" {
  description = "ID of an OAuth App associated with the organization. It is recommended to use an OAuth App with least-privileged roles in automated environments."
  default     = ""
}

variable "client_secret" {
  description = "Secret of the OAuth App, associated with the organization. It is recommended to use an OAuth App with least-privileged roles in automated environments."
  default     = ""
}

variable "org_id" {
  description = "Organization Identifier."
  default     = ""
}

variable "aws_account_number" {
  description = "The AWS account number."
  default     = ""
}

variable "sddc_name" {
  description = "Name of SDDC."
  default     = "sddc-test"
}

variable "sddc_region" {
  description = "The AWS or VMC specific region."
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

variable "public_ip_displayname" {
  description = "Display name for public IP."
  default     = "public-ip-test"
}

variable "host_instance_type" {
  description = "The instance type for the ESX hosts in the primary cluster of the SDDC. Possible values: I3_METAL, I4I_METAL."
  default     = ""
}

variable "sddc_primary_cluster_num_hosts" {
  description = "The number of hosts in the primary cluster of the SDDC."
  default     = 1
}

variable "provider_type" {
  description = "Determines what additional properties are available based on cloud provider. Default value : AWS"
  default     = "AWS"
}

variable "sddc_type" {
  description = "Denotes the SDDC type, if the value is null or empty, the type is considered as default. Possible values : '1NODE', 'DEFAULT'. "
  default     = "1NODE"
}

variable "cluster_num_hosts" {
  description = "The number of hosts in the cluster."
  default     = 3
}

variable "site_recovery_srm_extension_key_suffix" {
  description = "The custom extension suffix for SRM must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The suffix is appended to com.vmware.vcDr- to form the full extension key"
  default     = ""
}

variable "srm_node_srm_extension_key_suffix" {
  description = "The custom extension suffix for SRM must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The suffix is appended to com.vmware.vcDr- to form the full extension key"
  default     = "test"
}
