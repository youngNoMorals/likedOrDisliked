package commons

import (
	"database/sql"
	"fmt"
)

// const
//  @Description: 链接数据库参数
//  @return unc
const (
	host = "localhost"
	port = 5432
	sqlUser = "postgres"
	password = "root"
	dbname ="postgres"
)

// OpenConnection
//  @Description: 链接数据库
//  @return *sql.DB
func OpenConnection() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, sqlUser, password, dbname)

	db,err :=sql.Open("postgres",psqlInfo)
	if err !=nil{
		panic(err)
	}

	err = db.Ping()
	if err !=nil{
		panic(err)
	}

	return db
}