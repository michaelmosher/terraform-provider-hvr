provider "hvr" {
  connection_type = "postgresql"
  hub_database_host = "localhost"
  hub_database_name = "hvrhub"
  hub_database_username = "michaelmosher"
}

resource "hvr_location" "this" {
  name = "hello"
  class = "db2i"

  database_connection {
    node = "192.168.0.247"
    name = "AMTEST"
    username = "user"
    password = "pass"
  }

  remote_machine {
    node = "localhost"
    port = 15131
    username = "cdc"
    password = "N/A"
  }
}

resource "hvr_channel" "this" {
  name = "hello->world"
  description = "Replicates data from `hello` source to `world` target."
}
