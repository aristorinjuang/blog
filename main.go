package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/aristorinjuang/blog/controllers"
	"github.com/aristorinjuang/blog/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	controllers.Sessions = make(map[string]*models.User)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	models.Db, models.Err = sql.Open("mysql", os.Getenv("DATABASE"))
	if models.Err != nil {
		log.Fatal(models.Err)
	}

	models.Err = models.Db.Ping()
	if models.Err != nil {
		log.Fatal(models.Err)
	}
}

func main() {
	http.HandleFunc("/single", controllers.Controller)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/dashboard", controllers.Dashboard)
	http.HandleFunc("/", controllers.Index)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8000", nil)
}
