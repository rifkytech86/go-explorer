package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gitlab.com/naonweh-studio/bubbme-backend/api/middleware"
	"gitlab.com/naonweh-studio/bubbme-backend/api/route/cms"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	_ "gitlab.com/naonweh-studio/bubbme-backend/cmd/docs"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"time"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, cache *bootstrap.RedisClient, gin *gin.Engine) {

	publicRouter := gin.Group("api/v1")
	// All Public APIs
	NewSignupRouter(env, timeout, db, cache, publicRouter)
	NewLoginRouter(env, timeout, db, cache, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("api/v1")

	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	NewProfileRouter(env, timeout, db, protectedRouter)
	NewTaskRouter(env, timeout, db, protectedRouter)
	NewDiaryRouter(env, timeout, db, protectedRouter)

}

func SetupCMS(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, cache *bootstrap.RedisClient, gin *gin.Engine) {
	cmsRoute := gin.Group("api/v1/cms")
	cms.NewAuthorization(env, timeout, db, cmsRoute)

	protectedRouter := gin.Group("api/v1/cms")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	cms.NewUsers(env, timeout, db, protectedRouter)
	cms.NewUsersApp(env, timeout, db, protectedRouter)
	cms.NewMood(env, timeout, db, protectedRouter)
	cms.NewDiary(env, timeout, db, protectedRouter)
	cms.NewPetStatus(env, timeout, db, protectedRouter)
	cms.NewItemType(env, timeout, db, protectedRouter)
	cms.NewItem(env, timeout, db, protectedRouter)
	cms.NewCoinSource(env, timeout, db, protectedRouter)
	cms.NewPointSource(env, timeout, db, protectedRouter)
	cms.NewUserCoin(env, timeout, db, protectedRouter)
	cms.NewUserPoint(env, timeout, db, protectedRouter)

}

func SetupSwagger(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, cache *bootstrap.RedisClient, gin *gin.Engine) {
	protectedRouter := gin.Group("api/v1/swagger")
	protectedRouter.GET("*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
