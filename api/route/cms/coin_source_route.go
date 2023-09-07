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

func NewCoinSource(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewCoinSourceRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.CoinSourceController{
		CoinSourceUseCase: cms_usecase.NewCoinSourceUseCase(ur, env, timeout),
	}
	group.GET("/coin-source", pc.GetCoinSource)
	group.POST("/coin-source", pc.CreateCoinSource)
	group.PATCH("/coin-source/:id", pc.UpdateCoinSource)
	group.DELETE("/coin-source/:id", pc.DeleteCoinSource)
}
