package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type itemRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewItemRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IItemRepository {
	return &itemRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *itemRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.Item, err error) {
	query := `SELECT item_id, item_name, item_desc, item_value, item_type_id, created_at, updated_at, created_by, updated_by FROM game_item`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "item_name LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_desc LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_value LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_type_id LIKE ?")
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
	var listData []*cms_domain.Item

	for results.Next() {
		var data = &cms_domain.Item{}
		err = results.Scan(
			&data.ItemID,
			&data.ItemName,
			&data.ItemDesc,
			&data.ItemValue,
			&data.ItemTypeID,
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
func (ur *itemRepo) GetTotalItem(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(item_id) as total from game_item`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "item_name LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_desc LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_value LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "item_type_id LIKE ?")
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

func (ur *itemRepo) CreateItem(ctx context.Context, payload *cms_domain.Item) (int64, error) {
	query := `INSERT INTO game_item(item_name, item_desc, item_value, item_type_id, created_at, created_by)
				VALUES(?,?,?,?,NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.ItemName,
		payload.ItemDesc,
		payload.ItemValue,
		payload.ItemTypeID,
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

func (ur *itemRepo) UpdateItem(ctx context.Context, payload *cms_domain.Item, id int) error {

	query := "UPDATE game_item set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.ItemName.String != "" {
		columns = append(columns, " item_name = ?")
		params = append(params, payload.ItemName)
	}

	if payload.ItemDesc.String != "" {
		columns = append(columns, " item_desc = ?")
		params = append(params, payload.ItemDesc)
	}

	if payload.ItemValue.Int64 != 0 {
		columns = append(columns, " item_value = ?")
		params = append(params, payload.ItemValue)
	}

	if payload.ItemTypeID.Int64 != 0 {
		columns = append(columns, " item_type_id = ?")
		params = append(params, payload.ItemTypeID)
	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE item_id = ?"
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

func (ur *itemRepo) DeleteItem(ctx context.Context, id int) error {
	query := "DELETE FROM game_item WHERE item_id = ?"
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
