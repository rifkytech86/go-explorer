package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type userPointRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewUserPointRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IUserPointRepository {
	return &userPointRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *userPointRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.UserPoint, err error) {
	query := `SELECT user_point_id, user_id, user_point_value, created_at, updated_at, created_by, updated_by FROM user_point`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " user_coin_value LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, " user_point_value LIKE ?")
		params = append(params, "%"+search+"%")

		query += " WHERE " + "(" + strings.Join(conditions, " OR ") + ")"
	}

	query += " LIMIT ?  OFFSET ?"
	offset := (page - 1) * limit
	params = append(params, limit)
	params = append(params, offset)
	results, err := ur.database.Conn.Query(query, params...)

	if err != nil {
		return nil, err
	}
	var listData []*cms_domain.UserPoint

	for results.Next() {
		var data = &cms_domain.UserPoint{}
		err = results.Scan(
			&data.UserPointID,
			&data.UserID,
			&data.UserPointValue,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.CreatedBy,
			&data.UpdatedBy,
		)
		listData = append(listData, data)
		if err != nil {
			return nil, err
		}
	}
	return listData, nil
}
func (ur *userPointRepo) GetTotal(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(user_point_id) as total from user_point`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " user_id LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, " user_point_value LIKE ?")
		params = append(params, "%"+search+"%")

		query += " WHERE " + "(" + strings.Join(conditions, " OR ") + ")"
	}
	var count int
	stmt := ur.database.Conn.QueryRow(query, params...)
	if err != nil {
		return 0, err
	}
	err = stmt.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ur *userPointRepo) Create(ctx context.Context, payload *cms_domain.UserPoint) (int64, error) {
	query := `INSERT INTO user_point(user_id, user_point_value, created_at, created_by)
				VALUES(?,?, NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.UserID,
		payload.UserPointValue,
		payload.CreatedBy)
	if err != nil {
		logger.Error(err.Error(), nil)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		logger.Error(err.Error(), nil)
		return 0, err
	}

	return lastId, nil
}

func (ur *userPointRepo) Update(ctx context.Context, payload *cms_domain.UserPoint, id int) error {

	query := "UPDATE user_point set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.UserPointValue.Int64 != 0 {
		columns = append(columns, " user_point_value = ?")
		params = append(params, payload.UserPointValue)
	}
	if payload.UserID.Int64 != 0 {
		columns = append(columns, " user_id = ?")
		params = append(params, payload.UserID)
	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE user_point_id = ?"
	params = append(params, id)
	stmt, err := tx.Prepare(query)
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(params...)
	if err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}
	return nil
}

func (ur *userPointRepo) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM user_point WHERE user_point_id = ?"
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error(err.Error(), nil)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}
	return nil
}
