package connection

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"github.com/joho/godotenv"
	"os"
)

var DB *sql.DB

var DBConfig *database

type database struct {
	Host          string
	Port          string
	Database      string
	Username      string
	Password      string
	TenantPrefix  string
	DefaultTenant string
	BasicTenant   string
}

func Init() {
	
	info := getDatabaseInfo()
	
	var err error
	DB, err = sql.Open("mysql", info.Username+":"+info.Password+"@tcp("+info.Host+":"+info.Port+")/"+info.Database+"?charset=utf8")
	
	if err != nil {
		log.Fatalf("Error on initializing connection connection: %s", err.Error())
	}
	
	DB.SetMaxIdleConns(100)
	
	err = DB.Ping() // This DOES open a connection if necessary. This makes sure the connection is accessible
	if err != nil {
		log.Fatalf("Error on open")
	}
	
}

func ChangeDB() error {
	if DBConfig == nil {
		DBConfig = getDatabaseInfo()
	}
	
	var err error
	DB, err = sql.Open("mysql", DBConfig.Username+":"+DBConfig.Password+"@tcp("+DBConfig.Host+":"+DBConfig.Port+")/"+DBConfig.Database+"?charset=utf8")
	
	if err != nil {
		log.Fatalf("Error on initializing connection connection: %s", err.Error())
	}
	
	DB.SetMaxIdleConns(100)
	
	err = DB.Ping() // This DOES open a connection if necessary. This makes sure the connection is accessible
	
	return err
}

func getDatabaseInfo() *database {
	
	//从.env文件中读取
	err := godotenv.Load("./.env")
	if err != nil {
		//godotenv.Load("/www/go/scripts/src/main/.env")
		log.Fatal("Error loading .env file, please create .env file in current path")
	}
	
	DBConfig = &database{
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_TENANT_PREFIX"),
		os.Getenv("DEFAULT_TENANT_DATABASE"),
		os.Getenv("BASIC_TENANT_DATABASE"),
	}
	
	return DBConfig
}
