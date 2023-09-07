package repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"

	"gitlab.com/naonweh-studio/bubbme-backend/domain"
)

type taskRepository struct {
	database   mysql.MysqlClient
	collection string
}

func NewTaskRepository(db mysql.MysqlClient, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	//collection := tr.database.Collection(tr.collection)
	//
	//_, err := collection.InsertOne(c, task)

	return nil
}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	//collection := tr.database.Collection(tr.collection)
	//
	//var tasks []domain.Task
	//
	//idHex, err := primitive.ObjectIDFromHex(userID)
	//if err != nil {
	//	return tasks, err
	//}
	//
	//cursor, err := collection.Find(c, bson.M{"userID": idHex})
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = cursor.All(c, &tasks)
	//if tasks == nil {
	//	return []domain.Task{}, err
	//}

	//return tasks, err
	return nil, nil
}
