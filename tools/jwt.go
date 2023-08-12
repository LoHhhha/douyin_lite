package tools

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// Claims
// use userId and userName to form token
type Claims struct {
	Id   int64
	Name string
	jwt.StandardClaims
}

var issuer = "DouYin"
var key = "lohhhha"

// GetToken
// @param: id(int64), name(string)
// @return: token(string), err(error)
func GetToken(id int64, name string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour)
	claims := Claims{
		Id:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken
// @param: token(string)
// @return: claims(*Claims), err(error)
// When token invalid return nil, err.
// such as token was out-of-date
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// check str`s SigningMethod
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// CheckToken
// @param: token(string)
// @return: valid(bool)
// When token parse not err, we say the token is valid.
func CheckToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}
