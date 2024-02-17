package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Transacao struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type Resposta struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

func TransacaoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var transacao Transacao
	if err := json.NewDecoder(r.Body).Decode(&transacao); err != nil {
		buildError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	// Validar requisição
	if transacao.Tipo != "c" && transacao.Tipo != "d" {
		buildError(w, "Tipo de transação inválido", http.StatusUnprocessableEntity)
		return
	}

	if transacao.Valor <= 0 || len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 {
		buildError(w, "Dados da transação inválidos", http.StatusUnprocessableEntity)
		return
	}

	// Atualizar o saldo
	saldo, limite, err := incluirTransacao(id, transacao)
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

	resposta := Resposta{
		Limite: limite,
		Saldo:  saldo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposta)
}

func incluirTransacao(id string, transacao Transacao) (int, int, error) {
	var saldo, limite int
	err := db.QueryRow("SELECT * FROM atualizar_saldo($1, $2, $3, $4)", id, transacao.Valor, transacao.Tipo, transacao.Descricao).Scan(&saldo, &limite)
	if err != nil {
		return 0, 0, err
	}
	return saldo, limite, nil
}
