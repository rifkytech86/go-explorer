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

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(db, domain.CollectionTask)
	tc := &controller.TaskController{
		TaskUsecase: usecase.NewTaskUsecase(tr, timeout),
	}
	group.GET("/task", tc.Fetch)
	group.POST("/task", tc.Create)
}
