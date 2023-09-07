package repository

import (
	"context"
	"database/sql"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/dto/mysql"
)

type diaryRepository struct {
	database   mysql.MysqlClient
	collection string
}

func NewDiaryRepository(db mysql.MysqlClient, collection string) domain.DiaryRepository {
	return &diaryRepository{
		database:   db,
		collection: collection,
	}
}

func (dr diaryRepository) FetchDiaryRepository(ctx context.Context, diary *domain.DiaryRequest) ([]*domain.Diary, error) {
	query := `SELECT
				diary_id,
				user_id,
				diary_desc,
				diary_link_video,
				created_at,
				updated_at,
				diary_is_active
				from diary  `

	stmt, err := dr.database.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []*domain.Diary{}, err
	}

	result := make([]*domain.Diary, 0)
	for rows.Next() {
		t := new(domain.Diary)
		err = rows.Scan(
			&t.DiaryID,
			&t.UserID,
			&t.DiaryDesc,
			&t.DiaryLinkVideo,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DiaryIsActive,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (dr diaryRepository) CreateDiaryRepository(ctx context.Context, diary *domain.DiaryReq) (lastID int64, err error) {
	tx, err := dr.database.Conn.Begin()
	if err != nil {
		return 0, err
	}
	// Prepare the SQL statement
	stmt, err := tx.Prepare("INSERT INTO diary (user_id, diary_desc, diary_link_video, diary_is_active) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		diary.UserID,
		diary.DiaryDesc,
		diary.DiaryLinkVideo,
		1,
	)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil

}

func (dr diaryRepository) FindUserRepository(ctx context.Context, userEmail string) (domain.User, error) {
	var data domain.User
	query := `SELECT user_id, user_name, user_email, user_dob, user_school, user_is_verify, updated_at, created_at  from user where user_email = ? `

	err := dr.database.Conn.QueryRowContext(ctx, query, userEmail).Scan(
		&data.UserID,
		&data.UserName,
		&data.UserEmail,
		&data.UserDOB,
		&data.UserSchool,
		&data.UserIsVerify,
		&data.UpdatedAt,
		&data.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, sql.ErrNoRows
		} else {
			return data, err
		}
	}
	return data, nil
}

func (dr diaryRepository) FindUserByIDRepository(ctx context.Context, userID string) (domain.User, error) {
	var data domain.User
	query := `SELECT user_id, user_name, user_email, user_dob, user_school, user_is_verify, updated_at, created_at  from user where user_id = ? `

	err := dr.database.Conn.QueryRowContext(ctx, query, userID).Scan(
		&data.UserID,
		&data.UserName,
		&data.UserEmail,
		&data.UserDOB,
		&data.UserSchool,
		&data.UserIsVerify,
		&data.UpdatedAt,
		&data.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, sql.ErrNoRows
		} else {
			return data, err
		}
	}
	return data, nil
}

func (dr diaryRepository) UpdateDiaryRepository(c context.Context, diary *domain.DiaryReq) error {
	// Create a new transaction
	tx, err := dr.database.Conn.Begin()
	if err != nil {
		return err
	}
	// Prepare the SQL statement
	stmt, err := tx.Prepare("UPDATE diary set diary_desc = ?, diary_link_video = ?, updated_at = NOW() where diary_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		diary.DiaryDesc,
		diary.DiaryLinkVideo,
		diary.UserID,
	)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		return err
	}
	// Commit the transaction if no error occurred
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (dr diaryRepository) DeleteDiaryRepository(c context.Context, diaryID string) error {
	// Create a new transaction
	tx, err := dr.database.Conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE diary set diary_is_active = ?, updated_at = NOW() where diary_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(
		0,
		diaryID,
	)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		return err
	}

	// Commit the transaction if no error occurred
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
