package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupAttendanceRoutes(
	r *gin.Engine,
	checkInOutCtrl *controllers.CheckInOutController,
	getMyAttendancesCtrl *controllers.GetMyAttendancesController,
	getEmployeeAttendancesCtrl *controllers.GetEmployeeAttendancesController,
	jwtSecret string,
) {
	api := r.Group("/api/v1/attendance")
	
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/check", checkInOutCtrl.Handle)
		protected.GET("/history", getMyAttendancesCtrl.Handle)

		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
		{
			adminOnly.GET("/history/:id", getEmployeeAttendancesCtrl.Handle)
		}
	}
}