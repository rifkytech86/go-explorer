package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strings"
)

type usersApp struct {
	database   mysql.MysqlClient
	collection string
}

func NewUsersAppRepo(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IUserAppRepository {
	return &usersApp{
		database:   db,
		collection: collection,
	}
}

func (ur *usersApp) GetUserApp(ctx context.Context, page int, limit int, search string) (users []*cms_domain.UsersApp, err error) {
	query := `SELECT user_id, user_name, 
       	user_phone, user_is_verify, user_otp, 
       	updated_at, created_at, 
       	created_by, updated_by from user`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "user_name LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "user_phone LIKE ?")
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
	var userApps []*cms_domain.UsersApp

	for results.Next() {
		var userApp = &cms_domain.UsersApp{}
		err = results.Scan(
			&userApp.UserID,
			&userApp.UserName,
			&userApp.UserPhone,
			&userApp.UserISVerify,
			&userApp.UserOTP,
			&userApp.UserUpdatedAt,
			&userApp.UserCreatedAt,
			&userApp.UserCreatedBy,
			&userApp.UserUpdatedBy,
		)
		userApps = append(userApps, userApp)
		if err != nil {
			return nil, err
		}
	}
	return userApps, nil
}
func (ur *usersApp) GetTotalUserApp(ctx context.Context) (total int, err error) {
	query := `SELECT count(user_id) as total from user`
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

func (ur *usersApp) CreateUserApp(ctx context.Context, payload *cms_domain.UsersApp) (lastID int64, err error) {
	query := `INSERT INTO user(
                 user_name, user_phone, user_is_verify, created_by,
                 created_at) 
				VALUES(?,?,?,?,NOW())`

	result, err := ur.database.Conn.Exec(query,
		payload.UserName.String,
		payload.UserPhone.String,
		payload.UserISVerify.Int64,
		payload.UserCreatedBy.Int64,
	)
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

func (ur *usersApp) UpdateUserApp(ctx context.Context, payload *cms_domain.UsersApp, id int) error {

	query := "UPDATE user set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.UserName.String != "" {
		columns = append(columns, " user_name = ?")
		params = append(params, payload.UserName.String)

	}

	if payload.UserPhone.String != "" {
		columns = append(columns, " user_phone = ?")
		params = append(params, payload.UserPhone.String)
	}

	columns = append(columns, " user_is_verify = ?")
	params = append(params, payload.UserISVerify.Int64)

	columns = append(columns, " updated_by = ?")
	params = append(params, payload.UserUpdatedBy.Int64)

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE user_id = ?"
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

func (ur *usersApp) DeleteUserApp(ctx context.Context, id int) error {
	query := "DELETE FROM user WHERE user_id = ?"
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

func (ur *usersApp) GetUserAppByPhone(ctx context.Context, phone string) (user *cms_domain.UsersApp, err error) {
	query := `SELECT user_id, user_name, user_phone, user_is_verify, user_otp, created_at, created_by, updated_at, updated_by from user WHERE 1=1 `

	params := []interface{}{}
	if phone != "" {
		query += " AND user_phone = ? "
		params = append(params, phone)
	}

	results := ur.database.Conn.QueryRowContext(ctx, query, params...)

	var setUser = &cms_domain.UsersApp{}
	err = results.Scan(
		&setUser.UserID,
		&setUser.UserName,
		&setUser.UserPhone,
		&setUser.UserISVerify,
		&setUser.UserOTP,
		&setUser.UserCreatedAt,
		&setUser.UserCreatedBy,
		&setUser.UserUpdatedAt,
		&setUser.UserUpdatedBy,
	)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	return setUser, nil
}
