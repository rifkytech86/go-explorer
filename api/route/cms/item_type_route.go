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

func NewItemType(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewItemTypeRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.ItemTypeController{
		ItemTypeUseCase: cms_usecase.NewItemTypeUseCase(ur, env, timeout),
	}
	group.GET("/item-type", pc.GetItemType)
	group.POST("/item-type", pc.CreateItemType)
	group.PATCH("/item-type/:id", pc.UpdateItemType)
	group.DELETE("/item-type/:id", pc.DeleteItemType)
}
