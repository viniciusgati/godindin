package database

import "time"

type Resumo struct {
	Id       int
	Valor    int
	Debito   bool
	DataHora time.Time
}
