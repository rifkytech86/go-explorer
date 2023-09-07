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

func NewAuthorization(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewAuthorization(db, nil, domain.CollectionUser)
	pc := &cms_controllers.AuthorizationController{
		AuthorizationUseCase: cms_usecase.NewAuthorization(ur, env, timeout),
	}
	group.POST("/login", pc.Login)
	group.POST("/logout", pc.Logout)
}
