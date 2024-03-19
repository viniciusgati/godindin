package database

type Cartao struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	PrimeiroDiaDeCompras int    `json:"primeiroDiaDeCompras"`
	diaDoVencimento      int    `json:"diaDoVencimento"`
}

func IsValidoParaCriar(cartao Cartao) bool {
	if cartao.Name == "" {
		return false
	}
	if len(cartao.Name) > 50 || len(cartao.Name) < 3 {
		return false
	}
	return cartao.PrimeiroDiaDeCompras > 0 && cartao.diaDoVencimento > 0
}
