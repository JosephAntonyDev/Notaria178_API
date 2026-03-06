package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	actInfra "github.com/JosephAntonyDev/Notaria178_API/internal/act/infra"
	attendanceInfra "github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra"
	auditInfra "github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra"
	branchInfra "github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra"
	clientInfra "github.com/JosephAntonyDev/Notaria178_API/internal/client/infra"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	dashboardInfra "github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra"
	documentInfra "github.com/JosephAntonyDev/Notaria178_API/internal/document/infra"
	"github.com/JosephAntonyDev/Notaria178_API/internal/integration/adapters"
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

	// ── Redis (opcional) ────────────────────────────────────────────────
	var cachePort cache.CachePort
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB := 0
		if v := os.Getenv("REDIS_DB"); v != "" {
			redisDB, _ = strconv.Atoi(v)
		}
		rc, err := cache.NewRedisCache(redisAddr, redisPassword, redisDB)
		if err != nil {
			log.Printf("Advertencia: Redis no disponible, continuando sin caché: %v", err)
		} else {
			defer rc.Close()
			cachePort = rc
		}
	} else {
		log.Println("REDIS_ADDR no configurado, el servidor iniciará sin caché Redis")
	}

	r := gin.Default()

	r.Use(core.SetupCORS())

	userInfra.SetupDependencies(r, db, jwtSecret)
	attendanceInfra.SetupDependencies(r, db, jwtSecret)
	actInfra.SetupDependencies(r, db, jwtSecret, cachePort)
	clientInfra.SetupDependencies(r, db, jwtSecret)
	branchInfra.SetupDependencies(r, db, jwtSecret)
	documentInfra.SetupDependencies(r, db, jwtSecret)
	dashboardInfra.SetupDependencies(r, db, jwtSecret, cachePort)

	// Módulos que exponen sus use cases para integración cruzada
	logActionUC := auditInfra.SetupDependencies(r, db, jwtSecret)
	createNotifUC := notificationInfra.SetupDependencies(r, db, jwtSecret)

	// Adaptadores que cumplen las interfaces de work/domain/events
	auditAdapter := adapters.NewAuditLoggerAdapter(logActionUC)
	notifAdapter := adapters.NewNotifierAdapter(createNotifUC)

	workInfra.SetupDependencies(r, db, jwtSecret, auditAdapter, notifAdapter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Notaría 178 iniciado en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
