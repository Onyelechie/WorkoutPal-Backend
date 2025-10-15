package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) ReadUserByEmail(email string) (*model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow("SELECT id, username, email, password, name, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return &user, nil
}

func (u *userRepository) ReadUsers() ([]*model.User, error) {
	rows, err := u.db.Query("SELECT id, username, email, name, height, height_metric, weight, weight_metric, avatar_url FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		var avatarURL sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
		if err != nil {
			return nil, err
		}
		if avatarURL.Valid {
			user.Avatar = avatarURL.String
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *userRepository) ReadUserByID(id int64) (*model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow("SELECT id, username, email, name, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return &user, nil
}

func (u *userRepository) CreateUser(request model.CreateUserRequest) (*model.User, error) {
	var user model.User
	err := u.db.QueryRow(`
		INSERT INTO users (username, email, password, name, height, height_metric, weight, weight_metric, avatar_url) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, username, email, name, height, height_metric, weight, weight_metric, avatar_url`,
		request.Username, request.Email, request.Password, request.Name,
		request.Height, request.HeightMetric, request.Weight, request.WeightMetric, request.Avatar).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name,
		&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &user.Avatar)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, errors.New("user already exists")
			}
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(request model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow(`
		UPDATE users SET username=$2, email=$3, name=$4, height=$5, height_metric=$6, weight=$7, weight_metric=$8, avatar_url=$9
		WHERE id=$1 RETURNING id, username, email, name, height, height_metric, weight, weight_metric, avatar_url`,
		request.ID, request.Username, request.Email, request.Name,
		request.Height, request.HeightMetric, request.Weight, request.WeightMetric, request.Avatar).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name,
		&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return &user, nil
}

func (u *userRepository) DeleteUser(request model.DeleteUserRequest) error {
	result, err := u.db.Exec("DELETE FROM users WHERE id = $1", request.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
