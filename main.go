package main

import (
	"fmt"
	"log"
	"net/http"
	"store/database"
	"store/routes"
	"text/template"

	"github.com/gorilla/mux"
)

var tpl *template.Template

type User struct {
	Username string
	Password string
}

var authenticatedUser User

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", ShowLoginFormHandler).Methods("GET")
	r.HandleFunc("/login", LoginAuthHandler).Methods("POST")

	r.HandleFunc("/backup", DbBackupHandler).Methods("POST")

	r.HandleFunc("/employees/{id}", routes.GetEmployeeByIDHandler).Methods("GET")
	r.HandleFunc("/employees", routes.GetEmployeesHandler).Methods("GET")

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":3000", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****Hello World!*****")
}

func ShowLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Println("Erro ao executar o template:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***** LoginAuthHandler running *****")
	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("username:", username, "password:", password)

	usersAndPasswords, err := database.GetUserAndPassword(username, password)
	if err != nil {
		log.Println("Erro ao obter usuários e senhas do banco de dados:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	valid := false
	for _, u := range usersAndPasswords {
		if u.Name == username && u.Password == password {
			valid = true
			authenticatedUser = User{Username: username, Password: password}
			break
		}
	}

	if valid {
		_, err := database.DBconnection(username, password)
		if err != nil {
			log.Println("Erro ao conectar ao banco de dados:", err)
			http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
			return
		}

		err = tpl.ExecuteTemplate(w, "dashboard.html", nil)
		if err != nil {
			log.Println("Erro ao executar o template:", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		}
	} else {
		fmt.Println("Usuário ou senha incorretos. Tente novamente.")
		err := tpl.ExecuteTemplate(w, "login.html", "Verifique seu usuário e senha.")
		if err != nil {
			log.Println("Erro ao executar o template:", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		}
	}
}

func DbBackupHandler(w http.ResponseWriter, r *http.Request) {
	if authenticatedUser.Username == "" {
		log.Println("Erro: usuário não autenticado.")
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	err := database.MakeDbBackup(authenticatedUser.Username, authenticatedUser.Password)
	if err != nil {
		log.Println("Erro ao realizar backup do banco de dados:", err)
		http.Error(w, "Erro ao realizar backup do banco de dados", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "dashboard.html", "Backup do banco de dados realizado com sucesso!")
	if err != nil {
		log.Println("Erro ao renderizar o template:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}
