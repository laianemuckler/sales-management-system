package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"store/models"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func DBconnection(username, password string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=localhost port=5433 user=%s password=%s dbname=loja sslmode=disable", username, password)

	db, err = sql.Open("postgres", connectionString)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database Connected!")
	}

	return db, err
}

func GetEmployees(username, password string) ([]models.Employee, error) {
	connection, err := DBconnection(username, password)
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

func GetUserAndPassword(username, password string) ([]models.Employee, error) {
	connection, err := DBconnection(username, password)
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

func MakeDbBackup(username, password string) error {
	if username != "adminvendas" {
		return fmt.Errorf("usuário não tem permissão para fazer backup")
	}

	backupDir := "C:\\Users\\Free\\Desktop" // Exemplo: Caminho para a área de trabalho do Windows

	backupFile := filepath.Join(backupDir, "backup.sql")

	cmd := exec.Command("pg_dump", "-U", username, "-d", "loja")

	outFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de backup: %v", err)
	}
	defer outFile.Close()
	cmd.Stdout = outFile

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("erro ao executar pg_dump: %v", err)
	}

	fmt.Println("Backup do banco de dados realizado com sucesso:", backupFile)

	return nil
}
