package jwttoken

import (
	"context"
	"errors"
	"main/database"
	"main/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	db             = database.DBinstance()
	userCollection = database.OpenCollection(db, "users")
	key            = "SECRET_KEY"
)

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserType  string
	UserID    string
	jwt.RegisteredClaims
}

func GenerateAllTokens(user model.User) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		UserType:  *user.UserType,
		UserID:    user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(168))),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		return
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(key))
	if err != nil {
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userID string) error {
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObject := make(map[string]interface{})
	updateObject["token"] = signedToken
	updateObject["refresh_token"] = signedRefreshToken
	updateObject["updated_at"] = updatedAt

	_, err := userCollection.UpdateDocument(context.TODO(), userID, updateObject)

	return err
}

func ValidateToken(clientToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(clientToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		err = errors.New("failed to parse token")
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err = errors.New("the token is invalid")
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		err = errors.New("the token is expired")
		return
	}

	return claims, nil
}
