package main

import (
	"godindin/database"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db := database.CreateSqliteDatabase()

	horaAtual := time.Now()

	println(horaAtual.Day())

	if db == nil {
		log.Fatal("Erro ao criar banco de dados")
	}
	defer db.Close()
	database.SeedDatabase(db)

	// total de despesas
	total_despesas := database.SelectTotalDeDespesas(db)
	total_receitas := database.SelectTotalDeReceitas(db)

	println("Total de despesas: ", total_despesas)
	println("Total de receitas: ", total_receitas)
	println("Hoje" + horaAtual.String())
}
