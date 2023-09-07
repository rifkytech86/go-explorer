package cms

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/api/controller/cms_controllers"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/repository/cms_repository"
	"gitlab.com/naonweh-studio/bubbme-backend/usecase/cms_usecase"
	"time"
)

func NewUserCoin(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewUserCoinRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.UserCoinController{
		UserCoinUseCase: cms_usecase.NewUserCoinUseCase(ur, env, timeout),
	}
	group.GET("/user-coin", pc.Get)
	group.POST("/user-coin", pc.Create)
	group.PATCH("/user-coin/:id", pc.Update)
	group.DELETE("/user-coin/:id", pc.Delete)
}
