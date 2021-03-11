package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aristorinjuang/blog/helpers"
	"github.com/aristorinjuang/blog/models"
	"github.com/aristorinjuang/blog/views"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
)

var data struct {
	Article  *models.Article
	Articles []*models.Article
}

// Sessions is a user sessions.
var Sessions map[string]*models.User

func getCookie(r *http.Request) *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: "",
	}

	for _, c := range r.Cookies() {
		if c.Name == "session" {
			cookie.Value = c.Value
			break
		}
	}

	return cookie
}

func permission(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")[1]
	cookie := getCookie(r)

	switch path {
	case "dashboard":
		if Sessions[cookie.Value] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	case "login":
	case "register":
		if Sessions[cookie.Value] != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}

// Index is a homepage controller.
func Index(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")[1]

	if path == "" {
		data.Articles = models.Articles()

		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "index.gohtml", data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
	} else {
		data.Article = models.FindArticle(path)

		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "single.gohtml", data); err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
	}
}

// Register is a controller to register a user.
func Register(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "register.gohtml", nil); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
	case http.MethodPost:
		user := &models.User{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := helpers.ValidateEmail(user.Email); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user.Password = password
		user = user.Find()

		if user.ID == 0 {
			user = user.Create()
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Login is a controller for users to log in.
func Login(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "login.gohtml", nil); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
	case http.MethodPost:
		user := &models.User{
			Email: r.FormValue("email"),
		}
		user = user.Find()

		if user.ID == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password")))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Logout is a controller for users to log out.
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	if Sessions[cookie.Value] != nil {
		delete(Sessions, cookie.Value)
	}

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Dashboard is a controller for users to list articles.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("model")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "article":
			if delete != "" && id != 0 {
				article := &models.Article{
					ID: id,
				}

				user.DeleteArticle(article)

				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			} else {
				data.Article = &models.Article{
					Image:   "",
					Title:   "",
					Content: "",
				}

				if user != nil && id != 0 {
					data.Article = user.FindArticle(id)
				}

				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "article.gohtml", data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				io.WriteString(w, response.String())
			}
		default:
			if user != nil {
				data.Articles = user.FindArticles()
			}

			var response bytes.Buffer
			if err := views.Tmpl.ExecuteTemplate(&response, "dashboard.gohtml", data); err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			io.WriteString(w, response.String())
		}
	case http.MethodPost:
		switch model {
		case "article":
			if user != nil {
				if id == 0 {
					article := &models.Article{
						Image:     r.FormValue("image"),
						Slug:      slug.Make(r.FormValue("title")),
						Title:     r.FormValue("title"),
						Content:   r.FormValue("content"),
						Author:    *user,
						CreatedAt: time.Now(),
					}

					user.CreateArticle(article)
				} else {
					article := &models.Article{
						ID:      id,
						Image:   r.FormValue("image"),
						Slug:    slug.Make(r.FormValue("title")),
						Title:   r.FormValue("title"),
						Content: r.FormValue("content"),
					}

					user.UpdateArticle(article)
				}
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}
