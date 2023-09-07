package route

import (
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/api/controller"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/repository"
	"gitlab.com/naonweh-studio/bubbme-backend/usecase"
)

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, cache *bootstrap.RedisClient, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, nil, domain.CollectionUser)
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, timeout),
		Env:          env,
		Cache:        cache,
	}
	group.POST("/login", lc.Login)
}
