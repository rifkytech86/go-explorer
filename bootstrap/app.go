package bootstrap

import (
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/mongo"
)

type Application struct {
	Env        *Env
	Mongo      mongo.Client
	Mysql      mysql.MysqlClient
	RedisCache *RedisClient
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Mysql = NewMysqlDatabase(app.Env)

	app.RedisCache = NewRedisClient(app.Env)
	return *app
}

//func (app *Application) CloseDBConnection() {
//	//CloseMongoDBConnection(app.Mongo)
//	CloseMysqlDBConnection(app.Mysql)
//}
