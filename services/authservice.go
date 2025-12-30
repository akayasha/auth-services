package services

import (
	"auth-services/config"
	"auth-services/models"
	"auth-services/repository"
	"auth-services/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var userRepo repository.UserRepository
var refreshRepo repository.RefreshTokenRepository

// initUserRepo ensures the user repository is initialized
func initUserRepo() {
	if userRepo == nil {
		userRepo = repository.NewUserRepository(config.DB)
	}
}

func initRefreshRepo() {
	if refreshRepo == nil {
		refreshRepo = repository.NewRefreshTokenRepository(config.DB)
	}
}

// RegisterUser registers a new user, including email verification and OTP
func RegisterUser(username, email, password, role, firstName, lastName string, dob time.Time) (*models.User, error) {
	initUserRepo()

	if err := models.ValidateRole(models.Role(role)); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP()
	otpHash, _ := utils.HashOTP(otp)
	exp := time.Now().Add(5 * time.Minute)

	user := &models.User{
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.Role(role),
		OTPHash:      otpHash,
		OTPExpiresAt: &exp,
		Dob:          dob,
	}

	if err := userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("email already registered")
	}

	_ = SendOTPEmail(email, otp)

	return user, nil
}

// LoginUser authenticates the user and returns a JWT token
func LoginUser(email, password string) (map[string]interface{}, error) {
	initUserRepo()
	initRefreshRepo()

	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, fmt.Errorf("account locked")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		user.FailedLoginCount++

		if user.FailedLoginCount >= 5 {
			lock := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lock
		}

		_ = userRepo.UpdateUser(user)
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsEmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	// Reset counters
	user.FailedLoginCount = 0
	user.LockedUntil = nil
	_ = userRepo.UpdateUser(user)

	// 🔐 Revoke old refresh tokens
	_ = refreshRepo.RevokeByUserUUID(user.UUID)

	// 🔐 Issue new tokens
	tokens, err := GenerateAuthTokens(user)
	if err != nil {
		return nil, err
	}

	// 🧾 Audit log
	LogAuthEvent(&user.UUID, "LOGIN_SUCCESS")

	return map[string]interface{}{
		"access_token":  tokens["access_token"],
		"refresh_token": tokens["refresh_token"],
	}, nil
}

// ResendOTP resends the OTP for email verification
func ResendOTP(email string) error {
	initUserRepo()

	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return nil // silent
	}

	if user.OTPExpiresAt != nil && time.Until(*user.OTPExpiresAt) > 2*time.Minute {
		return fmt.Errorf("wait before requesting new OTP")
	}

	otp := utils.GenerateOTP()
	hash, _ := utils.HashOTP(otp)
	exp := time.Now().Add(5 * time.Minute)

	user.OTPHash = hash
	user.OTPExpiresAt = &exp

	_ = userRepo.UpdateUser(user)
	return SendOTPEmail(email, otp)
}

// VerifyEmail verifies the OTP entered by the user
func VerifyEmail(email, otp string) (string, error) {
	initUserRepo()

	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid OTP")
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return "", fmt.Errorf("account locked")
	}

	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return "", fmt.Errorf("OTP expired")
	}

	if !utils.VerifyOTP(otp, user.OTPHash) {
		user.FailedOTPCount++

		if user.FailedOTPCount >= 3 {
			lock := time.Now().Add(10 * time.Minute)
			user.LockedUntil = &lock
		}

		_ = userRepo.UpdateUser(user)
		return "", fmt.Errorf("invalid OTP")
	}

	user.IsEmailVerified = true
	user.OTPHash = ""
	user.OTPExpiresAt = nil
	user.FailedOTPCount = 0
	user.LockedUntil = nil

	_ = userRepo.UpdateUser(user)

	return "email verified", nil
}

func GenerateAuthTokens(user *models.User) (map[string]string, error) {
	initRefreshRepo()

	accessToken, err := utils.GenerateJWT(user.UUID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshPlain := uuid.New().String()
	refreshHash, _ := utils.HashToken(refreshPlain)

	rt := &models.RefreshToken{
		UserUUID:  user.UUID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_ = refreshRepo.Create(rt)

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshPlain,
	}, nil
}

func RefreshAccessToken(refreshToken string) (map[string]string, error) {
	initRefreshRepo()
	initUserRepo()

	hash, err := utils.HashToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// ✅ CHECK REDIS BLACKLIST FIRST
	if blacklisted, _ := IsRefreshTokenBlacklisted(hash); blacklisted {
		return nil, fmt.Errorf("refresh token revoked")
	}

	rt, err := refreshRepo.FindValidByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	user, err := userRepo.FindByUUID(rt.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 🔐 BLACKLIST OLD TOKEN
	_ = BlacklistRefreshToken(hash, rt.ExpiresAt)

	// 🔐 REVOKE DB TOKEN
	_ = refreshRepo.Revoke(rt.ID)

	// 🔐 ISSUE NEW TOKENS
	tokens, err := GenerateAuthTokens(user)
	if err != nil {
		return nil, err
	}

	// 🧾 AUDIT
	LogAuthEvent(&user.UUID, "TOKEN_REFRESH")

	return tokens, nil
}

func Logout(refreshToken string) error {
	initRefreshRepo()

	hash, err := utils.HashToken(refreshToken)
	if err != nil {
		return nil
	}

	rt, err := refreshRepo.FindValidByHash(hash)
	if err != nil {
		return nil
	}

	// 🔐 Blacklist in Redis
	_ = BlacklistRefreshToken(hash, rt.ExpiresAt)

	// 🔐 Revoke in DB
	_ = refreshRepo.Revoke(rt.ID)

	return nil
}
