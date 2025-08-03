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
	mentorRepo := repositories.CreateMentorRepositoryImpl(db)
	transactionManager := repositories.CreateTransactionManager(db)
	rbacRepo := repositories.CreateRBACRepository(db)
	courseCategoryRepo := repositories.CreateCourseCategoryRepository(db)

	bcryptUtil := utils.CreateBcryptUtil()
	gomailUtil := utils.CreateGomailUtil()
	cloudinaryUtil := utils.CreateCloudinaryUtil()
	jwtUtil := utils.CreateJWTUtil()

	mentorService := services.CreateMentorService(transactionManager, userRepo, mentorRepo, bcryptUtil, jwtUtil, cloudinaryUtil, gomailUtil)
	adminService := services.CreateAdminService(userRepo, adminRepo, studentRepo, mentorRepo, transactionManager, cloudinaryUtil, bcryptUtil, jwtUtil, gomailUtil)
	studentService := services.CreateStudentService(userRepo, studentRepo, transactionManager, bcryptUtil, gomailUtil, cloudinaryUtil, jwtUtil)
	courseCategoryService := services.CreateCourseCategoryService(courseCategoryRepo, transactionManager)

	studentHandler := handlers.CreateStudentHandler(studentService)
	adminHandler := handlers.CreateAdminHandler(adminService)
	mentorHandler := handlers.CreateMentorHandler(mentorService)
	courseCategoryHandler := handlers.CreateCourseCategoryHandler(courseCategoryService)

	cfg := routers.RouteConfig{
		App:                   app,
		StudentHandler:        studentHandler,
		AdminHandler:          adminHandler,
		CourseCategoryHandler: courseCategoryHandler,
		MentorHandler:         mentorHandler,
		RBACRepository:        rbacRepo,
		TokenUtil:             jwtUtil,
	}
	cfg.Setup()
}
