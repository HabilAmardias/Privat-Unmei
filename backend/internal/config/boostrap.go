package config

import (
	"privat-unmei/internal/db"
	"privat-unmei/internal/handlers"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/repositories/cache"
	"privat-unmei/internal/routers"
	"privat-unmei/internal/services"
	"privat-unmei/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func Bootstrap(db *db.CustomDB, rc *redis.Client, logger logger.CustomLogger, app *gin.Engine, upg *websocket.Upgrader) {
	userRepo := repositories.CreateUserRepository(db)
	studentRepo := repositories.CreateStudentRepository(db)
	adminRepo := repositories.CreateAdminRepository(db)
	mentorRepo := repositories.CreateMentorRepositoryImpl(db)
	transactionManager := repositories.CreateTransactionManager(db, logger)
	rbacRepo := repositories.CreateRBACRepository(db)
	courseCategoryRepo := repositories.CreateCourseCategoryRepository(db)
	mentorAvailabilityRepo := repositories.CreateCourseAvailabilityRepository(db)
	courseRepo := repositories.CreateCourseRepository(db)
	topicRepo := repositories.CreateTopicRepository(db)
	courseRequestRepo := repositories.CreateCourseRequestRepository(db)
	courseRatingRepo := repositories.CreateCourseRatingRepository(db)
	courseScheduleRepo := repositories.CreateCourseScheduleRepository(db)
	chatRepo := repositories.CreateChatRepository(db)
	paymentRepo := repositories.CreatePaymentRepository(db)
	discountRepo := repositories.CreateDiscountRepository(db)
	additionalCostRepo := repositories.CreateAdditionalCostRepository(db)
	rbaccache := cache.CreateRBACCache(rc)

	bcryptUtil := utils.CreateBcryptUtil()
	gomailUtil := utils.CreateGomailUtil()
	cloudinaryUtil := utils.CreateCloudinaryUtil()
	jwtUtil := utils.CreateJWTUtil()
	googleUtil := utils.CreateGoogleUtil()

	mentorService := services.CreateMentorService(transactionManager, userRepo, mentorRepo, topicRepo, courseCategoryRepo, mentorAvailabilityRepo, courseRequestRepo, courseRepo, paymentRepo, adminRepo, bcryptUtil, jwtUtil, cloudinaryUtil, gomailUtil, logger)
	adminService := services.CreateAdminService(userRepo, adminRepo, studentRepo, mentorRepo, transactionManager, cloudinaryUtil, bcryptUtil, jwtUtil, gomailUtil)
	studentService := services.CreateStudentService(userRepo, studentRepo, adminRepo, transactionManager, bcryptUtil, gomailUtil, cloudinaryUtil, jwtUtil, googleUtil)
	courseCategoryService := services.CreateCourseCategoryService(userRepo, adminRepo, courseCategoryRepo, transactionManager)
	courseService := services.CreateCourseService(courseRepo, courseCategoryRepo, topicRepo, mentorRepo, transactionManager, courseRequestRepo)
	courseRatingService := services.CreateCourseRatingService(courseRepo, courseRatingRepo, courseRequestRepo, mentorRepo, transactionManager)
	courseRequestService := services.CreateCourseRequestService(courseRequestRepo, courseRepo, courseScheduleRepo, mentorAvailabilityRepo, userRepo, studentRepo, mentorRepo, paymentRepo, discountRepo, additionalCostRepo, transactionManager, gomailUtil)
	chatService := services.CreateChatService(chatRepo, userRepo, studentRepo, mentorRepo, transactionManager)
	paymentService := services.CreatePaymentService(paymentRepo, adminRepo, userRepo, mentorRepo, transactionManager)
	discountService := services.CreateDiscountService(discountRepo, userRepo, adminRepo, transactionManager)
	additionalCostService := services.CreateAdditionalCostService(additionalCostRepo, userRepo, adminRepo, transactionManager)

	studentHandler := handlers.CreateStudentHandler(studentService)
	adminHandler := handlers.CreateAdminHandler(adminService)
	mentorHandler := handlers.CreateMentorHandler(mentorService, courseService)
	courseCategoryHandler := handlers.CreateCourseCategoryHandler(courseCategoryService)
	courseHandler := handlers.CreateCourseHandler(courseService)
	courseRatingHandler := handlers.CreateCourseRatingHandler(courseRatingService)
	courseRequestHandler := handlers.CreateCourseRequestHandler(courseRequestService)
	chatHandler := handlers.CreateChatHandler(chatService, upg, rc, logger)
	paymentHandler := handlers.CreatePaymentHandler(paymentService)
	discountHandler := handlers.CreateDiscountHandler(discountService)
	additionalCostHandler := handlers.CreateAdditionalCostHandler(additionalCostService)

	cfg := routers.RouteConfig{
		App:                   app,
		StudentHandler:        studentHandler,
		AdminHandler:          adminHandler,
		CourseCategoryHandler: courseCategoryHandler,
		CourseHandler:         courseHandler,
		MentorHandler:         mentorHandler,
		CourseRatingHandler:   courseRatingHandler,
		CourseRequestHandler:  courseRequestHandler,
		ChatHandler:           chatHandler,
		PaymentHandler:        paymentHandler,
		DiscountHandler:       discountHandler,
		AdditionalCostHandler: additionalCostHandler,
		RBACRepository:        rbacRepo,
		RBACCacheRepository:   rbaccache,
		TokenUtil:             jwtUtil,
		Logger:                logger,
	}
	cfg.Setup()
}
