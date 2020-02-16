package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"text/template"
)

// Checks if the user is logged in and has a session, if not err is not nil
func session(w http.ResponseWriter, r *http.Request) (sess Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

func generateHTML(w http.ResponseWriter, data interface{}, layout string, filenames ...string) {
	var files []string
	for _, file := range filenames {
		// files = append(files, fmt.Sprintf("athena/templates/%s.html", file))
		files = append(files, fmt.Sprintf("template/%s.html", file))

	}

	templates := template.Must(template.ParseFiles(files...))

	templates.ExecuteTemplate(w, layout, data)
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
