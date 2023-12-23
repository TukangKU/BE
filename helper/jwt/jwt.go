package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(idUser uint, role string) (string, error) {
	var claim = jwt.MapClaims{}
	claim["id"] = idUser
	claim["role"] = role
	claim["iat"] = time.Now().UnixMilli()
	claim["exp"] = time.Now().Add(time.Millisecond * 1).UnixMilli()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// token := jwt.New(jwt.SigningMethodHS256)
	strToken, err := token.SignedString([]byte("$!1gnK3yyy!!!"))
	if err != nil {
		return "", err
	}

	return strToken, nil
}

func ExtractToken(t *jwt.Token) (uint, error) {
	var userID uint

	expiredTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}

	var eTime = *expiredTime

	if t.Valid && eTime.Compare(time.Now()) > 0 {
		var tokenClaims = t.Claims.(jwt.MapClaims)
		userID = uint(tokenClaims["id"].(float64))

		return userID, nil
	}

	return 0, errors.New("token tidak valid")
}

// func GenerateJWT(idUser uint, role string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"id":   idUser,
// 		"iat":  time.Now().Unix(),
// 		"exp":  time.Now().Add(time.Hour).Unix(), // Set expiration to 1 hour
// 		"role": role,
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	strToken, err := token.SignedString([]byte("$!1gnK3yyy!!!"))
// 	if err != nil {
// 		return "", err
// 	}

//		return strToken, nil
//	}
func ExtractTokenRole(t *jwt.Token) (string, error) {

	expiredTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return "", errors.New("harap login lagi")
	}

	var eTime = *expiredTime

	if t.Valid && eTime.Compare(time.Now()) > 0 {
		var tokenClaims = t.Claims.(jwt.MapClaims)
		role := tokenClaims["role"]
		str := fmt.Sprintf("%v", role)

		return str, nil
	}

	return "", errors.New("token tidak valid")
}
