package hvr

import (
	"fmt"

	"github.com/michaelmosher/go-hvr-sdk/hvrhub"
	hvrhub_postgresql "github.com/michaelmosher/go-hvr-sdk/hvrhub/postgresql"
)

type Config struct {
	PostgresqlDatabaseHost     string
	PostgresqlDatabasePort     int
	PostgresqlDatabaseName     string
	PostgresqlDatabaseUsername string
	PostgresqlDatabasePassword string
}

type Service interface {
	GetLocation(string) (hvrhub.Location, error)
	NewLocation(hvrhub.Location) error
	UpdateLocation(hvrhub.Location) error
	DeleteLocation(string) error

	GetChannel(string) (hvrhub.Channel, error)
	NewChannel(hvrhub.Channel) error
	UpdateChannel(hvrhub.Channel) error
	DeleteChannel(string) error

	GetLocationGroup(string, string) (hvrhub.LocationGroup, error)
	NewLocationGroup(hvrhub.LocationGroup) error
	UpdateLocationGroup(hvrhub.LocationGroup) error
	DeleteLocationGroup(string, string) error

	GetLocationGroupMember(string, string, string) (hvrhub.LocationGroupMember, error)
	NewLocationGroupMember(hvrhub.LocationGroupMember) error
	DeleteLocationGroupMember(string, string, string) error
}

func (c Config) Client() (Service, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable user=%s password=%s",
		c.PostgresqlDatabaseHost,
		c.PostgresqlDatabasePort,
		c.PostgresqlDatabaseName,
		c.PostgresqlDatabaseUsername,
		c.PostgresqlDatabasePassword,
	)

	pgClient, err := hvrhub_postgresql.New(connStr)
	if err != nil {
		return nil, fmt.Errorf("error creating hub connection: %s", err)
	}

	hub := hvrhub.Service{Client: pgClient}
	return hub, nil
}
