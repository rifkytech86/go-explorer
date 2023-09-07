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

func NewUsers(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewUsers(db, nil, domain.CollectionUser)
	pc := &cms_controllers.UsersController{
		UsersUseCase: cms_usecase.NewUsersAdminUseCase(ur, env, timeout),
	}
	group.GET("/users", pc.GetUserAdmin)
	group.POST("/users", pc.CreateUserAdmin)
	group.PATCH("/users/:id", pc.UpdateUserAdmin)
	group.DELETE("/users/:id", pc.DeleteUserAdmin)
}
