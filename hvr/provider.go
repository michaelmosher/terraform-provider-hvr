package hvr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"connection_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"postgresql"}, false),
				Description:  descriptions["connection_type"],
			},
			"hub_database_host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PGHOST", "localhost"),
				Description: descriptions["database_host"],
			},
			"hub_database_port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5432,
				ValidateFunc: validation.IsPortNumber,
				Description:  descriptions["database_port"],
			},
			"hub_database_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PGDATABASE", "postgres"),
				Description: descriptions["database_name"],
			},
			"hub_database_username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PGUSER", "postgres"),
				Description: descriptions["database_username"],
			},
			"hub_database_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PGPASSWORD", nil),
				Description: descriptions["database_password"],
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hvr_location": resourceHVRLocation(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"connection_type":   "Use to specify the means of connecting to the HVR Hub.",
		"database_host":     "Use to specify the name of the Hub database server.",
		"database_port":     "Use to specify the port number of the Hub database server.",
		"database_name":     "Use to specify the name of Hub database.",
		"database_username": "Use to specify the user name to connect to the Hub database.",
		"database_password": "Use to speicfy the user password for connecting to the Hub database.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	if d.Get("connection_type") == "postgresql" {
		cfg := Config{
			PostgresqlDatabaseHost:     d.Get("hub_database_host").(string),
			PostgresqlDatabasePort:     d.Get("hub_database_port").(int),
			PostgresqlDatabaseName:     d.Get("hub_database_name").(string),
			PostgresqlDatabaseUsername: d.Get("hub_database_username").(string),
			PostgresqlDatabasePassword: d.Get("hub_database_password").(string),
		}

		return cfg.Client()
	}

	return nil, fmt.Errorf("unreachable error?")
}
