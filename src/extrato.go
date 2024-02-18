package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type StatementBalance struct {
	Total int       `json:"total"`
	Date  time.Time `json:"data_extrato"`
	Limit int       `json:"limite"`
}

type StatementTransaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type StatementResponse struct {
	Balance      StatementBalance       `json:"saldo"`
	Transactions []StatementTransaction `json:"ultimas_transacoes"`
}

func StatementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	balance, transactions, err := getStatement(id)
	if err != nil {
		if err == sql.ErrNoRows {
			buildError(w, "Extrato não encontrado para o cliente", http.StatusNotFound)
		} else {
			buildError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := StatementResponse{
		Balance:      balance,
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getStatement(id string) (StatementBalance, []StatementTransaction, error) {
	rows, err := db.Query(`
        SELECT c.saldo, c.limite, t.valor, t.tipo, t.descricao, t.realizada_em
        FROM clientes c
        LEFT JOIN transacoes t ON c.id = t.id_cliente
        WHERE c.id = $1
        ORDER BY t.realizada_em DESC
        LIMIT 10`, id)
	if err != nil {
		return StatementBalance{}, nil, err
	}
	defer rows.Close()

	var balance StatementBalance
	var transactions []StatementTransaction
	for rows.Next() {
		var total, limit int
		var transaction StatementTransaction
		if err := rows.Scan(&total, &limit, &transaction.Value, &transaction.Type, &transaction.Description, &transaction.CreatedAt); err != nil {
			return newBalance(total, limit), nil, nil
		}
		if balance.Total == 0 {
			balance = newBalance(total, limit)
		}
		transactions = append(transactions, transaction)
	}
	if len(transactions) == 0 {
		return StatementBalance{}, nil, sql.ErrNoRows
	}
	return balance, transactions, nil
}

func newBalance(total, limit int) StatementBalance {
	return StatementBalance{
		Total: total,
		Date:  time.Now(),
		Limit: limit,
	}
}
