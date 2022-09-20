package utils

import (
	"database/sql"
	"time"
)

//authSession stores http session information.
type authSession struct {
	Key       string
	Data      []byte
	CreatedOn time.Time
	UpdatedOn time.Time
	ExpiresOn time.Time
}

func (s *authSession) Create(db *sql.DB) error {
	var query = `
	BEGIN TRANSACTION;
	  INSERT INTO sessions  (key, data, created_on, updated_on, expires_on)
		VALUES ($1,$2,now(),now(),$3);
	COMMIT;
	`
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, s.Key, s.Data, s.ExpiresOn)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *authSession) FindByKey(db *sql.DB) error {
	var query = `
	SELECT * from sessions  WHERE key=$1 LIMIT 1;
	`
	return db.QueryRow(query, s.Key).Scan(
		&s.Key,
		&s.Data,
		&s.CreatedOn,
		&s.UpdatedOn,
		&s.ExpiresOn,
	)
}

func (s *authSession) Update(db *sql.DB) error {
	var query = `
BEGIN TRANSACTION;
  UPDATE sessions
    data = $2,
    updated_on = now(),
  WHERE key=$1;
COMMIT;
	`

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, s.Key, s.Data)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *authSession) Delete(db *sql.DB) error {
	var query = `
BEGIN TRANSACTION;
   DELETE FROM sessions
  WHERE key=$1;
COMMIT;
`
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, s.Key)
	if err != nil {
		return err
	}
	return tx.Commit()
}
