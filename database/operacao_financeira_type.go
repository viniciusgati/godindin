package database

type OperacaoFinanceira struct {
	Id        int    `json:"id"`
	Valor     int    `json:"valor"`
	IdOrigem  int    `json:"id_origem"`
	Descricao string `json:"descricao"`
	Debito    bool   `json:"debito"`
	Fixo      bool   `json:"fixo"`
}
