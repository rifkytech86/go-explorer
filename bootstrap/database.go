package bootstrap

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"log"
	"time"

	"gitlab.com/naonweh-studio/bubbme-backend/mongo"
)

func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

func NewMysqlDatabase(env *Env) mysql.MysqlClient {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbHost := env.DBMysqlHost
	dbPort := env.DBMysqlPort
	dbUser := env.DBMysqlUser
	dbPass := env.DBMysqlPass
	dbName := env.DBMysqlName

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	dbClusterHost := env.DBMysqlClusterHost
	dbClusterPort := env.DBMysqlClusterPort
	dbClusterUser := env.DBMysqlClusterUser
	dbClusterPass := env.DBMysqlClusterPass
	dbClusterName := env.DBMysqlClusterName
	connectionCluster := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbClusterUser, dbClusterPass, dbClusterHost, dbClusterPort, dbClusterName)
	client, err := mysql.NewClient(connection, connectionCluster)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

//
//func CloseMysqlDBConnection(client mysql.Client) {
//	if client == nil {
//		return
//	}
//
//	err := client.Disconnect(context.TODO())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Println("Connection to Mysql closed.")
//}
