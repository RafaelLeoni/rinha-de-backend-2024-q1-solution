package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	_ "net/http/pprof"
)

type Erro struct {
	Message   string    `json:"erro"`
	Timestamp time.Time `json:"ocorreu_em"`
}

var (
	db *sql.DB
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("A variável de ambiente DB_URL não está definida")
	}

	err := initDB(dbURL)
	if err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/clientes/{id}/transacoes", TransacaoHandler).Methods("POST")
	r.HandleFunc("/clientes/{id}/extrato", ExtratoHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), nil))
}

func initDB(dbURL string) error {
	var err error
	db, err = sql.Open("postgres", dbURL)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso")
	return nil
}

func buildError(w http.ResponseWriter, msg string, stausCode int) {
	log.Println(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stausCode)
	json.NewEncoder(w).Encode(&Erro{
		Message:   msg,
		Timestamp: time.Now(),
	})
}
