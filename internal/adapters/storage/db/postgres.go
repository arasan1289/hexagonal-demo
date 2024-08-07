package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
)

// Conn implements the DB interface using GORM.
type Conn struct {
	*gorm.DB
	url    string
	config *config.Container
}

// New creates a new GORM connection.
func New(config *config.Container, logger logger.Interface) (*Conn, error) {
	// Connect to the PostgreSQL database
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connectionString,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 logger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	db = db.InstanceSet("config", config)
	return &Conn{db, connectionString, config}, nil
}

// Set sets the connection pool.
func (c *Conn) Set() error {
	db, err := c.DB.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return nil
}

// Migrate migrates the models to the DB
// TODO: Switch to any migrator interface
func (c *Conn) Migrate(models ...interface{}) error {
	// Auto migrate the models
	return c.DB.AutoMigrate(models...)
}

// Close closes the connection
func (c *Conn) Close() error {
	db, err := c.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
