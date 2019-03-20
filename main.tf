provider "vmc" {
  refresh_token = ""

  # for staging environment only
  # vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
  # csp_url       = "https://console-stg.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  id = ""
}

data "vmc_connected_accounts" "accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}

resource "vmc_sddc" "sddc_1" {
  org_id = "${data.vmc_org.my_org.id}"

  # storage_capacity    = 100
  sddc_name           = ""
  vpc_cidr            = "10.2.0.0/16"
  num_host            = 1
  provider_type       = "AWS"
  region              = "US_EAST_1"
  vxlan_subnet        = "192.168.1.0/24"
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  sddc_template_id = ""
  deployment_type  = "SingleAZ"

  account_link_sddc_config = [
    {
      customer_subnet_ids  = [""]
      connected_account_id = "${data.vmc_connected_accounts.accounts.ids.0}"
    },
  ]
}
