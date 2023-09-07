package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type moodRepo struct {
	database   mysql.MysqlClient
	collection string
}

func NewMoodRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IMoodRepository {
	return &moodRepo{
		database:   db,
		collection: collection,
	}
}

func (ur *moodRepo) Get(ctx context.Context, page int, limit int, search string) (users []*cms_domain.Mood, err error) {
	query := `SELECT id, name, created_at, updated_at, created_by, updated_by
       from user_mood`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "name LIKE ?")
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
	var listData []*cms_domain.Mood

	for results.Next() {
		var data = &cms_domain.Mood{}
		err = results.Scan(
			&data.ID,
			&data.Name,
			&data.CreatedAt,
			&data.UpdatedAT,
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
func (ur *moodRepo) GetTotalUserMood(ctx context.Context) (total int, err error) {
	query := `SELECT count(id) as total from user_mood`
	var count int
	stmt, err := ur.database.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	err = stmt.QueryRowContext(ctx).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ur *moodRepo) CreateUserApp(ctx context.Context, payload *cms_domain.Mood) (lastID int64, err error) {
	query := `INSERT INTO user_mood(name, created_at, created_by)
				VALUES(?,NOW(), ?)`
	result, err := ur.database.Conn.Exec(query,
		payload.Name,
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

func (ur *moodRepo) UpdateUserApp(ctx context.Context, payload *cms_domain.Mood, id int) error {

	query := "UPDATE user_mood set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.Name.String != "" {
		columns = append(columns, " name = ?")
		params = append(params, payload.Name)

	}

	if payload.UpdatedBy.Int64 != 0 {
		columns = append(columns, " updated_by = ?")
		params = append(params, payload.UpdatedBy.Int64)

	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE id = ?"
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

func (ur *moodRepo) DeleteUserApp(ctx context.Context, id int) error {
	query := "DELETE FROM user_mood WHERE id = ?"
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
