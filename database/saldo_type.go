package database

import "time"

type Saldo struct {
	Id       string    `json:"id"`
	DataHora time.Time `json:"data_hora"`
	Valor    int       `json:"valor"`
}
