// Package version contains global variables for
// nolint: gochecknoglobals, gochecknoinits
package config

const (
	FileName = ".mamba"
	ConfigType = ".json"
	// ServiceName defines short service name
	ServiceName = "Mamba"
	// DefaultPostgresPort defines default port for PostgreSQL
	DefaultPostgresPort = 5432
	// DefaultMySQLPort defines default port for MySQL
	DefaultMySQLPort = 3306
	// Base declared base templates
	Base = "base"
	// GKE declared GKE accounts/cluster/deployment
	GKE = "gke"
	// API declared type API
	API = "api"
	// APIGateway declared type API gateway: REST
	APIGateway = "rest"
	// APIgRPC declared type API: gRPC
	APIgRPC = "grpc"
	// Contract declared contract API example
	Contract = "contract"
	// Storage declared type Storage
	Storage = "storage"
	// StoragePostgres declared storage driver type: postgres
	StoragePostgres = "postgres"
	// StorageMySQL declared storage driver type: mysql
	StorageMySQL = "mysql"
	// RELEASE returns the release version
	RELEASE = "v0.1.1"
	// DATE returns the release date
	DATE = "UNKNOWN"
)
