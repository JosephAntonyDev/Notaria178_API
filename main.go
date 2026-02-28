package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core"
	userInfra "github.com/JosephAntonyDev/Notaria178_API/internal/user/infra"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Error fatal: JWT_SECRET no está configurado en el archivo .env")
	}

	db, err := core.GetDBPool()
	if err != nil {
		log.Fatalf("Error fatal al conectar con la base de datos: %v", err)
	}
	defer db.Close() 

	r := gin.Default()

	r.Use(core.SetupCORS())

	userInfra.SetupDependencies(r, db, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Notaría 178 iniciado en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}