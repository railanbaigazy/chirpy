package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		ID:       newID,
		Body:     body,
		AuthorID: userID,
	}

	dbStructure.Chirps[newID] = newChirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps(authorIDStr string) ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	if authorIDStr == "" {
		for _, chirp := range dbStructure.Chirps {
			chirps = append(chirps, chirp)
		}
		return chirps, nil
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		return nil, err
	}

	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
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

func (db *DB) DeleteChirp(chirpID int, userID int) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[chirpID]
	if !ok {
		return errors.New("chirp not found")
	}

	if chirp.AuthorID != userID {
		return errors.New("access denied")
	}
	delete(dbStructure.Chirps, chirpID)

	if err = db.writeDB(dbStructure); err != nil {
		return err
	}
	return nil
}
