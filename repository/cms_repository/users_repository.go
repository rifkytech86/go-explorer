package cms_repository

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"strings"
)

type users struct {
	database   mysql.MysqlClient
	collection string
}

func NewUsers(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) cms_domain.IUserAdminRepository {
	return &users{
		database:   db,
		collection: collection,
	}
}

func (ur *users) GetUserAdmin(ctx context.Context, page int, limit int, search string) (users []*cms_domain.Users, err error) {
	query := `SELECT id, email, username, password, name, created_at, updated_at, is_login from user_admin`

	params := []interface{}{}
	var conditions []string
	if search != "" {
		conditions = append(conditions, "email LIKE ?")
		params = append(params, "%"+search+"%")

		conditions = append(conditions, "username LIKE ?")
		params = append(params, "%"+search+"%")

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
	var userAdmins []*cms_domain.Users
	for results.Next() {
		var userAdmin = &cms_domain.Users{}
		err = results.Scan(
			&userAdmin.ID,
			&userAdmin.Email,
			&userAdmin.Username,
			&userAdmin.Password,
			&userAdmin.Name,
			&userAdmin.CreatedAt,
			&userAdmin.UpdatedAt,
			&userAdmin.IsLogin,
		)
		userAdmins = append(userAdmins, userAdmin)
		if err != nil {
			return nil, err
		}
	}
	return userAdmins, nil
}
func (ur *users) GetTotalUserAdmin(ctx context.Context) (total int, err error) {
	query := `SELECT count(id) as total from user_admin`
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

func (ur *users) CreateUserAdmin(ctx context.Context, payload *cms_domain.Users) (lastID int64, err error) {
	query := `INSERT INTO user_admin(email, username, name, password, is_login, created_at, created_by) VALUES(?,?,?,?,?,NOW(), ?)`
	result, err := ur.database.Conn.Exec(query, payload.Email, payload.Username, payload.Name, payload.Password, payload.IsLogin, payload.CreatedBy)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (ur *users) UpdateUserAdmin(ctx context.Context, payload *cms_domain.Users, id int) error {

	query := "UPDATE user_admin set "
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	params := []interface{}{}
	var columns []string
	if payload.Email != "" {
		columns = append(columns, " email = ?")
		params = append(params, payload.Email)

	}
	if payload.Username != "" {
		columns = append(columns, " username = ?")
		params = append(params, payload.Username)
	}
	if payload.Name != "" {
		columns = append(columns, " name = ?")
		params = append(params, payload.Name)
	}

	query += " " + strings.Join(columns, ",") + ", updated_at = NOW() WHERE id = ?"
	params = append(params, id)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(params...)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (ur *users) DeleteUserAdmin(ctx context.Context, id int) error {
	query := "DELETE FROM user_admin WHERE id = ?"
	tx, err := ur.database.Conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (ur *users) GetUserAdminByEmail(ctx context.Context, email string) (user *cms_domain.Users, err error) {
	query := `SELECT id, email, username, password, name, created_at, updated_at, is_login from user_admin WHERE 1=1 `

	params := []interface{}{}
	if email != "" {
		query += " AND email = ? "
		params = append(params, email)
	}

	results := ur.database.Conn.QueryRowContext(ctx, query, params...)

	var userAdmin = &cms_domain.Users{}
	err = results.Scan(
		&userAdmin.ID,
		&userAdmin.Email,
		&userAdmin.Username,
		&userAdmin.Password,
		&userAdmin.Name,
		&userAdmin.CreatedAt,
		&userAdmin.UpdatedAt,
		&userAdmin.IsLogin,
	)
	if err != nil {
		return nil, err
	}

	return userAdmin, nil
}
