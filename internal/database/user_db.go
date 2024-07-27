package database

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserResp struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type LoginResp struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

func (db *DB) Login(email string, password string, secretKey []byte, expiresInSeconds int) (LoginResp, error) {
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
	if expiresInSeconds <= 0 || expiresInSeconds > 24*3600 {
		expiresInSeconds = 24 * 3600
	}
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return LoginResp{}, errors.New("error signing the token")
	}
	loginResp := LoginResp{
		ID:    user.ID,
		Email: user.Email,
		Token: tokenString,
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

	_, exists := dbStructure.Users[id]
	if !exists {
		return UserResp{}, errors.New("user not found")
	}

	user := User{
		ID:       id,
		Email:    strings.ToLower(newEmail),
		Password: newPassword,
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
