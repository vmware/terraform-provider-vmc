provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}
# Empty data source defined in order to store the org display name and name in terraform state
data "vmc_org" "my_org" {
}

data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = replace(upper(var.sddc_region), "-", "_")
}

resource "vmc_sddc" "sddc_1" {
  sddc_name           = var.sddc_name
  vpc_cidr            = var.vpc_cidr
  num_host            = var.sddc_num_hosts
  provider_type       = var.provider_type
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = var.deployment_type
  host_instance_type = var.host_instance_type

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0],data.vmc_customer_subnets.my_subnets.ids[1]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}