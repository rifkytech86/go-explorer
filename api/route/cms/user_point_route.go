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

func NewUserPoint(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewUserPointRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.UserPointController{
		UserPointUseCase: cms_usecase.NewUserPointCase(ur, env, timeout),
	}
	group.GET("/user-point", pc.Get)
	group.POST("/user-point", pc.Create)
	group.PATCH("/user-point/:id", pc.Update)
	group.DELETE("/user-point/:id", pc.Delete)
}
