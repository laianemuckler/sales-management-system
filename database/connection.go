package database

import (
	"database/sql"
	"fmt"
	"log"
	"store/models"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func DBconnection() (*sql.DB, error) {
	connectionString := "host=localhost port=5433 user= password= dbname=loja sslmode=disable"

	db, err = sql.Open("postgres", connectionString)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database Connected!")
	}

	return db, err
}

func GetEmployees() ([]models.Employee, error) {
	connection, err := DBconnection()
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	rows, err := db.Query("SELECT * FROM funcionarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var employees []models.Employee

	for rows.Next() {
		var e models.Employee
		err = rows.Scan(&e.ID, &e.Name, &e.CPF, &e.Password, e.Occupation)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func GetUserAndPassword() ([]models.Employee, error) {
		connection, err := DBconnection()
		if err != nil {
			return nil, err
		}
		defer connection.Close()
	
		rows, err := connection.Query("SELECT fun_nome, fun_senha FROM funcionarios")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
	
		var usersAndPasswords []models.Employee
	
		for rows.Next() {
			var u models.Employee
			err := rows.Scan(&u.Name, &u.Password)
			if err != nil {
				return nil, err
			}
	
			usersAndPasswords = append(usersAndPasswords, u)
		}
	
		if err = rows.Err(); err != nil {
			return nil, err
		}
	
		return usersAndPasswords, nil
}
