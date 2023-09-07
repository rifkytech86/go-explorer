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

func NewPointSource(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewPointSourceRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.PointSourceController{
		PointSourceUseCase: cms_usecase.NewPointSourceUseCase(ur, env, timeout),
	}
	group.GET("/point-source", pc.GetPointSource)
	group.POST("/point-source", pc.CreatePointSource)
	group.PATCH("/point-source/:id", pc.UpdatePointSource)
	group.DELETE("/point-source/:id", pc.DeletePointSource)
}
