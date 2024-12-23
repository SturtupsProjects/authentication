package token

import (
	"authentification/config"
	"github.com/golang-jwt/jwt"
	"strconv"
)

type Claims struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	PhoneNumber string `json:"phone_number"`
	CompanyId   string `json:"company_id"`
	Role        string `json:"role"`
	jwt.StandardClaims
}

func ConfigToken(config *config.Config) error {

	exAcc, err := strconv.Atoi(config.EXPIRED_ACCESS)
	if err != nil {
		return err
	}

	refAcc, err := strconv.Atoi(config.EXPIRED_REFRESH)
	if err != nil {
		return err
	}

	AccessSecretKey = config.ACCESS_TOKEN
	RefreshSecretKey = config.ACCESS_TOKEN
	ExpiredAccess = exAcc
	ExpiredRefresh = refAcc

	return nil
}
