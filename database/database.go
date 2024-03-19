package database

import (
	"database/sql"
	"log"
)

// SelectTotalDeReceitas retrieves the total sum of revenue from the database.
func SelectTotalDeReceitas(db *sql.DB) int {
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

// SelectTotalDeDespesas retrieves the total sum of expenses from the database.
func SelectTotalDeDespesas(db *sql.DB) int {
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

// CheckCartaoExists checks if a card with the given name exists in the database.
func CheckCartaoExists(db *sql.DB, name string) bool {
	rows, err := db.Query("select * from cartoes where name = ?", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

// CheckGastoFixoExists checks if an expense with the given description and value exists in the database.
func CheckGastoFixoExists(db *sql.DB, gastoFixo OperacaoFinanceira) bool {
	rows, err := db.Query("select * from operacoes_financeiras where descricao = ? and valor = ?", gastoFixo.Descricao, gastoFixo.Valor)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return rows.Next()
}

// CreateSqliteDatabase creates a new SQLite database and initializes the necessary tables.
func CreateSqliteDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	DoStatement(db, `
		CREATE TABLE IF NOT EXISTS cartoes (id integer not null primary key, name text, primeiroDiaDeCompras integer not null, diaDoVencimento integer not null);
	`)

	DoStatement(db, `
		CREATE TABLE IF NOT EXISTS operacoes_financeiras (id integer not null primary key, id_origem int not null default 0, valor integer not null, descricao text, debito boolean not null default true, fixo boolean not null default false);
	`)

	DoStatement(db, `
		CREATE TABLE IF NOT EXISTS contas (id integer not null primary key, valor integer not null, debito boolean not null default true, data_hora datetime not null default current_timestamp);
	`)

	DoStatement(db, `
		CREATE TABLE IF NOT EXISTS resumo (id integer not null primary key, valor integer not null, debito boolean not null default true, data_hora datetime not null default current_timestamp);
	`)

	return db
}

// DoStatement executes the given SQL statement on the database.
func DoStatement(db *sql.DB, sqlStmt string) {
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

// SeedDatabase populates the database with initial data.
func SeedDatabase(db *sql.DB) {
	CreateOrUpdateByCartao(db, Cartao{Name: "vinicius", PrimeiroDiaDeCompras: 13, diaDoVencimento: 21})
	CreateOrUpdateByCartao(db, Cartao{Name: "franciele", PrimeiroDiaDeCompras: 1, diaDoVencimento: 6})

	//gastos fixos
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "condominio", Valor: 320, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "financiamento casa", Valor: 2150, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "nina diarista", Valor: 560, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "internet", Valor: 100, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "luz", Valor: 120, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "agua", Valor: 200, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "cartao 2", Valor: 560, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "financiamento luz solar", Valor: 490, Fixo: true, Debito: true})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "emprestimo mamãe s2", Valor: 330, Fixo: true, Debito: true})

	//salarios
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "salario vinicius", Valor: 4368, Fixo: true, Debito: false})
	CreateOrUpdateByOperacaoFinanceira(db, OperacaoFinanceira{Descricao: "adiantamento vinicius", Valor: 3068, Fixo: true, Debito: false})
}

// CreateOrUpdateByCartao creates or updates a card in the database based on its name.
func CreateOrUpdateByCartao(db *sql.DB, cartao Cartao) {
	if CheckCartaoExists(db, cartao.Name) {
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

// CreateOrUpdateByOperacaoFinanceira creates or updates an expense in the database based on its description and value.
func CreateOrUpdateByOperacaoFinanceira(db *sql.DB, gastoFixo OperacaoFinanceira) {
	if CheckGastoFixoExists(db, gastoFixo) {
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

// InsertCartao inserts a new card into the database.
func InsertCartao(db *sql.DB, cartao Cartao) {
	_, err := db.Exec("insert into cartoes(name, primeiroDiaDeCompras, diaDoVencimento) values(?, ?, ?)", cartao.Name, cartao.PrimeiroDiaDeCompras, cartao.diaDoVencimento)
	if err != nil {
		log.Fatal(err)
	}
}

// InsertOperacaoFinanceira inserts a new expense into the database.
func InsertOperacaoFinanceira(db *sql.DB, gastoFixo OperacaoFinanceira) {
	_, err := db.Exec("insert into operacoes_financeiras(valor, descricao, fixo, debito) values(?,?,?,?)", gastoFixo.Valor, gastoFixo.Descricao, gastoFixo.Fixo, gastoFixo.Debito)
	if err != nil {
		log.Fatal(err)
	}
}

// InsertSaldo inserts a new balance into the database.
func InsertSaldo(db *sql.DB, saldo Saldo) {
	_, err := db.Exec("insert into saldo(data_hora, valor) values(?,?)", saldo.DataHora, saldo.Valor)
	if err != nil {
		log.Fatal(err)
	}
}

func SelectCartoes(db *sql.DB) []Cartao {
	rows, err := db.Query("select * from cartoes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var cartoes []Cartao
	for rows.Next() {
		var cartao Cartao
		err := rows.Scan(&cartao.Id, &cartao.Name, &cartao.PrimeiroDiaDeCompras, &cartao.diaDoVencimento)
		if err != nil {
			log.Fatal(err)
		}
		cartoes = append(cartoes, cartao)
	}
	return cartoes
}
