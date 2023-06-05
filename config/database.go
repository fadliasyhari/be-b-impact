package config

// import (
// 	"fmt"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type DbConnection interface {
// 	Conn() *gorm.DB
// 	Migrate(model ...any) error
// }

// type dbConnection struct {
// 	db  *gorm.DB
// 	cfg *Config
// }

// func (d *dbConnection) initDb() error {
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", d.cfg.Host, d.cfg.Port, d.cfg.User, d.cfg.Password, d.cfg.Name)
// 	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		return err
// 	}

// 	d.db = conn
// 	return nil
// }

// func (d *dbConnection) Conn() *gorm.DB {
// 	return d.db
// }

// func (d *dbConnection) Migrate(model ...any) error {
// 	err := d.Conn().AutoMigrate(model...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func NewDbConnection(cfg *Config) (DbConnection, error) {
// 	conn := &dbConnection{
// 		cfg: cfg,
// 	}
// 	err := conn.initDb()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return conn, nil
// }
