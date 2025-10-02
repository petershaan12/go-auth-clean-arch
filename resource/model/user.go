package model

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type (
	UserHasChangePassword string
)

const (
	UserTable     = "users"
	UserDeletedAt = "deleted_at"
)

type (
	User struct {
		Id             int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
		Username       string     `json:"username" gorm:"type:varchar(100)"`
		Email          string     `json:"email" gorm:"type:varchar(100)"`
		Fullname       string     `json:"fullname" gorm:"type:varchar(255)"`
		RoleId         int64      `json:"role_id" gorm:"type:int(11)"`
		SessionVersion int        `json:"session_version" gorm:"column:session_version;default:1"`
		Password       string     `json:"password,omitempty" gorm:"type:varchar(255)"`
		CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
		DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	}

	CreateUserRequest struct {
		Username string `json:"username" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email,max=100"`
		Fullname string `json:"fullname" validate:"required,max=255"`
		RoleID   int64  `json:"role_id" validate:"required"`
		Password string `json:"password" validate:"required,min=6,max=128"`
	}

	UpdateUserRequest struct {
		Username string `json:"username" validate:"omitempty,min=3,max=100"`
		Email    string `json:"email" validate:"omitempty,email,max=100"`
		Fullname string `json:"fullname" validate:"omitempty,max=255"`
		RoleID   int64  `json:"role_id" validate:"omitempty"`
		Password string `json:"password" validate:"omitempty,min=6,max=128"`
	}

	UserMethodRepository interface {
		WithContext(ctx context.Context) UserMethodRepository
		FindAllBy(filter []*GormWhere) (result []*User, err error)
		FindBy(filter []*GormWhere) (result *User, err error)
		Count(filter []*GormWhere) (total int64, err error)
		Create(req *User) (result *User, err error)
		Update(data *User) (result *User, err error)
		Delete(id int64) (err error)
		IncrementSessionVersion(ctx context.Context, userId int64) error
		GetSessionVersion(ctx context.Context, userId int64) (int, error)
	}

	UserMethodService interface {
		List(ctx context.Context) (result []*User, err error)
		Create(tx *gorm.DB, req *CreateUserRequest) (result *User, err error)
		Update(c echo.Context, tx *gorm.DB, id int, req *UpdateUserRequest) (result *User, err error)
		Delete(c echo.Context, tx *gorm.DB, id int) (err error)
		GetByID(c echo.Context, id int) (result *User, err error)
	}
)
