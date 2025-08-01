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
	transactionManager := repositories.CreateTransactionManager(db)

	bcryptUtil := utils.CreateBcryptUtil()
	gomailUtil := utils.CreateGomailUtil()
	cloudinaryUtil := utils.CreateCloudinaryUtil()
	jwtUtil := utils.CreateJWTUtil()

	studentService := services.CreateStudentService(userRepo, studentRepo, transactionManager, bcryptUtil, gomailUtil, cloudinaryUtil, jwtUtil)

	studentHandler := handlers.CreateStudentHandler(studentService)
	cfg := routers.RouteConfig{
		App:            app,
		StudentHandler: studentHandler,
		TokenUtil:      jwtUtil,
	}
	cfg.Setup()
}
