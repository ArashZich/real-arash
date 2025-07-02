package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"

	"gitag.ir/armogroup/armo/services/reality/notification"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, input LoginRequest) (response LoginResponse, err response.ErrorResponse)
	Register(ctx context.Context, input RegisterRequest) (response LoginResponse, err response.ErrorResponse)
	ResetPassword(ctx context.Context, input ResetPasswordRequest) (response LoginResponse, err response.ErrorResponse)
	ChangePassword(ctx context.Context, input ChangePasswordRequest) (
		response LoginResponse, err response.ErrorResponse,
	)
	Impersonate(ctx context.Context, id string, currentAccessToken string) (response LoginResponse, err response.ErrorResponse)
	UserInfo(ctx context.Context, accessToken string) (user models.User, err response.ErrorResponse)
	ApprovePhone(ctx context.Context, code string) (ok string, err response.ErrorResponse)
	ApproveEmail(ctx context.Context, token string) (ok string, err response.ErrorResponse)
	DeleteTokens(ctx context.Context, accessTokens []string) (response []string, err error)
	Logout(ctx context.Context, accessTokens string) (response string, err response.ErrorResponse)
	RefreshToken(ctx context.Context, input RefreshTokenRequest) (tokens RefreshTokenResponse, err response.ErrorResponse)
	GetAllTokensByUserId(ctx context.Context, userID string) (tokens []models.Token, err response.ErrorResponse)
	authenticate(ctx context.Context, username, password string) Identity
	generateAccessTokenJWT(identity Identity) (string, error)
	generateRefreshTokenJWT(identity Identity) (string, error)
	generateTokens(ctx context.Context, identity Identity, removeCurrentAccessToken string) (
		accessToken string, refreshToken string, err response.ErrorResponse,
	)
	canCreateToken(ctx context.Context, userID uint, offset int) (ok bool, err error)
	checkAndDeleteVerificationBySessionCodeAndPhone(ctx context.Context, sessionCode string, phone string) (err response.ErrorResponse)
	checkAndDeleteVerificationByCode(ctx context.Context, code string) (phone string, err response.ErrorResponse)
	useInviteCode(ctx context.Context, code string) (invite models.Invite, err response.ErrorResponse)
}

type Identity interface {
	GetID() string
	GetFullName() string
	GetPhone() string
	GetRoles() []string
}

type service struct {
	db                     *gorm.DB
	logger                 log.Logger
	notifier               notification.Notifier
	accessTokenSigningKey  string
	refreshTokenSigningKey string
	tokenExpiration        int
	refreshTokenExpiration int
}

func MakeService(
	db *gorm.DB, logger log.Logger, notifier notification.Notifier,
	accessTokenSigningKey, refreshTokenSigningKey string, accessTokenExpiration, refreshTokenExpiration int,
) Service {

	return &service{
		db,
		logger,
		notifier,
		accessTokenSigningKey,
		refreshTokenSigningKey,
		accessTokenExpiration,
		refreshTokenExpiration,
	}
}
