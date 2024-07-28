package database

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 int       `json:"id"`
	Email              string    `json:"email"`
	Password           []byte    `json:"password"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

type UserResp struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type LoginResp struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResp struct {
	Token string `json:"token"`
}

func (db *DB) CreateUser(email string, password []byte) (UserResp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResp{}, err
	}

	if ok, _ := dbStructure.userExists(email); ok {
		return UserResp{}, errors.New("user already exists")
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    strings.ToLower(email),
		Password: password,
	}
	userResp := UserResp{
		ID:    id,
		Email: strings.ToLower(email),
	}

	dbStructure.Users[id] = user
	if err = db.writeDB(dbStructure); err != nil {
		return UserResp{}, err
	}

	return userResp, nil
}

func (db *DB) Login(email string, password string, secretKey []byte) (LoginResp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return LoginResp{}, err
	}

	ok, user := dbStructure.userExists(email)
	if !ok {
		return LoginResp{}, errors.New("no such user found")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return LoginResp{}, errors.New("incorrect password")
	}

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return LoginResp{}, errors.New("error signing the token")
	}

	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return LoginResp{}, errors.New("error generating random bytes")
	}
	refreshToken := hex.EncodeToString(bytes)
	refreshTokenExpiry := now.Add(60 * 24 * time.Hour)

	user.RefreshToken = refreshToken
	user.RefreshTokenExpiry = refreshTokenExpiry
	dbStructure.Users[user.ID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return LoginResp{}, err
	}

	loginResp := LoginResp{
		ID:           user.ID,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
	}

	return loginResp, nil
}

func (dbStructure *DBStructure) userExists(email string) (bool, User) {
	for _, userVal := range dbStructure.Users {
		if email == userVal.Email {
			return true, userVal
		}
	}
	return false, User{}
}

func (db *DB) UpdateUser(id int, newEmail string, newPassword []byte) (UserResp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResp{}, err
	}

	userBody, exists := dbStructure.Users[id]
	if !exists {
		return UserResp{}, errors.New("user not found")
	}

	user := User{
		ID:                 id,
		Email:              strings.ToLower(newEmail),
		Password:           newPassword,
		RefreshToken:       userBody.RefreshToken,
		RefreshTokenExpiry: userBody.RefreshTokenExpiry,
	}

	dbStructure.Users[id] = user
	if err := db.writeDB(dbStructure); err != nil {
		return UserResp{}, err
	}

	userResp := UserResp{
		ID:    user.ID,
		Email: user.Email,
	}

	return userResp, nil
}

func (db *DB) RefreshAccessToken(refreshToken string, secretKey []byte) (RefreshResp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshResp{}, err
	}

	userID := 0
	exists := false
	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			if time.Now().After(user.RefreshTokenExpiry) {
				return RefreshResp{}, errors.New("expired refresh token")
			} else {
				userID = user.ID
				exists = true
				break
			}
		}
	}
	if !exists {
		return RefreshResp{}, errors.New("refresh token not found")
	}

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		Subject:   strconv.Itoa(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return RefreshResp{}, errors.New("error signing the token")
	}

	return RefreshResp{Token: tokenString}, nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	userBody := User{}
	exists := false
	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			if time.Now().After(user.RefreshTokenExpiry) {
				return errors.New("expired refresh token")
			} else {
				userBody = user
				exists = true
				break
			}
		}
	}
	if !exists {
		return errors.New("refresh token not found")
	}

	userBody.RefreshToken = ""
	userBody.RefreshTokenExpiry = time.Time{}
	dbStructure.Users[userBody.ID] = userBody

	if err = db.writeDB(dbStructure); err != nil {
		return err
	}
	return nil
}
