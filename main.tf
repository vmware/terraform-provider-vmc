provider "vmc" {
  refresh_token = "ac5140ea-1749-4355-a892-56cff4893be0"
  csp_url       = "https://console-stg.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  id = "05e0a625-3293-41bb-a01f-35e762781c2a"
}

resource "vmc_sddc" "sddc_1" {
  org_id        = "${data.vmc_org.my_org.id}"
  sddc_name     = "terraform sddc 1"
  num_host      = 4
  provider_type = "ZEROCLOUD"
  region        = "US_WEST_1"
}
