package database

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		initialData := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
		return db.writeDB(initialData)
	}
	return nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	file, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, file, 0644)
}

func (db *DB) loadDB() (DBStructure, error) {
	dbStructure := DBStructure{}
	file, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}

	err = json.Unmarshal(file, &dbStructure)
	return dbStructure, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		ID:   newID,
		Body: body,
	}

	dbStructure.Chirps[newID] = newChirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, Chirp{ID: chirp.ID, Body: chirp.Body})
	}
	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("id not found")
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, userVal := range dbStructure.Users {
		if email == userVal.Email {
			return User{}, errors.New("user already exists")
		}
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:    id,
		Email: strings.ToLower(email),
	}

	dbStructure.Users[id] = user
	if err = db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}
