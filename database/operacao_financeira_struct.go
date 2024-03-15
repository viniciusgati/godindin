package database

type OperacaoFinanceira struct {
	Id        int
	Valor     int
	IdOrigem  int
	Descricao string
	Debito    bool
	Fixo      bool
}
