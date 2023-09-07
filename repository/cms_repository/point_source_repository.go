package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type pointSourceRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewPointSourceRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IPointSourceRepository {
	return &pointSourceRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *pointSourceRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.PointSource, err error) {
	query := `SELECT point_source_id, point_source_name, created_at, updated_at, created_by, updated_by FROM point_source`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " point_source_name LIKE ?")
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
	var listData []*cms_domain.PointSource

	for results.Next() {
		var data = &cms_domain.PointSource{}
		err = results.Scan(
			&data.PointSourceID,
			&data.PointSourceName,
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
func (ur *pointSourceRepo) GetTotalPointSource(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(point_source_id) as total from point_source`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, " point_source_name LIKE ?")
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

func (ur *pointSourceRepo) CreatePointSource(ctx context.Context, payload *cms_domain.PointSource) (int64, error) {
	query := `INSERT INTO point_source(point_source_name, created_at, created_by)
				VALUES(?, NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.PointSourceName,
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

func (ur *pointSourceRepo) UpdatePointSource(ctx context.Context, payload *cms_domain.PointSource, id int) error {

	query := "UPDATE point_source set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.PointSourceName.String != "" {
		columns = append(columns, " point_source_name = ?")
		params = append(params, payload.PointSourceName)
	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE point_source_id = ?"
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

func (ur *pointSourceRepo) DeletePointSource(ctx context.Context, id int) error {
	query := "DELETE FROM point_source WHERE point_source_id = ?"
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
