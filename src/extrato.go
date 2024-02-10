package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type TransacaoExtrato struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type ExtratoResposta struct {
	Saldo             Saldo              `json:"saldo"`
	UltimasTransacoes []TransacaoExtrato `json:"ultimas_transacoes"`
}

func ExtratoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	saldo, ultimasTransacoes, err := obterExtrato(id)
	if err != nil {
		if err == sql.ErrNoRows {
			buildError(w, "Extrato não encontrado para o cliente", http.StatusNotFound)
		} else {
			buildError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	resposta := ExtratoResposta{
		Saldo:             saldo,
		UltimasTransacoes: ultimasTransacoes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposta)
}

func obterExtrato(id string) (Saldo, []TransacaoExtrato, error) {
	rows, err := db.Query(`
        SELECT c.saldo, c.limite, t.valor, t.tipo, t.descricao, t.realizada_em
        FROM clientes c
        LEFT JOIN transacoes t ON c.id = t.id_cliente
        WHERE c.id = $1
        ORDER BY t.realizada_em DESC
        LIMIT 10`, id)
	if err != nil {
		return Saldo{}, nil, err
	}
	defer rows.Close()

	var saldo Saldo
	var ultimasTransacoes []TransacaoExtrato
	for rows.Next() {
		var total, limite, valor int
		var tipo, descricao string
		var realizadaEm time.Time
		if err := rows.Scan(&total, &limite, &valor, &tipo, &descricao, &realizadaEm); err != nil {
			return newSaldo(total, limite), nil, nil
		}
		if saldo.Total == 0 {
			saldo = newSaldo(total, limite)
		}
		ultimasTransacoes = append(ultimasTransacoes, TransacaoExtrato{
			Valor:       valor,
			Tipo:        tipo,
			Descricao:   descricao,
			RealizadaEm: realizadaEm,
		})
	}
	if len(ultimasTransacoes) == 0 {
		return Saldo{}, nil, sql.ErrNoRows
	}
	return saldo, ultimasTransacoes, nil
}

func newSaldo(total, limite int) Saldo {
	return Saldo{
		Total:       total,
		DataExtrato: time.Now(),
		Limite:      limite,
	}
}
