package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

// Custom type for handling BYTEA data in PostgreSQL
type ByteaData []byte

// Value implements driver.Valuer interface for PostgreSQL BYTEA
func (b ByteaData) Value() (driver.Value, error) {
	if len(b) == 0 {
		return nil, nil
	}
	// Return the raw bytes - pq driver will handle the BYTEA encoding
	return []byte(b), nil
}

// Scan implements sql.Scanner interface for PostgreSQL BYTEA
func (b *ByteaData) Scan(value interface{}) error {
	if value == nil {
		*b = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		*b = ByteaData(v)
		return nil
	case string:
		// Handle hex-encoded BYTEA (\\x format)
		if strings.HasPrefix(v, "\\x") {
			decoded, err := hex.DecodeString(v[2:])
			if err != nil {
				return err
			}
			*b = ByteaData(decoded)
			return nil
		}
		*b = ByteaData(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into ByteaData", value)
	}
}

// Helper function to convert binary data to data URL with correct MIME type detection
func binaryToDataURL(binaryData []byte) string {
	if len(binaryData) == 0 {
		return ""
	}
	
	// Detect image type from binary data
	var mimeType string
	if len(binaryData) >= 8 && binaryData[0] == 0x89 && string(binaryData[1:4]) == "PNG" {
		mimeType = "image/png"
	} else if len(binaryData) >= 2 && binaryData[0] == 0xFF && binaryData[1] == 0xD8 {
		mimeType = "image/jpeg"
	} else {
		mimeType = "image/jpeg" // default fallback
	}
	
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(binaryData)
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) ReadUserByEmail(email string) (*model.User, error) {
	var user model.User
	var avatarData ByteaData
	err := u.db.QueryRow("SELECT id, username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Age, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarData, &user.IsPrivate, &user.ShowMetricsToFollowers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	// Convert binary data back to data URL for frontend compatibility
	user.Avatar = binaryToDataURL([]byte(avatarData))
	return &user, nil
}

func (u *userRepository) ReadUsers() ([]*model.User, error) {
	rows, err := u.db.Query("SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		var avatarData ByteaData
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Age, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarData, &user.IsPrivate, &user.ShowMetricsToFollowers)
		if err != nil {
			return nil, err
		}
		// Convert binary data back to data URL for frontend compatibility
		user.Avatar = binaryToDataURL([]byte(avatarData))
		users = append(users, &user)
	}
	return users, nil
}

func (u *userRepository) ReadUserByID(id int64) (*model.User, error) {
	var user model.User
	var avatarData ByteaData
	err := u.db.QueryRow("SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Age, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarData, &user.IsPrivate, &user.ShowMetricsToFollowers)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	// Convert binary data back to data URL for frontend compatibility
	user.Avatar = binaryToDataURL([]byte(avatarData))
	return &user, nil
}

func (u *userRepository) CreateUser(request model.CreateUserRequest) (*model.User, error) {
	var user model.User
	
	// Convert data URL to binary data for storage
	var avatarBinary []byte
	if request.Avatar != "" {
		if strings.HasPrefix(request.Avatar, "data:image/") {
			// Extract base64 data from data URL
			commaIndex := strings.Index(request.Avatar, ",")
			if commaIndex != -1 {
				base64Data := request.Avatar[commaIndex+1:]
				var err error
				avatarBinary, err = base64.StdEncoding.DecodeString(base64Data)
				if err != nil {
					return nil, errors.New("invalid base64 avatar data")
				}
			}
		}
	}
	
	var avatarData ByteaData
	avatarParam := ByteaData(avatarBinary)
	err := u.db.QueryRow(`
		INSERT INTO users (username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::bytea, COALESCE($11, FALSE), COALESCE($12, FALSE)) 
		RETURNING id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers`,
				request.Username, request.Email, request.Password, request.Name, request.Age,
				request.Height, request.HeightMetric, request.Weight, request.WeightMetric, avatarParam, request.IsPrivate, request.ShowMetricsToFollowers).Scan(
				&user.ID, &user.Username, &user.Email, &user.Name, &user.Age,
				&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarData, &user.IsPrivate, &user.ShowMetricsToFollowers)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, errors.New("user already exists")
			}
		}
		return nil, err
	}
	
	// Convert binary data back to data URL for response
	user.Avatar = binaryToDataURL([]byte(avatarData))
	return &user, nil
}

func (u *userRepository) UpdateUser(request model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	
	// Convert data URL to binary data for storage
	var avatarBinary []byte
	if request.Avatar != "" {
		if strings.HasPrefix(request.Avatar, "data:image/") {
			// Extract base64 data from data URL
			commaIndex := strings.Index(request.Avatar, ",")
			if commaIndex != -1 {
				base64Data := request.Avatar[commaIndex+1:]
				var err error
				avatarBinary, err = base64.StdEncoding.DecodeString(base64Data)
				if err != nil {
					return nil, errors.New("invalid base64 avatar data")
				}
			}
		}
	}
	
	var avatarData ByteaData
	
	// Use ByteaData for proper PostgreSQL BYTEA handling
	avatarParam := ByteaData(avatarBinary)
	
	err := u.db.QueryRow(`
		UPDATE users SET username=$2, email=$3, name=$4, age=$5, height=$6, height_metric=$7, weight=$8, weight_metric=$9, avatar_data=$10::bytea, is_private=$11, show_metrics_to_followers=$12
		WHERE id=$1 RETURNING id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers`,
				request.ID, request.Username, request.Email, request.Name, request.Age,
				request.Height, request.HeightMetric, request.Weight, request.WeightMetric, avatarParam, request.IsPrivate, request.ShowMetricsToFollowers).Scan(
				&user.ID, &user.Username, &user.Email, &user.Name, &user.Age,
				&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarData, &user.IsPrivate, &user.ShowMetricsToFollowers)
	if err != nil {
		log.Printf("UpdateUser: SQL query error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	// Convert binary data back to data URL for response
	user.Avatar = binaryToDataURL([]byte(avatarData))
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
