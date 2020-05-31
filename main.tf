provider "hvr" {
  connection_type = "postgresql"
  hub_database_host = "db-cdc-hvrhub.test.edgeapps.net"
  hub_database_name = "hvrhub"
  hub_database_username = "hvr"
}

resource "hvr_location" "this" {
  name = "hello"
  class = "db2i"

  database_connection {
    node = "AMTEST"
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
