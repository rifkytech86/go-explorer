package cms_repository

import (
	"context"
	"database/sql"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
)

type authorization struct {
	database   mysql.MysqlClient
	collection string
}

func NewAuthorization(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) domain.IUserRepository {
	return &authorization{
		database:   db,
		collection: collection,
	}
}

func (ur *authorization) GetUser(ctx context.Context, email string) (user dto.Users, err error) {
	query := `SELECT id, email, username, password, name, created_at, updated_at from user_admin where email = ?`
	err = ur.database.Conn.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		logger.Error(err.Error(), nil)
		if err == sql.ErrNoRows {
			return user, sql.ErrNoRows
		} else {
			return user, err
		}
	}

	return user, nil
}

func (ur *authorization) UpdateIsLogin(ctx context.Context, email string, isLogin int) (err error) {
	query := `UPDATE user_admin set is_login = ? where email = ?`

	_, err = ur.database.Conn.Exec(query, isLogin, email)
	if err != nil {
		return err
	}
	return nil
}
