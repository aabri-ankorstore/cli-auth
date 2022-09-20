package utils

import (
	"database/sql"
	"encoding/base32"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
	"time"
)

const migrationQL = `
BEGIN TRANSACTION;
  CREATE TABLE IF NOT EXISTS sessions(
    key string,
    data blob,
    created_on time,
    updated_on time,
    expires_on time);
COMMIT;
`

// AuthStore is the session storage implementation for gorilla/sessions using
// embedded SQL database(ql).
type AuthStore struct {
	store   *sql.DB
	Codecs  []securecookie.Codec
	Options *sessions.Options
}

// NewAuthStore initializes QLStore with the given keyPairs
func NewAuthStore(store *sql.DB, path string, maxAge int, keyPairs ...[]byte) *AuthStore {
	return &AuthStore{
		store:  store,
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   path,
			MaxAge: maxAge,
		},
	}
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (db *AuthStore) MaxAge(age int) {
	db.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range db.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// Get fetches a session for a given name after it has been added to the registry.
func (db *AuthStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(db, name)
}

// New returns a new session
func (db *AuthStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(db, name)
	opts := *db.Options
	session.Options = &(opts)
	session.IsNew = true

	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, db.Codecs...)
		if err == nil {
			err = db.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	db.MaxAge(db.Options.MaxAge)
	return session, err
}

// Save saves the session into a ql database
func (db *AuthStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Set delete if max-age is < 0
	if session.Options.MaxAge < 0 {
		if err := db.Delete(r, w, session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Generate a random session ID key suitable for storage in the DB
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}

	if err := db.save(session); err != nil {
		return err
	}

	// Keep the session ID key in a cookie so it can be looked up in DB later.
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, db.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

//load fetches a session by ID from the database and decodes its content into session.Values
func (db *AuthStore) load(session *sessions.Session) error {
	s := authSession{Key: session.ID}
	err := s.FindByKey(db.store)
	if err != nil {
		return err
	}
	return securecookie.DecodeMulti(session.Name(), string(s.Data),
		&session.Values, db.Codecs...)
}

func (db *AuthStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		db.Codecs...)
	if err != nil {
		return err
	}
	var expiresOn time.Time
	exOn := session.Values["expires_on"]
	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
		if expiresOn.Sub(time.Now().Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
			expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
		}
	}
	s := authSession{
		Key:       session.ID,
		Data:      []byte(encoded),
		ExpiresOn: expiresOn,
	}
	if session.IsNew {
		return s.Create(db.store)
	}
	return s.Update(db.store)
}

func (db *AuthStore) destroy(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	options := *db.Options
	options.MaxAge = -1
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", &options))
	for k := range session.Values {
		delete(session.Values, k)
	}
	s := authSession{Key: session.ID}
	return s.Delete(db.store)
}

// Delete deletes session.
func (db *AuthStore) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return db.destroy(r, w, session)
}

//Migrate creates the session table if the table does not exist yet.
func Migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(migrationQL)
	if err != nil {
		return err
	}
	return tx.Commit()

}
