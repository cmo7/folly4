package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Engine string

const (
	MySQL    Engine = "mysql"
	SQLite   Engine = "sqlite"
	Postgres Engine = "postgres"
)

type ConnectionData struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Engine   Engine
	File     string
}

func init() {
	// Default configuration values
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "root")
	viper.SetDefault("database.name", "folly")
	viper.SetDefault("database.engine", "sqlite")
	viper.SetDefault("database.file", "folly.db")
}

func ConnectWithConfig(config *gorm.Config, connectionData *ConnectionData) (*gorm.DB, error) {
	var connection *gorm.DB
	var err error
	switch connectionData.Engine {
	case MySQL:
		connection, err = connectMySQL(connectionData, config)
	case SQLite:
		connection, err = connectSQLite(connectionData, config)
	case Postgres:
		connection, err = connectPostgres(connectionData, config)
	default:
		connection, err = nil, fmt.Errorf("unsupported engine: %s", connectionData.Engine)
	}
	return connection, err
}

func Connect(config *gorm.Config) (*gorm.DB, error) {

	connectionData := &ConnectionData{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.name"),
		Engine:   Engine(viper.GetString("database.engine")),
		File:     viper.GetString("database.file"),
	}

	var connection *gorm.DB
	var err error
	switch connectionData.Engine {
	case MySQL:
		connection, err = connectMySQL(connectionData, config)
	case SQLite:
		connection, err = connectSQLite(connectionData, config)
	case Postgres:
		connection, err = connectPostgres(connectionData, config)
	default:
		connection, err = nil, fmt.Errorf("unsupported engine: %s", connectionData.Engine)
	}
	return connection, err
}

func connectMySQL(connectionData *ConnectionData, config *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		connectionData.Username,
		connectionData.Password,
		connectionData.Host,
		connectionData.Port,
		connectionData.Database,
	)
	return gorm.Open(mysql.Open(dsn), config)
}

func connectSQLite(connectionData *ConnectionData, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(connectionData.File), config)
}

func connectPostgres(connectionData *ConnectionData, config *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		connectionData.Host,
		connectionData.Username,
		connectionData.Password,
		connectionData.Database,
		connectionData.Port,
	)
	return gorm.Open(postgres.Open(dsn), config)
}
