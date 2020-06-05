package hvr

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/michaelmosher/go-hvr-sdk/hvrhub"
)

func expandLocation(d *schema.ResourceData) hvrhub.Location {
	location := hvrhub.Location{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	var connection map[string]interface{}
	for _, c := range d.Get("database_connection").([]interface{}) {
		connection = c.(map[string]interface{})
	}

	switch d.Get("class").(string) {
	case "db2i":
		location.Class = "db2i"
		location.DBName = fmt.Sprintf("%s~%s", connection["node"], connection["name"])
		location.DBUser = fmt.Sprintf("%s/%s", connection["username"], connection["password"])

	case "postgresql":
		location.Class = "postgresql"
		location.DBName = fmt.Sprintf("%s~%d~%s", connection["node"], connection["port"], connection["name"])
		location.DBUser = fmt.Sprintf("%s/%s", connection["username"], connection["password"])
	}

	if rm, isDefined := d.GetOk("remote_machine"); isDefined {
		var remoteMachine map[string]interface{}

		for _, r := range rm.([]interface{}) {
			remoteMachine = r.(map[string]interface{})
		}

		location.RemoteNode = remoteMachine["node"].(string)
		location.RemotePort = remoteMachine["port"].(int)
		location.RemoteLogin = remoteMachine["username"].(string)
		location.RemotePassword = remoteMachine["password"].(string)
	}

	return location
}

func flattenDatabaseConnection(l hvrhub.Location) []interface{} {
	connection := make([]interface{}, 1, 1)

	switch l.Class {
	case "db2i":
		connection[0] = map[string]interface{}{
			"node":     strings.Split(l.DBName, "~")[0],
			"port":     0,
			"name":     strings.Split(l.DBName, "~")[1],
			"username": strings.Split(l.DBUser, "/")[0],
			"password": strings.Split(l.DBUser, "/")[1],
		}

	case "postgresql":
		port, _ := strconv.ParseInt(strings.Split(l.DBName, "~")[1], 0, 32)
		connection[0] = map[string]interface{}{
			"node":     strings.Split(l.DBName, "~")[0],
			"port":     port,
			"name":     strings.Split(l.DBName, "~")[2],
			"username": strings.Split(l.DBUser, "/")[0],
			"password": strings.Split(l.DBUser, "/")[1],
		}
	}

	return connection
}

func flattenRemoteMachine(l hvrhub.Location) []interface{} {
	connection := make([]interface{}, 1, 1)

	connection[0] = map[string]interface{}{
		"node":     l.RemoteNode,
		"port":     l.RemotePort,
		"username": l.RemoteLogin,
		"password": l.RemotePassword,
	}

	return connection
}
