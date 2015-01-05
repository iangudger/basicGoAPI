/*
 * Functions for querying the database. Each function should represent the most basic useful database operation.
 *
 * The functions in the package should only be accessed through the common package.
 */
package database

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	alphaNumeric            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sessionIdLength         = 25
	generatedPasswordLength = 15
)

// Errors
var UsernameInUse = errors.New("You have already registered with this email address.")

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(b)
}

func newSessionID() string {
	return randSeq(sessionIdLength)
}

func genPassword() string {
	return randSeq(generatedPasswordLength)
}

type DB struct {
	conn *sql.DB
}

func NewDB(constr string) (DB, error) {
	conn, err := connect(constr)
	if err != nil {
		return DB{}, err
	}
	return DB{conn}, nil
}

// Establishes connection to Postgres database
func connect(constr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", constr)
	return conn, err
}

func (db DB) GetUserID(email string, password string) (userid int, err error) {
	log.Printf("Looking up user: %s\n", email)

	// Get hashed password
	var hashed string
	err = db.conn.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&userid, &hashed)
	if err != nil {
		return
	}

	// Check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return
}

func (db DB) GetSessionUserID(sessionid string) (userid int, err error) {
	log.Printf("Looking up userid for sessionid: %s\n", sessionid)

	err = db.conn.QueryRow("SELECT userid FROM sessions WHERE sessions.id = $1", sessionid).Scan(&userid)
	return
}

func (db DB) NewSession(userid int) (sessionid string, err error) {
	// Create a new unique sessionid
	for numRows := 0; ; {
		sessionid = newSessionID()
		err = db.conn.QueryRow("SELECT count(*) FROM sessions WHERE sessions.id = $1", sessionid).Scan(&numRows)
		if err != nil {
			return
		}
		if numRows == 0 {
			break
		}
	}

	// Insert new sessionid
	_, err = db.conn.Exec(
		"INSERT INTO sessions (userid, id) VALUES ($1, $2)",
		userid,
		sessionid,
	)

	return
}

func (db DB) RegisterUser(email string) (password string, err error) {
	// Check if email is already in use
	var numRows int
	err = db.conn.QueryRow("SELECT count(*) FROM users WHERE email = $1", email).Scan(&numRows)
	if err != nil {
		return
	}
	if numRows != 0 {
		err = UsernameInUse
		return
	}

	// Add new user
	password = genPassword()
	hashed, err := hashPassword(password)
	if err != nil {
		return
	}
	log.Printf(
		"Adding user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	_, err = db.conn.Exec(
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		email,
		hashed,
	)
	return
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (db DB) GetEmailFromSessionID(sessionid string) (email string, err error) {
	err = db.conn.QueryRow(
		"SELECT email FROM sessions, users WHERE sessions.id = $1 AND sessions.userid = users.id",
		sessionid,
	).Scan(&email)
	return
}

func (db DB) GetEmail(sessionid string) (email string, err error) {
	err = db.conn.QueryRow(
		"SELECT email FROM sessions, users WHERE sessions.id = $1 AND sessions.userid = users.id",
		sessionid,
	).Scan(&email)
	return
}

func (db DB) ResetPassword(email string) (password string, err error) {
	password = genPassword()
	hashed, err := hashPassword(password)
	if err != nil {
		return
	}
	log.Printf(
		"Resetting password for user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	err = db.changePassword(email, hashed)
	return
}

func (db DB) changePassword(email, hashed string) error {
	return db.conn.QueryRow(
		"UPDATE users SET password = $2 WHERE email = $1 RETURNING password;",
		email,
		hashed,
	).Scan(&hashed)
}

func (db DB) ChangePassword(email, newPassword string) (err error) {
	hashed, err := hashPassword(newPassword)
	if err != nil {
		return
	}
	log.Printf(
		"Changing password for user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	err = db.changePassword(email, hashed)
	return
}

func (db DB) UpdateLocation(userid int, lon, lat float64) error {
	_, err := db.conn.Exec(
		"INSERT INTO locations (userid, location) VALUES ($1, point($2, $3));",
		userid,
		lon,
		lat,
	)
	return err
}

func (db DB) AddListing(userid int, title, description, price string, lon, lat float64) error {
	_, err := db.conn.Exec(
		"INSERT INTO listings (userid, title, description, price, location) VALUES ($1, $2, $3, $4, point($5, $6));",
		userid,
		title,
		description,
		price,
		lon,
		lat,
	)
	return err
}

func (db DB) Logout(sessionid string) error {
	return db.conn.QueryRow(
		"DELETE FROM sessions WHERE sessions.id = $1 RETURNING sessions.id;",
		sessionid,
	).Scan(&sessionid)
}

func (db DB) UpdateSession(sessionid string) error {
	log.Println("Updating last date for sessionid", sessionid)
	return db.conn.QueryRow(
		"UPDATE sessions SET last = now() WHERE sessions.id = $1 RETURNING sessions.id;",
		sessionid,
	).Scan(&sessionid)
}
