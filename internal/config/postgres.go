package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func NewDatabase(viper *viper.Viper) *sqlx.DB {
	var idleConnection int
	var maxConnection int
	var maxLifeTimeConnection int

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD") // viper.GetString("database.password")
	host := os.Getenv("DB_HOST")         //viper.GetString("database.host")
	port := os.Getenv("DB_PORT")         // viper.GetInt("database.port")
	database := os.Getenv("DB_NAME")     //viper.GetString("database.name")
	//idleConnection := viper.GetInt("database.pool.idle")
	//maxConnection := viper.GetInt("database.pool.max")
	//maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	if viper.IsSet("database.pool.idle") {
		idleConnection = viper.GetInt("database.pool.idle")
	} else {
		idleConnection = 10
	}

	if viper.IsSet("database.pool.max") {
		maxConnection = viper.GetInt("database.pool.max")
	} else {
		maxConnection = 100
	}

	if viper.IsSet("database.pool.lifetime") {
		maxLifeTimeConnection = viper.GetInt("database.pool.lifetime")
	} else {
		maxLifeTimeConnection = 100
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", username, password, host, port, database, os.Getenv("DB_PARAMS"))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed connect database: %v", err)
	}

	db.SetMaxIdleConns(idleConnection)
	db.SetMaxOpenConns(maxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	fmt.Println("Database Connected")

	return db
}
