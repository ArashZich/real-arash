package endpoints

import (
	"context"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/golang-jwt/jwt"
)

func (s *service) generateAccessTokenJWT(identity Identity) (string, error) {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"Audience":  config.AppConfig.AppUrl,
			"ExpiresAt": time.Now().Add(time.Duration(s.tokenExpiration) * time.Minute).Unix(),
			"Id":        identity.GetID(),
			"IssuedAt":  time.Now().Unix(),
			"Issuer":    identity.GetFullName(),
			"NotBefore": time.Now().Unix(),
			"Subject":   identity.GetPhone(),
			"Roles":     identity.GetRoles(),
		},
	).SignedString([]byte(s.accessTokenSigningKey))

	if err != nil {
		s.logger.Error("sign access token", err)
		return "", err
	}
	return token, err
}

func (s *service) generateRefreshTokenJWT(identity Identity) (string, error) {
	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"Audience":  config.AppConfig.AppUrl,
			"ExpiresAt": time.Now().Add(time.Duration(s.refreshTokenExpiration) * time.Minute).Unix(),
			"Id":        identity.GetID(),
			"IssuedAt":  time.Now().Unix(),
			"Issuer":    identity.GetFullName(),
			"NotBefore": time.Now().Unix(),
			"Subject":   identity.GetPhone(),
			"Roles":     identity.GetRoles(),
		},
	).SignedString([]byte(s.refreshTokenSigningKey))

	if err != nil {
		s.logger.Error("sign refresh token", err)
		return "", err
	}
	return refreshToken, err
}

func (s *service) canCreateToken(ctx context.Context, userID uint, offset int) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&models.Token{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		s.logger.Error("count tokens", err)
		return false, err
	}
	return count < int64(config.AppConfig.MaxLoginDeviceCount)+int64(offset), err
}

func (s *service) generateTokens(ctx context.Context, identity Identity, removeCurrentAccessToken string) (
	string, string, response.ErrorResponse,
) {
	var accessToken string
	var refreshToken string
	var user models.User
	_, user, err := s.findUser(ctx, identity.GetPhone())
	if err != nil {
		s.logger.With(ctx).Error(err)
		return accessToken, refreshToken, response.GormErrorResponse(err, "خطایی در سرور رخ داده است")
	}

	if user.SuspendedAt.Valid {
		return accessToken, refreshToken, response.ErrorBadRequest(nil, "حساب کاربری شما مسدود شده است")
	}

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		s.logger.With(ctx).Error(tx.Error)
		return accessToken, refreshToken, response.ErrorInternalServerError(nil, "خطایی در سرور رخ داده است")
	}

	// TODO: later we can get device information to store login history
	userID, er := strconv.Atoi(identity.GetID())
	if er != nil {
		s.logger.Error("convert string to int", er)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError("خطایی در سرور رخ داده است")
	}

	offset := 0
	if removeCurrentAccessToken != "" {
		offset = 1
	}

	canCreateToken, er := s.canCreateToken(ctx, uint(userID), offset)
	if er != nil {
		s.logger.Error("can create token", er)
		tx.Rollback()
		return "", "", response.GormErrorResponse(er, "خطایی در سرور رخ داده است")
	}

	if !canCreateToken {
		if config.AppConfig.AutoDeleteDevice {
			var token models.Token
			err = tx.WithContext(ctx).Where("user_id = ?", userID).Order("created_at asc").First(&token).Error
			if err != nil {
				s.logger.Error("find token", err)
				tx.Rollback()
				return "", "", response.GormErrorResponse(err, "خطایی در سرور رخ داده است")
			}

			err = tx.WithContext(ctx).Delete(&token).Error
			if err != nil {
				s.logger.Error("delete token", err)
				tx.Rollback()
				return "", "", response.GormErrorResponse(err, "خطایی در سرور رخ داده است")
			}
		} else {
			tx.Rollback()
			//	get all tokens from db
			var tokens []models.Token
			s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens)
			return "", "", response.ErrorBadRequest(tokens, "تعداد دستگاه های مجاز شما به پایان رسیده است. لطفا یکی از دستگاه های قبلی را حذف کنید")
		}
	}

	accessToken, err = s.generateAccessTokenJWT(identity)
	if err != nil {
		s.logger.Error("generate access token", err)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError(nil, "خطایی در سرور رخ داده است")
	}

	refreshToken, err = s.generateRefreshTokenJWT(identity)
	if err != nil {
		s.logger.Error("generate refresh token", err)
		tx.Rollback()
		return "", "", response.ErrorInternalServerError(nil, "خطایی در سرور رخ داده است")
	}

	token := models.Token{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err = tx.WithContext(ctx).Create(&token).Error; err != nil {
		s.logger.Error("create token", err)
		tx.Rollback()
		return "", "", response.GormErrorResponse(err, "خطایی در ساخت توکن رخ داده است")
	}

	var currentToken models.Token
	if removeCurrentAccessToken != "" {
		if err = tx.WithContext(ctx).Where("access_token = ?", removeCurrentAccessToken).First(&currentToken).
			Error; err != nil {
			s.logger.Error("find current token", err)
			tx.Rollback()
			return "", "", response.GormErrorResponse(err, "خطایی در ساخت توکن رخ داده است")
		}

		// delete current token
		if err = tx.WithContext(ctx).Delete(&currentToken).Error; err != nil {
			s.logger.Error("delete current token", err)
			tx.Rollback()
			return "", "", response.GormErrorResponse(err, "خطایی در ساخت توکن رخ داده است")
		}
	}

	if err = tx.Commit().Error; err != nil {
		s.logger.Error("commit transaction", err)
		return "", "", response.GormErrorResponse(err, "خطایی در تکمیل ساخت توکن رخ داده است")
	}

	return accessToken, refreshToken, response.ErrorResponse{}
}
