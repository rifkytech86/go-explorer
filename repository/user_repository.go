package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"

	"gitlab.com/naonweh-studio/bubbme-backend/domain"
)

type userRepository struct {
	database   mysql.MysqlClient
	collection string
	cache      *bootstrap.RedisClient
}

func NewUserRepository(db mysql.MysqlClient, cache *bootstrap.RedisClient, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
		cache:      cache,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) (int64, error) {
	query := `INSERT INTO user(user_name, user_phone, user_otp, user_is_verify, created_at) VALUES(?,?,?,0, NOW())`
	result, err := ur.database.Conn.Exec(query, user.UserName.String, user.UserPhone.String, user.UserOTP.Int64)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	//collection := ur.database.Collection(ur.collection)
	//
	//opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	//cursor, err := collection.Find(c, bson.D{}, opts)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//var users []domain.User
	//
	//err = cursor.All(c, &users)
	//if users == nil {
	//	return []domain.User{}, err
	//}
	//
	//return users, err
	return nil, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//collection := ur.database.Collection(ur.collection)
	//var user domain.User
	//err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	//return user, err
	query := `SELECT user_id, user_name, user_email, user_dob, user_school, user_is_verify, updated_at, created_at  from user where user_email = ? `

	row := ur.database.Conn.QueryRowContext(ctx, query, email)
	t := new(domain.User)
	err := row.Scan(
		&t.UserID,
		&t.UserName,
		&t.UserEmail,
		&t.UserDOB,
		&t.UserSchool,
		&t.UserIsVerify,
		&t.UpdatedAt,
		&t.CreatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return domain.User{}, errors.New("user not found")
	case err != nil:
		return domain.User{}, err
	default:
		return *t, nil
	}
}

func (ur *userRepository) GetUserByPhone(ctx context.Context, phoneNumber string) (domain.User, error) {
	query := `SELECT user_id, user_name, user_email, user_phone, user_dob, user_school, user_is_verify, updated_at, created_at  from user where user_phone = ? `

	row := ur.database.Conn.QueryRowContext(ctx, query, phoneNumber)
	t := new(domain.User)
	err := row.Scan(
		&t.UserID,
		&t.UserName,
		&t.UserEmail,
		&t.UserPhone,
		&t.UserDOB,
		&t.UserSchool,
		&t.UserIsVerify,
		&t.UpdatedAt,
		&t.CreatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		logger.Error(sql.ErrNoRows.Error(), nil)
		return domain.User{}, sql.ErrNoRows
	case err != nil:
		logger.Error(err.Error(), nil)
		return domain.User{}, err
	default:
		return *t, nil
	}
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	//collection := ur.database.Collection(ur.collection)
	//
	//var user domain.User
	//
	//idHex, err := primitive.ObjectIDFromHex(id)
	//if err != nil {
	//	return user, err
	//}
	//
	//err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	//return user, err
	return domain.User{}, nil
}

func (ur *userRepository) SetExpiredToken(c context.Context, param domain.CacheOTP) error {
	err := ur.cache.SetRedis(c, fmt.Sprintf("otp:%s", param.Phone), param.OTP, 3)
	return err
}

func (ur *userRepository) GetCacheVerify(c context.Context, param domain.VerifyOTP) (*domain.CacheOTP, error) {
	result, err := ur.cache.GetRedis(c, fmt.Sprintf("otp:%s", param.UserPhone))
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("data not found")
		}
		return nil, err
	}

	var otpCache domain.CacheOTP
	x := []byte(result.(string))
	err = json.Unmarshal(x, &otpCache)
	if err != nil {
		return nil, err
	}

	return &otpCache, err
}

func (ur *userRepository) DeleteCacheVerify(c context.Context, phone string) error {
	err := ur.cache.DelRedis(c, fmt.Sprintf("otp:%s", phone))
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserVerify(c context.Context, userID int64, isVerify int) (int64, error) {
	query := `UPDATE user set user_is_verify = ? where user_id = ?`
	result, err := ur.database.Conn.Exec(query, isVerify, userID)
	if err != nil {
		logger.Error(err.Error(), nil)
		return 0, err
	}
	return result.RowsAffected()
}
func (ur *userRepository) UpdateOTPUserByPhone(c context.Context, phoneNumber string, otp int) error {
	query := `UPDATE user set user_otp = ? where user_phone = ?`
	_, err := ur.database.Conn.Exec(query, otp, phoneNumber)
	if err != nil {
		logger.Error(err.Error(), nil)
		return err
	}
	return nil
}
