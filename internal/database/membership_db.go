package database

import "errors"

func (db *DB) UpgradeUser(userID int) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[userID]
	if !ok {
		return errors.New("user doesn't exist")
	}

	user.IsChirpyRed = true
	dbStructure.Users[userID] = user
	if err = db.writeDB(dbStructure); err != nil {
		return err
	}

	return nil
}
