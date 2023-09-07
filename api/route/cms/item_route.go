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

func NewItem(env *bootstrap.Env, timeout time.Duration, db mysql.MysqlClient, group *gin.RouterGroup) {
	ur := cms_repository.NewItemRepo(db, nil, domain.CollectionUser)
	pc := &cms_controllers.ItemController{
		ItemUseCase: cms_usecase.NewItemUseCase(ur, env, timeout),
	}
	group.GET("/item", pc.GetItemType)
	group.POST("/item", pc.CreateItemType)
	group.PATCH("/item/:id", pc.UpdateItemType)
	group.DELETE("/item/:id", pc.DeleteItemType)
}
