package main

import (
	"godindin/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/is_it_on", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "0.0.1",
		})
	})
	r.GET("/total_despesas", func(c *gin.Context) {
		total_despesas := database.SelectTotalDeDespesas(db)
		c.JSON(http.StatusOK, gin.H{
			"total_despesas": total_despesas,
		})
	})

	r.GET("/total_receitas", func(c *gin.Context) {
		total_receitas := database.SelectTotalDeReceitas(db)
		c.JSON(http.StatusOK, gin.H{
			"total_receitas": total_receitas,
		})
	})

	r.GET("/cartoes", func(c *gin.Context) {
		cartoes := database.SelectCartoes(db)
		c.JSON(http.StatusOK, gin.H{
			"cartoes": cartoes,
		})
	})

	r.POST("/cartao", func(c *gin.Context) {
		var cartao database.Cartao
		c.BindJSON(&cartao)
		if database.CheckCartaoExists(db, cartao.Name) {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Cartão já existe",
			})
		} else {
			database.InsertCartao(db, cartao)
			c.JSON(http.StatusCreated, gin.H{
				"message": "Cartão criado com sucesso",
			})
		}
	})

	r.POST("/gasto_fixo", func(c *gin.Context) {
		var gastoFixo database.OperacaoFinanceira
		c.BindJSON(&gastoFixo)
		if database.CheckGastoFixoExists(db, gastoFixo) {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Gasto fixo já existe",
			})
		} else {
			database.InsertOperacaoFinanceira(db, gastoFixo)
			c.JSON(http.StatusCreated, gin.H{
				"message": "Gasto fixo criado com sucesso",
			})
		}
	})

	r.Run()

	// total de despesas
	total_despesas := database.SelectTotalDeDespesas(db)
	total_receitas := database.SelectTotalDeReceitas(db)

	println("Total de despesas: ", total_despesas)
	println("Total de receitas: ", total_receitas)
	println("Hoje" + horaAtual.String())
}
