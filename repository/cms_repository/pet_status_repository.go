package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type petStatusRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewPetStatusRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IPetStatusRepository {
	return &petStatusRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *petStatusRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.PetStatus, err error) {
	query := `SELECT pet_status_id, pet_status_name, created_at, updated_at, created_by, updated_by FROM game_pet_status`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "pet_status_name LIKE ?")
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
	var listData []*cms_domain.PetStatus

	for results.Next() {
		var data = &cms_domain.PetStatus{}
		err = results.Scan(
			&data.PetStatusID,
			&data.PetStatusName,
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
func (ur *petStatusRepo) GetTotalPetStatus(ctx context.Context, search string) (total int, err error) {

	query := `SELECT count(pet_status_id) as total from game_pet_status`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "pet_status_name LIKE ?")
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

func (ur *petStatusRepo) CreatePetStatus(ctx context.Context, payload *cms_domain.PetStatus) (int64, error) {
	query := `INSERT INTO game_pet_status(pet_status_name, created_at, created_by)
				VALUES(?,NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.PetStatusName,
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

func (ur *petStatusRepo) UpdatePetStatus(ctx context.Context, payload *cms_domain.PetStatus, id int) error {

	query := "UPDATE game_pet_status set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.PetStatusName.String != "" {
		columns = append(columns, " pet_status_name = ?")
		params = append(params, payload.PetStatusName)

	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE pet_status_id = ?"
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

func (ur *petStatusRepo) DeletePetStatus(ctx context.Context, id int) error {
	query := "DELETE FROM game_pet_status WHERE pet_status_id = ?"
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
