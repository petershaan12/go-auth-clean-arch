package repository

import (
	"context"
	"time"

	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  library.Database
	ctx context.Context
}

func NewUserRepository(db library.Database) model.UserMethodRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) baseQuery() *gorm.DB {
	return u.db.DB.Table(model.UserTable).Where("users.deleted_at IS NULL")
}

func (u *UserRepository) WithContext(ctx context.Context) model.UserMethodRepository {
	return &UserRepository{
		db:  u.db,
		ctx: ctx,
	}
}
func (u *UserRepository) Count(conditions []*model.GormWhere) (total int64, err error) {
	query := u.baseQuery()
	for _, condition := range conditions {
		query = query.Where(condition.Where, condition.Value...)
	}
	query = query.Count(&total)

	return total, query.Error
}

func (u *UserRepository) FindAllBy(conditions []*model.GormWhere) (result []*model.User, err error) {
	query := u.baseQuery()

	for _, condition := range conditions {
		query = query.Where(condition.Where, condition.Value...)
	}
	query = query.Order("created_at DESC")

	query = u.queryJoinHidePassword(query).
		Debug().Find(&result)

	return result, query.Error
}

func (u *UserRepository) queryJoinHidePassword(query *gorm.DB) *gorm.DB {
	return query.Select("users.id, users.username, users.email, users.fullname, users.role_id, users.created_at, users.updated_at, users.deleted_at, roles.name as role_name").
		Joins("left join roles on roles.id = users.role_id").Where("users.deleted_at IS NULL")
}

func (u *UserRepository) Create(data *model.User) (result *model.User, err error) {
	query := u.db.DB.Debug().Table(model.UserTable).Create(data)
	if query.Error != nil {
		return nil, query.Error
	}

	data.Password = ""
	return data, nil
}

func (u *UserRepository) Update(data *model.User) (result *model.User, err error) {
	query := u.db.DB.Debug().Table(model.UserTable).Where("id = ? AND deleted_at IS NULL", data.Id).Save(data)
	if query.Error != nil {
		return nil, query.Error
	}
	data.Password = ""
	return data, nil
}

func (u *UserRepository) Delete(id int64) (err error) {
	now := time.Now()
	query := u.db.DB.Debug().Table(model.UserTable).Where("id = ?", id).Update("deleted_at", now)
	return query.Error
}

func (u *UserRepository) FindBy(filter []*model.GormWhere) (result *model.User, err error) {
	query := u.baseQuery()
	// apply filters
	for _, f := range filter {
		if f == nil {
			continue
		}
		if len(f.Value) > 0 {
			query = query.Where(f.Where, f.Value...)
		} else {
			query = query.Where(f.Where)
		}
	}
	err = query.First(&result).Error
	return
}

func (u *UserRepository) IncrementSessionVersion(ctx context.Context, userId int64) error {
	return u.db.DB.WithContext(ctx).
		Table(model.UserTable).
		Where("id = ? AND deleted_at IS NULL", userId).
		Update("session_version", gorm.Expr("session_version + 1")).Error
}

func (u *UserRepository) GetSessionVersion(ctx context.Context, userId int64) (int, error) {
	var user model.User
	err := u.db.DB.WithContext(ctx).
		Table(model.UserTable).
		Select("session_version").
		Where("id = ? AND deleted_at IS NULL", userId).
		First(&user).Error

	if err != nil {
		return 0, err
	}

	return user.SessionVersion, nil
}
