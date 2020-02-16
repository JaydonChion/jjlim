package main

import (
	"database/sql"
	"time"
)

type User struct {
	Id                  int
	Uuid                string
	Name                string
	Email               string
	Password            sql.NullString
	CreatedAt           time.Time
	Phone               int
	Gender              sql.NullString
	Dob                 time.Time
	Mnum                string
	Verification_status bool
	CupSaved            int
}

type Session struct {
	Id        int
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid,user_id, created_at) values ($1, $2, $3) returning id, uuid,user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.UserId, &session.CreatedAt)
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.UserId, &session.CreatedAt)
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid,user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// Get the user from the session
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// Delete all sessions from database
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {

	statement := "insert into users (uuid, name, email, password, created_at, phone, gender, dob,mnum) values ($1, $2, $3, $4, $5, $6, $7, $8,$9) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	crypText, err := Encrypt(user.Password.String)
	if err != nil {
		panic(err)

		return
	}

	err = stmt.QueryRow(createUUID(), user.Name, user.Email, crypText, time.Now(), user.Phone, user.Gender, user.Dob, user.Mnum).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// Update user information in the database
// func (user *User) Update() (err error) {
// 	statement := "update users set phone = $2, distance = $3 where id = $1"
// 	stmt, err := Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(user.Id, user.Phone, user.Distance)

// 	statement = "update sessions set distance = $2 where user_id = $1"
// 	stmt, err = Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(user.Id, user.Distance)

// 	return
// }

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT name, email, phone created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Name, &user.Email, &user.Phone); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at,gender, verification_status FROM users WHERE email = $1", email).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.Gender, &user.Verification_status)
	// fmt.Println(err)
	return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	// user = User{}
	// err = Db.QueryRow("SELECT email FROM sessions WHERE uuid = $1", uuid).Scan(&user.Email)
	// err = Db.QueryRow("SELECT id, uuid, name,phone, gender, distance  FROM users WHERE email = $1", user.Email).
	// 	Scan(&user.Id, &user.Uuid, &user.Name, &user.Phone, &user.Gender, &user.Distance)
	err = Db.QueryRow("SELECT id,email, uuid, name,phone, gender FROM users WHERE email = (SELECT email FROM sessions WHERE uuid = $1);", uuid).
		Scan(&user.Id, &user.Email, &user.Uuid, &user.Name, &user.Phone, &user.Gender)
	return
}

func UserByID(id int) (user User, err error) {

	err = Db.QueryRow("SELECT id,email, uuid, name,phone, gender FROM users WHERE id = $1;", id).
		Scan(&user.Id, &user.Email, &user.Uuid, &user.Name, &user.Phone, &user.Gender)
	return

}

func setDoneIntro(userId int) (err error) {
	statement1 := "update users set fb_status=true where id = $1"
	stmt1, err := Db.Prepare(statement1)

	defer stmt1.Close()

	_, err = stmt1.Exec(userId)

	return
}
