package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/petershaan12/go-auth-clean-arch/internal/token"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/petershaan12/go-auth-clean-arch/resource/model"
	"gorm.io/gorm"
)

type AuthService struct {
	repo       model.UserMethodRepository
	env        library.Env
	tokenMaker token.Maker
}

func NewAuthService(repo model.UserMethodRepository, env library.Env, tokenMaker token.Maker) model.AuthMethodService {
	return &AuthService{
		repo:       repo,
		env:        env,
		tokenMaker: tokenMaker,
	}
}

func (a AuthService) Login(ctx context.Context, req *model.AuthReq) (result *model.User, err error) {
	// cari user (include password)
	filter := []*model.GormWhere{
		{
			Where: "users.email = ? AND users.deleted_at IS NULL",
			Value: []any{req.Email},
		},
	}
	result, err = a.repo.WithContext(ctx).FindBy(filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if result == nil {
		return nil, errors.New("user not found")
	}

	if !library.CheckPasswordHash(string(req.Password), result.Password) {
		return nil, errors.New("invalid credentials")
	}

	return result, nil
}

func (a AuthService) GenerateToken(ctx context.Context, user *model.User) (result *model.TokenOutput, err error) {

	roleStr := strconv.FormatInt(int64(user.RoleId), 10)
	accessTokenExpiry := library.AccessTokenExpiry()
	refreshTokenExpiry := library.RefreshTokenExpiry()

	accessToken, err := a.tokenMaker.CreateToken(
		user.Id,
		user.Email,
		roleStr,
		accessTokenExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := a.tokenMaker.CreateRefreshToken(
		user.Id,
		user.Email,
		roleStr,
		refreshTokenExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	result = &model.TokenOutput{
		BearerType:       "Bearer",
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiredToken:     time.Now().Add(accessTokenExpiry).Format(time.RFC3339),
		UserId:           int(user.Id),
		RequireTwoFactor: false,
	}

	return result, nil
}

func (a AuthService) VerifyRefreshToken(ctx context.Context, req *model.RefreshTokenReq) (*model.TokenOutput, error) {
	payload, err := a.tokenMaker.VerifyToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if payload.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected 'refresh', got '%s'", payload.TokenType)
	}

	newAccessToken, err := a.tokenMaker.CreateToken(
		payload.UserId,
		payload.Email,
		payload.RoleId,
		library.AccessTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	return &model.TokenOutput{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
		ExpiredToken: time.Now().Add(library.AccessTokenExpiry()).Format(time.RFC3339),
		UserId:       int(payload.UserId),
	}, nil
}

func (a AuthService) Logout(ctx context.Context, payload *token.Payload) error {
	return a.repo.WithContext(ctx).IncrementSessionVersion(ctx, payload.UserId)
}
