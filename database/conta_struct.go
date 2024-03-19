package database

import "time"

type Contas struct {
	Id       int
	Valor    int
	Debito   bool
	DataHora time.Time
}
