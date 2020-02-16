package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var Db *sql.DB

var (
	// host     = os.Getenv("DB_HOST")
	// port     = 5432
	// user     = os.Getenv("DB_USERNAME")
	// password = os.Getenv("DB_PASS")
	// dbname   = os.Getenv("DB_NAME")
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "chionjetherng"
	dbname   = "lootgather"
)

func init() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string, err error) {
	// cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	cryptbytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 14)
	cryptext = string(cryptbytes)
	return
}
