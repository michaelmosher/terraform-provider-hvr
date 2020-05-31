package hvr

import "fmt"

type Config struct {
	PostgresqlDatabaseHost     string
	PostgresqlDatabasePort     int
	PostgresqlDatabaseName     string
	PostgresqlDatabaseUsername string
	PostgresqlDatabasePassword string
}

type Service interface {
}

type postgresqlService struct {
}

func (c Config) Client() (Service, error) {
	if c.PostgresqlDatabaseHost != "" {
		return postgresqlService{}, nil
	}

	return nil, fmt.Errorf("one of (`postgresql_database_host`) must be defined")
}
