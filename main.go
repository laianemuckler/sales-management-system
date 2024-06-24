package main

import (
	"fmt"
	"log"
	"net/http"
	"store/database"
	"store/models"
	"strconv"
	"text/template"
	"time"

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

	r.HandleFunc("/login", ShowLoginFormHandler).Methods("GET")
	r.HandleFunc("/login", LoginAuthHandler).Methods("POST")

	r.HandleFunc("/dashboard", DashboardHandler).Methods("GET")

	r.HandleFunc("/backup", DbBackupHandler).Methods("POST")

	r.HandleFunc("/sell", ShowSellFormHandler).Methods("GET")
	r.HandleFunc("/sell", SellHandler).Methods("POST")
	r.HandleFunc("/sell-success", SellSuccessHandler).Methods("GET")

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":3000", r)
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
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		fmt.Println("Usuário ou senha incorretos. Tente novamente.")
		err := tpl.ExecuteTemplate(w, "login.html", "Verifique seu usuário e senha.")
		if err != nil {
			log.Println("Erro ao executar o template:", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		}
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		log.Println("Erro ao executar o template:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
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

func ShowSellFormHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/order.html")
}

func SellHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao analisar os dados do formulário", http.StatusBadRequest)
		return
	}

	orderIDInput := r.FormValue("id")
	totalValueInput := r.FormValue("total_value")
	employeeIDInput := r.FormValue("employee_id")

	orderID, err := strconv.Atoi(orderIDInput)
	if err != nil {
		http.Error(w, "Código da venda inválido", http.StatusBadRequest)
		return
	}

	totalValue, err := strconv.ParseFloat(totalValueInput, 64)
	if err != nil {
		http.Error(w, "Valor total inválido", http.StatusBadRequest)
		return
	}

	employeeID, err := strconv.Atoi(employeeIDInput)
	if err != nil {
		http.Error(w, "Código do funcionário inválido", http.StatusBadRequest)
		return
	}

	order := models.Order{
		ID:         orderID,
		Time:       time.Now(),
		TotalValue: totalValue,
		EmployeeID: employeeID,
	}

	err = database.InsertSale(authenticatedUser.Username, authenticatedUser.Password, order)
	if err != nil {
		http.Error(w, "Erro ao inserir venda no banco de dados", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sell-success", http.StatusSeeOther)
}

func SellSuccessHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "sell-success.html", nil)
	if err != nil {
		log.Println("Erro ao renderizar o template:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}
