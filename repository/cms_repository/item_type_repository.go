package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type itemTypeRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewItemTypeRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IItemTypeRepository {
	return &itemTypeRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *itemTypeRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.ItemType, err error) {
	query := `SELECT item_type_id, item_type_name, created_at, updated_at, created_by, updated_by FROM game_item_type`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "item_type_name LIKE ?")
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
	var listData []*cms_domain.ItemType

	for results.Next() {
		var data = &cms_domain.ItemType{}
		err = results.Scan(
			&data.ItemTypeID,
			&data.ItemTypeName,
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
func (ur *itemTypeRepo) GetTotalItemType(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(item_type_id) as total from game_item_type`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "item_type_name LIKE ?")
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

func (ur *itemTypeRepo) CreateItemType(ctx context.Context, payload *cms_domain.ItemType) (int64, error) {
	query := `INSERT INTO game_item_type(item_type_name, created_at, created_by)
				VALUES(?,NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.ItemTypeName,
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

func (ur *itemTypeRepo) UpdateItemType(ctx context.Context, payload *cms_domain.ItemType, id int) error {

	query := "UPDATE game_item_type set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.ItemTypeName.String != "" {
		columns = append(columns, " item_type_name = ?")
		params = append(params, payload.ItemTypeName)

	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE item_type_id = ?"
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

func (ur *itemTypeRepo) DeleteItemType(ctx context.Context, id int) error {
	query := "DELETE FROM game_item_type WHERE item_type_id = ?"
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
