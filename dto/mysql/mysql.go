package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/url"
	"time"
)

type MysqlClient struct {
	Conn        *sql.DB
	ConnCluster *sql.DB
}

func NewClient(connection string, connectionCluster string) (MysqlClient, error) {
	time.Local = time.UTC
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	val.Add("multiStatements", "true")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dsnCluster := fmt.Sprintf("%s?%s", connectionCluster, val.Encode())
	dbConnCluster, err := sql.Open(`mysql`, dsnCluster)
	err = dbConnCluster.Ping()
	if err != nil {
		fmt.Println("failed connection to dsnCluster")
		log.Fatal(err)
	}

	//
	//
	//path := fmt.Sprintf("mysql://%s", dsn)
	//m, err := migrate.New("file://./db/migrations", path)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := m.Up(); err != nil {
	//	if err == migrate.ErrNoChange {
	//		log.Println("No migrations to apply")
	//	} else {
	//		log.Fatal(err)
	//	}
	//}
	//
	//// cluster
	//dsnCluster := fmt.Sprintf("%s?%s", connectionCluster, val.Encode())
	//dbConnCluster, err := sql.Open(`mysql`, dsnCluster)
	//err = dbConnCluster.Ping()
	//if err != nil {
	//	fmt.Println("failed connection to dsnCluster")
	//	log.Fatal(err)
	//}
	//
	//pathCluster := fmt.Sprintf("mysql://%s", dsnCluster)
	//mCluster, errCLuster := migrate.New("file://./db/migrations", pathCluster)
	//if errCLuster != nil {
	//	fmt.Println("migrationError")
	//	log.Fatal(err)
	//}
	//
	//if err := mCluster.Up(); err != nil {
	//	if err == migrate.ErrNoChange {
	//		log.Println("No migrations to apply cluster")
	//	} else {
	//		log.Fatal(err)
	//	}
	//}

	return MysqlClient{Conn: dbConn, ConnCluster: dbConnCluster}, err
}
