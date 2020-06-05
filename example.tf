provider "hvr" {
  connection_type       = "postgresql"
  hub_database_host     = "localhost"
  hub_database_name     = "hvrhub"
  hub_database_username = "michaelmosher"
}

resource "hvr_location" "target" {
  name        = "epp01"
  class       = "postgresql"
  description = "Production Datalake"

  database_connection {
    node     = "your.db.host"
    name     = "ae_lake"
    port     = 5432
    username = "cdc"
    password = "(sensitive value)"
  }
}

resource "hvr_channel" "this" {
  name        = "my_channel"
  description = "Replicates data from a source to a target."
}

resource "hvr_location_group" "target" {
  channel_name      = hvr_channel.this.name
  group_name        = "TARGET"
  group_description = "targets live here"
}

resource "hvr_location_group_member" "target" {
  channel_name  = hvr_channel.this.name
  group_name    = hvr_location_group.target.group_name
  location_name = hvr_location.target.name
}



provider "aws" {}

variable "source_secret_name" {}

data "aws_secretsmanager_secret_version" "source" {
  secret_id = var.source_secret_name
}

locals {
  source_secret_data = jsondecode(data.aws_secretsmanager_secret_version.source.secret_string)
}

resource "hvr_location" "source" {
  name        = local.source_secret_data["hvr-name"]
  class       = local.source_secret_data["db-class"]
  description = local.source_secret_data["description"]

  database_connection {
    node     = local.source_secret_data["db-host-name"]
    name     = local.source_secret_data["db-name"]
    username = local.source_secret_data["db-user-name"]
    password = local.source_secret_data["db-user-password"]
  }

  remote_machine {
    node     = local.source_secret_data["cdc-host-name"]
    port     = local.source_secret_data["ssh-tunnel-port"]
    username = local.source_secret_data["os-user-name"]
    password = local.source_secret_data["os-user-password"]
  }
}

resource "hvr_location_group" "source" {
  channel_name      = hvr_channel.this.name
  group_name        = "SOURCE"
  group_description = "sources live here"
}

resource "hvr_location_group_member" "source" {
  channel_name  = hvr_channel.this.name
  group_name    = hvr_location_group.source.group_name
  location_name = hvr_location.source.name
}
