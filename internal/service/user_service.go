package service

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/model"
	"gorm.io/gorm"
)

type UserService struct {
	repo model.UserMethodRepository
	env  library.Env
}

func NewUserService(repo model.UserMethodRepository, env library.Env) model.UserMethodService {
	return &UserService{
		repo: repo,
		env:  env,
	}
}

func (u *UserService) List(ctx context.Context) (result []*model.User, err error) {
	var filter []*model.GormWhere
	result, err = u.repo.WithContext(ctx).FindAllBy(filter)
	if err != nil {
		return
	}
	return
}

func (u *UserService) Create(tx *gorm.DB, req *model.CreateUserRequest) (*model.User, error) {

	// 1) Cek duplikasi email
	ctx := tx.Statement.Context
	emailFilter := []*model.GormWhere{
		{Where: "users.email = ?", Value: []any{req.Email}},
	}
	emailTotal, err := u.repo.WithContext(ctx).Count(emailFilter)
	if err != nil {
		return nil, err
	}
	if emailTotal > 0 {
		return nil, errors.New("email already in use")
	}

	// 2) Cek duplikasi username
	usernameFilter := []*model.GormWhere{
		{Where: "users.username = ?", Value: []any{req.Username}},
	}
	usernameTotal, err := u.repo.WithContext(ctx).Count(usernameFilter)
	if err != nil {
		return nil, err
	}
	if usernameTotal > 0 {
		return nil, errors.New("username already in use")
	}

	// 3) Hash password
	hashPass, err := library.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 4) Map ke entity
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Fullname: req.Fullname,
		RoleId:   req.RoleID,
		Password: hashPass,
	}

	// 5) Simpan
	created, err := u.repo.WithContext(ctx).Create(user)
	if err != nil {
		return nil, err
	}

	// 6) Sanitasi sebelum return (kalau field sensitif belum ditandai json:"-")
	user.Password = ""

	return created, nil
}

func (u *UserService) Update(c echo.Context, tx *gorm.DB, id int, req *model.UpdateUserRequest) (*model.User, error) {
	// 1) Cek user exists
	ctx := tx.Statement.Context
	id64 := int64(id)
	userFilter := []*model.GormWhere{
		{Where: "users.id = ?", Value: []any{id}},
	}
	userTotal, err := u.repo.WithContext(ctx).Count(userFilter)
	if err != nil {
		return nil, err
	}
	if userTotal == 0 {
		return nil, errors.New("user not found")
	}

	// 2) Cek duplikasi email jika email diupdate
	if req.Email != "" {
		emailFilter := []*model.GormWhere{
			{Where: "users.email = ? AND users.id != ?", Value: []any{req.Email, id}},
		}
		emailTotal, err := u.repo.WithContext(ctx).Count(emailFilter)
		if err != nil {
			return nil, err
		}
		if emailTotal > 0 {
			return nil, errors.New("email already in use")
		}
	}

	// 3) Cek duplikasi username jika username diupdate
	if req.Username != "" {
		usernameFilter := []*model.GormWhere{
			{Where: "users.username = ? AND users.id != ?", Value: []any{req.Username, id}},
		}
		usernameTotal, err := u.repo.WithContext(ctx).Count(usernameFilter)
		if err != nil {
			return nil, err
		}
		if usernameTotal > 0 {
			return nil, errors.New("username already in use")
		}
	}

	// 4) Handle password update jika disediakan
	var pwdHash string
	if req.Password != "" {
		hashPass, err := library.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		pwdHash = hashPass
	}

	// 5) Ambil user existing untuk update
	existingUser, err := u.repo.WithContext(ctx).FindBy([]*model.GormWhere{
		{Where: "users.id = ?", Value: []any{id64}},
	})
	if err != nil {
		return nil, err
	}

	// 6) Update fields
	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Fullname != "" {
		existingUser.Fullname = req.Fullname
	}

	if req.RoleID != 0 {
		existingUser.RoleId = req.RoleID
	}
	if pwdHash != "" {
		existingUser.Password = pwdHash
	}

	// 7) Simpan update
	updated, err := u.repo.WithContext(ctx).Update(existingUser)
	if err != nil {
		return nil, err
	}

	// 8) Sanitasi sebelum return
	updated.Password = ""

	return updated, nil
}

func (u *UserService) Delete(c echo.Context, tx *gorm.DB, id int) error {
	// 1) Cek user exists
	ctx := tx.Statement.Context
	id64 := int64(id)
	userFilter := []*model.GormWhere{
		{Where: "users.id = ?", Value: []any{id}},
	}
	userTotal, err := u.repo.WithContext(ctx).Count(userFilter)
	if err != nil {
		return err
	}
	if userTotal == 0 {
		return errors.New("user not found")
	}

	// 2) Hapus user
	err = u.repo.WithContext(ctx).Delete(id64)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetByID(c echo.Context, id int) (*model.User, error) {
	id64 := int64(id)
	user, err := u.repo.WithContext(c.Request().Context()).FindBy([]*model.GormWhere{
		{Where: "users.id = ?", Value: []any{id64}},
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
