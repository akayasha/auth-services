package repository

import (
	"auth-services/models"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindValidByHash(hash string) (*models.RefreshToken, error)
	RevokeByUserUUID(userUUID string) error
	Update(token *models.RefreshToken) error
	Revoke(id uint) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func (r *refreshTokenRepository) Revoke(id uint) error {
	//TODO implement me
	return r.db.
		Model(&models.RefreshToken{}).
		Where("user_uuid = ?", id).
		Update("is_revoked", true).Error
}

func (r *refreshTokenRepository) Update(token *models.RefreshToken) error {
	return r.db.Save(token).Error
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindValidByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.
		Where("token_hash = ? AND is_revoked = false AND expires_at > ?", hash, time.Now()).
		First(&token).Error
	return &token, err
}

func (r *refreshTokenRepository) RevokeByUserUUID(userUUID string) error {
	return r.db.
		Model(&models.RefreshToken{}).
		Where("user_uuid = ?", userUUID).
		Update("is_revoked", true).Error
}
