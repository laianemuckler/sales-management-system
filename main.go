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

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")
	database.DBconnection()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/login", ShowLoginFormHandler).Methods("GET")
	r.HandleFunc("/login", LoginAuthHandler).Methods("POST")

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

	usersAndPasswords, err := database.GetUserAndPassword()
	if err != nil {
		log.Println("Erro ao obter usuários e senhas do banco de dados:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	valid := false
	for _, u := range usersAndPasswords {
		if u.Name == username && u.Password == password {
			valid = true
			break
		}
	}

	if valid {
		err := tpl.ExecuteTemplate(w, "dashboard.html", nil)
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