package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cartao struct {
	Id                   int
	Name                 string
	PrimeiroDiaDeCompras int
	diaDoVencimento      int
}

type OperacaoFinanceira struct {
	Id        int
	Valor     int
	Descricao string
	Debito    bool
	Fixo      bool
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := createSqliteDatabase()

	horaAtual := time.Now()

	println(horaAtual.Day())

	if db == nil {
		log.Fatal("Erro ao criar banco de dados")
	}
	defer db.Close()
	seedDatabase(db)

	// total de despesas
	total_despesas := selectTotalDeDespesas(db)
	total_receitas := selectTotalDeReceitas(db)

	println("Total de despesas: ", total_despesas)
	println("Total de receitas: ", total_receitas)

}

func selectTotalDeReceitas(db *sql.DB) int {
	rows, err := db.Query("select sum(coalesce(valor, 0)) from operacoes_financeiras where debito = false")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			log.Fatal(err)
		}
	}
	return total
}

func selectTotalDeDespesas(db *sql.DB) int {
	rows, err := db.Query("select sum(coalesce(valor, 0)) from operacoes_financeiras where debito = true")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			log.Fatal(err)
		}
	}
	return total

}

func checkCartaoExists(db *sql.DB, name string) bool {
	rows, err := db.Query("select * from cartoes where name = ?", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

func checkGastoFixoExists(db *sql.DB, gastoFixo OperacaoFinanceira) bool {
	rows, err := db.Query("select * from operacoes_financeiras where descricao = ? and valor = ?", gastoFixo.Descricao, gastoFixo.Valor)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

func createSqliteDatabase() *sql.DB {

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	doStatement(db, `
		CREATE TABLE IF NOT EXISTS cartoes (id integer not null primary key, name text, primeiroDiaDeCompras integer not null, diaDoVencimento integer not null);
	`)

	doStatement(db, `
		CREATE TABLE IF NOT EXISTS operacoes_financeiras (id integer not null primary key, valor integer not null, descricao text, debito boolean not null default true, fixo boolean not null default false);
	`)

	return db
}

func doStatement(db *sql.DB, sqlStmt string) {
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

func seedDatabase(db *sql.DB) {
	createOrUpdateByCartao(db, Cartao{Name: "vinicius", PrimeiroDiaDeCompras: 13, diaDoVencimento: 21})
	createOrUpdateByCartao(db, Cartao{Name: "franciele", PrimeiroDiaDeCompras: 1, diaDoVencimento: 6})

	//gastos fixos
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "condominio", Valor: 320, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "financiamento casa", Valor: 2150, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "nina diarista", Valor: 560, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "internet", Valor: 100, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "luz", Valor: 120, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "agua", Valor: 200, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "cartao 2", Valor: 560, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "financiamento luz solar", Valor: 490, Fixo: true, Debito: true})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "emprestimo mamãe s2", Valor: 330, Fixo: true, Debito: true})

	//salarios
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "salario vinicius", Valor: 4368, Fixo: true, Debito: false})
	createOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "adiantamento vinicius", Valor: 3068, Fixo: true, Debito: false})
}

func createOrUpdateByCartao(db *sql.DB, cartao Cartao) {
	if checkCartaoExists(db, cartao.Name) {
		_, err := db.Exec("update cartoes set primeiroDiaDeCompras = ?, diaDoVencimento = ? where name = ?", cartao.PrimeiroDiaDeCompras, cartao.diaDoVencimento, cartao.Name)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := db.Exec("insert into cartoes(name, primeiroDiaDeCompras, diaDoVencimento) values(?, ?, ?)", cartao.Name, cartao.PrimeiroDiaDeCompras, cartao.diaDoVencimento)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createOrUpdateByOperacaoFinanceira(db *sql.DB, gastoFixo OperacaoFinanceira) {
	if checkGastoFixoExists(db, gastoFixo) {
		_, err := db.Exec("update operacoes_financeiras set valor = ?, fixo = ?, debito = ? where descricao = ?", gastoFixo.Valor, gastoFixo.Fixo, gastoFixo.Debito, gastoFixo.Descricao)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := db.Exec("insert into operacoes_financeiras(valor, descricao, fixo, debito) values(?,?,?,?)", gastoFixo.Valor, gastoFixo.Descricao, gastoFixo.Fixo, gastoFixo.Debito)
		if err != nil {
			log.Fatal(err)
		}
	}
}
