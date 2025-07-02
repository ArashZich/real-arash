package models

import (
	"fmt"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/plugin/soft_delete"
)

type Token struct {
	ID           uint                  `gorm:"primaryKey;uniqueIndex:udx_tokens;"`
	UserID       int                   `json:"user_id"`
	User         *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	CreatedAt    time.Time             `json:"created_at"`
	DeletedAt    soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_tokens"`
}

func (t *Token) isValidAccessToken() bool {
	// Replace this secret with the same secret you used when creating the JWT
	accessSigningKey := []byte(config.AppConfig.AccessTokenSigningKey)

	token, err := jwt.Parse(
		t.AccessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("روش امضای غیرمنتظره: %v", token.Header["alg"])
			}
			return accessSigningKey, nil
		},
	)

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Assuming the "exp" claim stores the expiry date as Unix timestamp
		if f64, ook := claims["ExpiresAt"].(float64); ook {
			expiry := int64(f64)
			// Check if the token is expired
			if time.Unix(expiry, 0).Before(time.Now()) {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}

	return true
}

func (t *Token) IsValidRefreshToken() bool {
	refreshOffset := time.Minute * 5

	// You need to replace this secret with the same secret you used when creating the JWT
	refreshSigningKey := []byte(config.AppConfig.RefreshTokenSigningKey)

	token, err := jwt.Parse(
		t.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("روش امضای غیرمنتظره: %v", token.Header["alg"])
			}
			return refreshSigningKey, nil
		},
	)

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	expiry, ok := claims["ExpiresAt"].(float64)
	if !ok {
		return false
	}

	// Assuming the "exp" claim stores the expiry date as Unix timestamp
	expiryTime := time.Unix(int64(expiry), 0)

	// Check if the token is expired
	// if expiryTime.Before(time.Now().Add(time.Minute * refreshOffset)) {
	// 	return false
	// }

	// return true
	// Check if the token is expired and return the opposite of the result
	return !expiryTime.Before(time.Now().Add(time.Minute * refreshOffset))
}
