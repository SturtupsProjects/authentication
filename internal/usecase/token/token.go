package token

import (
	"authentification/internal/entity"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

var (
	AccessSecretKey  string
	RefreshSecretKey string
	ExpiredAccess    int
	ExpiredRefresh   int
)

func GenerateAccessToken(in *entity.LogInToken) (string, error) {
	claims := Claims{
		Id:          in.UserId,
		FirstName:   in.FirstName,
		PhoneNumber: in.PhoneNumber,
		CompanyId:   in.CompanyId,
		Role:        in.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpiredAccess)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv(AccessSecretKey)))
}

func GenerateRefreshToken(in *entity.LogInToken) (string, error) {
	claims := Claims{
		Id:          in.UserId,
		FirstName:   in.FirstName,
		PhoneNumber: in.PhoneNumber,
		CompanyId:   in.CompanyId,
		Role:        in.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpiredRefresh)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv(RefreshSecretKey)))
}

func ExtractAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(AccessSecretKey)), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func ExtractRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(RefreshSecretKey)), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func GetExpires() int {
	return ExpiredAccess
}
