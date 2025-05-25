package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"minitask1.go/internal/routes"
	"minitask1.go/pkg"
)

// @title 			Tickitz
// @version 		1.0
// @description		API for Tickitz APP
// @host			localhost:8080
// @BasePath		/

func main() {
	pg, err := pkg.Connect()
	if err != nil {
		log.Printf("[ERROR] Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		log.Println("Closing DB...")
		pg.Close()
	}()

	rdb := pkg.RedisConnect()

	router := routes.InitRouter(pg, rdb)

	router.Static("/public/img", "./public/img")
	// jalankan service
	router.Run("127.0.0.1:8080")
	// router.Run(":8080")
}
