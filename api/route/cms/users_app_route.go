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

func NewUsersApp(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewUsersAppRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.UsersAppController{
		UsersAppUseCase: cms_usecase.NewUsersAppUseCase(ur, env, timeout),
	}
	group.GET("/users-app", pc.GetUserApp)
	group.POST("/users-app", pc.CreateUserApp)
	group.PATCH("/users-app/:id", pc.UpdateUserApp)
	group.DELETE("/users-app/:id", pc.DeleteUserApp)
}
