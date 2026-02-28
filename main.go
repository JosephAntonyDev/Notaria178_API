package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	actInfra "github.com/JosephAntonyDev/Notaria178_API/internal/act/infra"
	attendanceInfra "github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra"
	auditInfra "github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra"
	branchInfra "github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra"
	clientInfra "github.com/JosephAntonyDev/Notaria178_API/internal/client/infra"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core"
	documentInfra "github.com/JosephAntonyDev/Notaria178_API/internal/document/infra"
	notificationInfra "github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra"
	userInfra "github.com/JosephAntonyDev/Notaria178_API/internal/user/infra"
	workInfra "github.com/JosephAntonyDev/Notaria178_API/internal/work/infra"
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
	attendanceInfra.SetupDependencies(r, db, jwtSecret)
	actInfra.SetupDependencies(r, db, jwtSecret)
	clientInfra.SetupDependencies(r, db, jwtSecret)
	branchInfra.SetupDependencies(r, db, jwtSecret)
	workInfra.SetupDependencies(r, db, jwtSecret)
	documentInfra.SetupDependencies(r, db, jwtSecret)
	notificationInfra.SetupDependencies(r, db, jwtSecret)
	_ = auditInfra.SetupDependencies(r, db, jwtSecret) // Devuelve *LogActionUseCase para inyección futura
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Notaría 178 iniciado en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
