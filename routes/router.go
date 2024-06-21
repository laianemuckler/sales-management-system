package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"store/database"
	"store/models"
)

func GetEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get employee"))
}

func GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	var employees []models.Employee
	var err error
	employees, err = database.GetEmployees()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: %s", err)
		return
	}
	json.NewEncoder(w).Encode(&employees)
}

