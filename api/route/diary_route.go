package route

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/api/controller"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/repository"
	"gitlab.com/naonweh-studio/bubbme-backend/usecase"
	"time"
)

func NewDiaryRouter(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	tr := repository.NewDiaryRepository(db, domain.CollectionTask)
	tc := &controller.DiaryController{
		DiaryUsecase: usecase.NewDiaryUsecase(tr, timeout),
	}
	group.GET("/diary", tc.FetchDiary)
	group.POST("/diary", tc.CreateDiary)
	group.PATCH("/diary/:id", tc.UpdateDiary)
	group.DELETE("/diary/:id", tc.DeleteDiary)
}
