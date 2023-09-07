package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type coinSourceRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewCoinSourceRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.ICoinSourceRepository {
	return &coinSourceRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *coinSourceRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.CoinSource, err error) {
	query := `SELECT coin_source_id, coin_source_name,  created_at, updated_at, created_by, updated_by FROM coin_source`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " coin_source_name LIKE ?")
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
	var listData []*cms_domain.CoinSource

	for results.Next() {
		var data = &cms_domain.CoinSource{}
		err = results.Scan(
			&data.CoinSourceID,
			&data.CoinSourceName,
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
func (ur *coinSourceRepo) GetTotalCoinSource(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(coin_source_id) as total from coin_source`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " coin_source_name LIKE ?")
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

func (ur *coinSourceRepo) CreateCoinSource(ctx context.Context, payload *cms_domain.CoinSource) (int64, error) {
	query := `INSERT INTO coin_source(coin_source_name, created_at, created_by)
				VALUES(?, NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.CoinSourceName,
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

func (ur *coinSourceRepo) UpdateCoinSource(ctx context.Context, payload *cms_domain.CoinSource, id int) error {

	query := "UPDATE coin_source set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.CoinSourceName.String != "" {
		columns = append(columns, " coin_source_name = ?")
		params = append(params, payload.CoinSourceName)
	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE coin_source_id = ?"
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

func (ur *coinSourceRepo) DeleteCoinSource(ctx context.Context, id int) error {
	query := "DELETE FROM coin_source WHERE coin_source_id = ?"
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
