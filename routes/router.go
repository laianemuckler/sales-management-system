package routes

import (
	"net/http"
)

func GetEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get a employee"))
}

func GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get employess"))
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get a product"))
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get products"))
}
