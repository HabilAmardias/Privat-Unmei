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
	courseAvailabilityRepo := repositories.CreateCourseAvailabilityRepository(db)
	courseRepo := repositories.CreateCourseRepository(db)
	topicRepo := repositories.CreateTopicRepository(db)
	courseRequestRepo := repositories.CreateCourseRequestRepository(db)
	courseRatingRepo := repositories.CreateCourseRatingRepository(db)

	bcryptUtil := utils.CreateBcryptUtil()
	gomailUtil := utils.CreateGomailUtil()
	cloudinaryUtil := utils.CreateCloudinaryUtil()
	jwtUtil := utils.CreateJWTUtil()
	googleUtil := utils.CreateGoogleUtil()

	mentorService := services.CreateMentorService(transactionManager, userRepo, mentorRepo, topicRepo, courseCategoryRepo, courseAvailabilityRepo, courseRepo, bcryptUtil, jwtUtil, cloudinaryUtil, gomailUtil)
	adminService := services.CreateAdminService(userRepo, adminRepo, studentRepo, mentorRepo, transactionManager, cloudinaryUtil, bcryptUtil, jwtUtil, gomailUtil)
	studentService := services.CreateStudentService(userRepo, studentRepo, transactionManager, bcryptUtil, gomailUtil, cloudinaryUtil, jwtUtil, googleUtil)
	courseCategoryService := services.CreateCourseCategoryService(courseCategoryRepo, transactionManager)
	courseService := services.CreateCourseService(courseAvailabilityRepo, courseRepo, courseCategoryRepo, topicRepo, transactionManager, courseRequestRepo)
	courseRatingService := services.CreateCourseRatingService(courseRepo, courseRatingRepo, courseRequestRepo, mentorRepo, transactionManager)

	studentHandler := handlers.CreateStudentHandler(studentService)
	adminHandler := handlers.CreateAdminHandler(adminService)
	mentorHandler := handlers.CreateMentorHandler(mentorService)
	courseCategoryHandler := handlers.CreateCourseCategoryHandler(courseCategoryService)
	courseHandler := handlers.CreateCourseHandler(courseService)
	courseRatingHandler := handlers.CreateCourseRatingHandler(courseRatingService)

	cfg := routers.RouteConfig{
		App:                   app,
		StudentHandler:        studentHandler,
		AdminHandler:          adminHandler,
		CourseCategoryHandler: courseCategoryHandler,
		CourseHandler:         courseHandler,
		MentorHandler:         mentorHandler,
		CourseRatingHandler:   courseRatingHandler,
		RBACRepository:        rbacRepo,
		TokenUtil:             jwtUtil,
	}
	cfg.Setup()
}
