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

func NewMood(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewMoodRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.MoodController{
		MoodUseCase: cms_usecase.NewMoodUseCase(ur, env, timeout),
	}
	group.GET("/mood", pc.GetUserApp)
	group.POST("/mood", pc.CreateUserApp)
	group.PATCH("/mood/:id", pc.UpdateUserApp)
	group.DELETE("/mood/:id", pc.DeleteUserApp)
}
