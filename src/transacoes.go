package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Transaction struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type TransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		buildError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	// Validar requisição
	if transaction.Type != "c" && transaction.Type != "d" {
		buildError(w, "Tipo de transação inválido", http.StatusUnprocessableEntity)
		return
	}

	if transaction.Value <= 0 || len(transaction.Description) < 1 || len(transaction.Description) > 10 {
		buildError(w, "Dados da transação inválidos", http.StatusUnprocessableEntity)
		return
	}

	// Salvar transação e recuperar saldo atualizado
	balance, limit, err := saveTransaction(id, transaction)
	if err != nil {
		if err.(*pq.Error).Message == "CLIENTE_NAO_ENCONTRADO" {
			buildError(w, "Cliente não encontrado", http.StatusNotFound)
		} else if err.(*pq.Error).Message == "LIMITE_EXECEDIDO" {
			buildError(w, "Limite excedido", http.StatusUnprocessableEntity)
		} else {
			buildError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := TransactionResponse{
		Limit:   limit,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func saveTransaction(id string, transaction Transaction) (int, int, error) {
	var balance, limit int
	err := db.QueryRow("SELECT * FROM atualizar_saldo($1, $2, $3, $4)", id, transaction.Value, transaction.Type, transaction.Description).Scan(&balance, &limit)
	if err != nil {
		return 0, 0, err
	}
	return balance, limit, nil
}
