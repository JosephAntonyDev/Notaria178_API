package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	attendanceRepo := repository.NewPostgresAttendanceRepository(db)
	
	checkInOutUseCase := app.NewCheckInOutUseCase(attendanceRepo)
	getMyAttendancesUseCase := app.NewGetMyAttendancesUseCase(attendanceRepo)
	getEmployeeAttendancesUseCase := app.NewGetEmployeeAttendancesUseCase(attendanceRepo)

	checkInOutCtrl := controllers.NewCheckInOutController(checkInOutUseCase)
	getMyAttendancesCtrl := controllers.NewGetMyAttendancesController(getMyAttendancesUseCase)
	getEmployeeAttendancesCtrl := controllers.NewGetEmployeeAttendancesController(getEmployeeAttendancesUseCase)

	routes.SetupAttendanceRoutes(r, checkInOutCtrl, getMyAttendancesCtrl, getEmployeeAttendancesCtrl, jwtSecret)
}
