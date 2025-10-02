package model

import (
	"context"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"github.com/petershaan12/go-auth-clean-arch/internal/token"
)

type (
	AuthReq struct {
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=6,max=128"`
		IPAddress string `json:"ip_address,omitempty"`
	}

	RefreshTokenReq struct {
		AccessToken  string `json:"at" query:"at"`
		RefreshToken string `json:"rt" query:"rt"`
	}

	TokenOutput struct {
		BearerType       string `json:"b"`
		AccessToken      string `json:"a"`
		RefreshToken     string `json:"r"`
		ExpiredToken     string `json:"e"`
		UserId           int    `json:"user_id"`
		RequireTwoFactor bool   `json:"require_two_factor"`
		Enable2FA        bool   `json:"enable_2fa"`
	}

	AuthMethodService interface {
		Login(ctx context.Context, req *AuthReq) (result *User, err error)
		GenerateToken(ctx context.Context, user *User) (result *TokenOutput, err error)
		VerifyRefreshToken(ctx context.Context, req *RefreshTokenReq) (result *TokenOutput, err error)
		Logout(ctx context.Context, payload *token.Payload) error
	}
)

func ValidatePassword(password string) error {
	// Decode base64 input
	passBytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return err
	}
	pass := string(passBytes)

	if len(pass) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	// simple regex checks for each required class
	hasNumber, _ := regexp.MatchString(`[0-9]`, pass)
	hasUpper, _ := regexp.MatchString(`[A-Z]`, pass)
	hasLower, _ := regexp.MatchString(`[a-z]`, pass)
	hasSpecial, _ := regexp.MatchString(`[!@#\$%\^&\*\(\)\-_=+\[\]\{\};:'",.<>/?]`, pass)

	if hasNumber && hasUpper && hasLower && hasSpecial {
		return nil
	}

	var b strings.Builder
	b.WriteString("Password must contain ")
	if !hasNumber {
		b.WriteString("at least one number, ")
	}
	if !hasUpper {
		b.WriteString("at least one uppercase letter, ")
	}
	if !hasLower {
		b.WriteString("at least one lowercase letter, ")
	}
	if !hasSpecial {
		b.WriteString("at least one special character, ")
	}

	return errors.New(strings.TrimSuffix(b.String(), ", ") + ".")
}
