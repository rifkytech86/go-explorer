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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, cache *bootstrap.RedisClient, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, cache, domain.CollectionUser)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, cache, timeout),
		Env:           env,
		Cache:         cache,
	}
	group.POST("/signup", sc.Signup)
	group.POST("/verify-otp", sc.VerifyOTP)
}
