package config

import (
	"database/sql"
	"privat-unmei/internal/handlers"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/routers"
	"privat-unmei/internal/services"
	"privat-unmei/internal/utils"

	"github.com/gin-gonic/gin"
)

func Bootstrap(db *sql.DB, app *gin.Engine) {
	userRepo := repositories.CreateUserRepository(db)
	studentRepo := repositories.CreateStudentRepository(db)
	adminRepo := repositories.CreateAdminRepository(db)
	transactionManager := repositories.CreateTransactionManager(db)

	bcryptUtil := utils.CreateBcryptUtil()
	gomailUtil := utils.CreateGomailUtil()
	cloudinaryUtil := utils.CreateCloudinaryUtil()
	jwtUtil := utils.CreateJWTUtil()

	adminService := services.CreateAdminService(userRepo, adminRepo, transactionManager, bcryptUtil, jwtUtil, gomailUtil)
	studentService := services.CreateStudentService(userRepo, studentRepo, transactionManager, bcryptUtil, gomailUtil, cloudinaryUtil, jwtUtil)

	studentHandler := handlers.CreateStudentHandler(studentService)
	adminHandler := handlers.CreateAdminHandler(adminService)
	cfg := routers.RouteConfig{
		App:            app,
		StudentHandler: studentHandler,
		AdminHandler:   adminHandler,
		TokenUtil:      jwtUtil,
	}
	cfg.Setup()
}
