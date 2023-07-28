package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/artificial-lua/example-account-go/dbconnector"
	"golang.org/x/crypto/pbkdf2"
)

type Account struct {
	email    string // must be pk
	password string
	hash     string
	salt     string
	name     string
	birth    date.Date
	gender   string
}

func (a *Account) GetEmail() string {
	return a.email
}

func (a *Account) GetPassword() string {
	return a.password
}

func (a *Account) GetHash() string {
	return a.hash
}

func (a *Account) GetSalt() string {
	return a.salt
}

func (a *Account) GetName() string {
	return a.name
}

func (a *Account) GetBirth() date.Date {
	return a.birth
}

func (a *Account) GetGender() string {
	return a.gender
}

const accountTableColumns = "email, hash, salt, name, birth, gender"

var emailError = errors.New("email is invalid")
var emailDuplicateError = errors.New("email is duplicate")
var passwordError = errors.New("password is invalid")
var nameError = errors.New("name is invalid")
var birthError = errors.New("birth is invalid")
var genderError = errors.New("gender is invalid")

func checkEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}
	return true
}

func checkEmailIsNotExists(email string) bool {
	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("SELECT email FROM account WHERE email = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		return false
	}

	return true
}

func checkPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func checkName(name string) bool {
	if len(name) < 1 {
		return false
	}
	return true
}

func checkBirth(birth date.Date) bool {
	if birth.After(time.Now()) {
		return false
	}
	return true
}

func checkGender(gender string) bool {
	if gender == "M" || gender == "F" {
		return true
	}
	return false
}

func cryptoString(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(dk)
}

func (a *Account) CryptoPassword(salt string) {
	if salt == "" {
		salt = "salt"
	}

	a.hash = cryptoString(a.password, salt)
	a.salt = salt
}

func MakeAccountObject(email string, password string, hash string, salt string, name string, birth string, gender string) (*Account, error) {
	dateBirth, err := date.ParseDate(birth)
	if err != nil {
		return nil, err
	}

	newAccount := &Account{
		email:    email,
		password: password,
		hash:     hash,
		salt:     salt,
		name:     name,
		birth:    dateBirth,
		gender:   gender,
	}

	if newAccount.hash == "" {
		newAccount.CryptoPassword(newAccount.salt)
	}

	return newAccount, nil
}

func checkUpdateAccount(ac *Account) error {
	if !checkPassword(ac.password) {
		return passwordError
	}

	if !checkName(ac.name) {
		return nameError
	}

	if !checkBirth(ac.birth) {
		return birthError
	}

	if !checkGender(ac.gender) {
		return genderError
	}
	return nil
}

func checkCreateAccount(ac *Account) error {
	if !checkEmail(ac.email) {
		return emailError
	}

	if !checkEmailIsNotExists(ac.email) {
		return emailDuplicateError
	}
	return checkUpdateAccount(ac)
}

func CreateAcccount(ac *Account) (sql.Result, error) {
	err := checkCreateAccount(ac)
	if err != nil {
		return nil, err
	}

	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the statement
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO account (%s) VALUES ($1, $2, $3, $4, $5, $6)", accountTableColumns))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Println("Setup the statement")

	// Execute the statement with parameters
	result, err := stmt.Exec(ac.email, ac.hash, ac.salt, ac.name, ac.birth.String(), ac.gender)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func ReadAccountByEmail(email string) (*Account, error) {
	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(fmt.Sprintf("SELECT %s FROM account WHERE email = $1", accountTableColumns))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		log.Fatal(err)
	}

	var hash string
	var salt string
	var name string
	var birth string
	var gender string

	for rows.Next() {
		err := rows.Scan(&email, &hash, &salt, &name, &birth, &gender)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	dateBirth, _ := date.ParseDate(birth[:10])

	return &Account{
		email:  email,
		hash:   hash,
		salt:   salt,
		name:   name,
		birth:  dateBirth,
		gender: gender,
	}, nil
}

func UpdateAccount(ac *Account) (sql.Result, error) {
	err := checkUpdateAccount(ac)
	if err != nil {
		return nil, err
	}

	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the statement
	stmt, err := db.Prepare(`UPDATE account
	SET 
		hash = $2,
		salt = $3,
		name = $4,
		birth = $5,
		gender = $6
	WHERE email = $1`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Println(ac)

	// Execute the statement with parameters
	result, err := stmt.Exec(
		ac.email,
		ac.hash,
		ac.salt,
		ac.name,
		ac.birth.String(),
		ac.gender,
	)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func DeleteAccountByEmail(email string) (sql.Result, error) {
	db, err := dbconnector.NewPostgreSQLConnector()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the statement
	stmt, err := db.Prepare("DELETE FROM account WHERE email = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the statement with parameters
	result, err := stmt.Exec(email)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
