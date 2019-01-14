provider "vmc" {
  refresh_token = ""
}

data "vmc_org" "my_org" {
  id = "058f47c4-92aa-417f-8747-87f3ed61cb45"
}

resource "vmc_sddc" "sddc_1" {
  org_id        = "${data.vmc_org.my_org.id}"
  sddc_name     = "terraform sddc 2"
  num_host      = 1
  provider_type = "ZEROCLOUD"
  region        = "US_WEST_2"
}
