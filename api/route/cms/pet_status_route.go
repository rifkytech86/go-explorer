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

func NewPetStatus(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewPetStatusRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.PetStatusController{
		PetStatusUseCase: cms_usecase.NewPetStatusUseCase(ur, env, timeout),
	}
	group.GET("/pet-status", pc.GetPetStatus)
	group.POST("/pet-status", pc.CreatePetStatus)
	group.PATCH("/pet-status/:id", pc.UpdatePetStatus)
	group.DELETE("/pet-status/:id", pc.DeletePetStatus)
}
