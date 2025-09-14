package routers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/handlers"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/middlewares"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/repositories/cache"
	"privat-unmei/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RouteConfig struct {
	App                   *gin.Engine
	StudentHandler        *handlers.StudentHandlerImpl
	AdminHandler          *handlers.AdminHandlerImpl
	CourseCategoryHandler *handlers.CourseCategoryHandlerImpl
	CourseHandler         *handlers.CourseHandlerImpl
	MentorHandler         *handlers.MentorHandlerImpl
	CourseRatingHandler   *handlers.CourseRatingHandlerImpl
	CourseRequestHandler  *handlers.CourseRequestHandlerImpl
	ChatHandler           *handlers.ChatHandlerImpl
	PaymentHandler        *handlers.PaymentHandlerImpl
	DiscountHandler       *handlers.DiscountHandlerImpl
	AdditionalCostHandler *handlers.AdditionalCostHandlerImpl
	RBACRepository        *repositories.RBACRepository
	RBACCacheRepository   *cache.RBACCacheRepository
	TokenUtil             *utils.JWTUtil
	Logger                logger.CustomLogger
}

func (c *RouteConfig) Setup() {
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)
	prometheus.MustRegister(httpRequestDuration, httpRequestsTotal)

	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}
	c.App.Use(cors.New(config))
	c.App.Use(middlewares.PrometheusMiddleware(httpRequestsTotal, httpRequestDuration))
	c.App.Use(middlewares.LoggerMiddleware(c.Logger))
	c.App.Use(middlewares.ErrorMiddleware(c.Logger))
	c.App.GET("/metrics", gin.WrapH(promhttp.Handler()))

	c.SetupPublicRoute()
	c.SetupPrivateRoute()
	c.SetupWebsocketRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	v1 := c.App.Group("/api/v1")
	v1.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dtos.Response{
			Success: true,
			Data:    "Welcome to Privat Unmei API",
		})
	})
	v1.POST("/register", c.StudentHandler.Register)
	v1.POST("/login", c.StudentHandler.Login)
	v1.GET("/verify", middlewares.AuthenticationMiddleware(c.TokenUtil, constants.ForVerification), c.StudentHandler.Verify)
	v1.POST("/reset-password/send", c.StudentHandler.SendResetTokenEmail)
	v1.POST("/reset-password/reset", middlewares.AuthenticationMiddleware(c.TokenUtil, constants.ForReset), c.StudentHandler.ResetPassword)
	v1.POST("/admin/login", c.AdminHandler.Login)
	v1.POST("/mentor/login", c.MentorHandler.Login)
	v1.GET("/course-categories", c.CourseCategoryHandler.GetCategoriesList)
	v1.GET("/courses/most-bought", c.CourseHandler.MostBoughtCourses)
	v1.GET("/auth/google", c.StudentHandler.GoogleLogin)
	v1.GET("/auth/google/callback", c.StudentHandler.GoogleLoginCallback)
	v1.GET("/courses", c.CourseHandler.ListCourse)
	v1.GET("/courses/:id", c.CourseHandler.CourseDetail)
	v1.GET("/mentors/:id", c.MentorHandler.GetMentorProfileForStudent)
	v1.GET("/courses/:id/reviews", c.CourseRatingHandler.GetCourseReview)
}

func (c *RouteConfig) SetupPrivateRoute() {
	v1 := c.App.Group("/api/v1")
	v1.Use(middlewares.AuthenticationMiddleware(c.TokenUtil, constants.ForLogin))
	v1.POST("/courses/:id/reviews", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseRatingResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRatingHandler.AddReview)
	v1.GET("/verify/send", c.StudentHandler.SendVerificationEmail)
	v1.PATCH("/courses/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseHandler.UpdateCourse)
	v1.GET("/me", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.StudentResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.StudentHandler.GetStudentProfile)
	v1.POST("/me/change-password", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.StudentResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.StudentHandler.ChangePassword)
	v1.GET("/students", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.StudentResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.StudentHandler.GetStudentList)
	v1.GET("/mentors", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.GetMentorList)
	v1.POST("/mentors", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.AddNewMentor)
	v1.GET("/mentors/password", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdminHandler.GenerateRandomPassword)
	v1.PATCH("/mentors/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.UpdateMentorForAdmin)
	v1.DELETE("/mentors/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.DeleteMentor)
	v1.POST("/course-categories", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseCategoryResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseCategoryHandler.CreateCategory)
	v1.PATCH("/course-categories/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.CourseCategoryResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseCategoryHandler.UpdateCategory)
	v1.DELETE("/course-categories/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.CourseCategoryResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseCategoryHandler.DeleteCategory)
	v1.POST("/mentors/me/change-password", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.ChangePassword)
	v1.PATCH("/students/me", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.StudentResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.StudentHandler.UpdateStudentProfile)
	v1.POST("/courses", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseHandler.AddNewCourse)
	v1.DELETE("/courses/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseHandler.DeleteCourse)
	v1.GET("/mentors/me/courses", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseHandler.MentorListCourse)
	v1.PATCH("/mentors/me", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.UpdateMentor)
	v1.POST("/courses/:id/course-requests", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.CreateReservation)
	v1.GET("/course-requests/:id/approve", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.AcceptCourseRequest)
	v1.GET("/course-requests/:id/reject", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.RejectCourseRequest)
	v1.GET("/course-requests/:id/confirm-payment", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.ConfirmPayment)
	v1.GET("/course-requests/:id/payment-detail", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.PaymentDetailResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.GetPaymentDetail)
	v1.GET("/mentors/me", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.MentorResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.MentorHandler.GetProfileForMentor)
	v1.GET("/mentors/me/course-requests", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.MentorCourseRequestList)
	v1.GET("/mentors/me/course-requests/:id", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.MentorCourseRequestDetail)
	v1.GET("/me/course-requests", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.StudentCourseRequestList)
	v1.GET("/me/course-requests/:id", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.CourseRequestHandler.StudentCourseRequestDetail)
	v1.POST("/chatrooms", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.ChatroomResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.ChatHandler.CreateChatroom)
	v1.GET("/chatrooms/me", c.ChatHandler.GetUserChatrooms)
	v1.GET("/chatrooms/:id", c.ChatHandler.GetChatroom)
	v1.GET("/chatrooms/:id/messages", c.ChatHandler.GetMessages)
	v1.POST("/chatrooms/:id/messages", c.ChatHandler.SendMessage)
	v1.GET("/courses/:id/mentor-availability", c.MentorHandler.GetDOWAvailability)
	v1.POST("/payment-methods", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.PaymentMethodResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.PaymentHandler.CreatePaymentMethod)
	v1.DELETE("/payment-methods/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.PaymentMethodResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.PaymentHandler.DeletePaymentMethod)
	v1.PATCH("/payment-methods/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.PaymentMethodResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.PaymentHandler.UpdatePaymentMethod)
	v1.GET("/payment-methods", c.PaymentHandler.GetAllPaymentMethod)
	v1.GET("/mentors/:id/payment-methods", c.PaymentHandler.GetMentorPaymentMethod)
	v1.POST("/discounts", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.DiscountResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.DiscountHandler.CreateNewDiscount)
	v1.PATCH("/discounts/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.DiscountResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.DiscountHandler.UpdateDiscountAmount)
	v1.DELETE("/discounts/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.DiscountResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.DiscountHandler.DeleteDiscount)
	v1.GET("/discounts", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.DiscountResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.DiscountHandler.GetAllDiscount)
	v1.POST("/additional-costs", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.AdditionalCostResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdditionalCostHandler.CreateNewAdditionalCost)
	v1.PATCH("/additional-costs/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.AdditionalCostResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdditionalCostHandler.UpdateCostAmount)
	v1.DELETE("/additional-costs/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.AdditionalCostResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdditionalCostHandler.DeleteCost)
	v1.GET("/additional-costs", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.AdditionalCostResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdditionalCostHandler.GetAllAdditionalCost)
	v1.POST("/admins/me/verify", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.AdminResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdminHandler.VerifyAdmin)
	v1.POST("/admins/me/change-password", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.AdminResource,
		c.RBACRepository,
		c.RBACCacheRepository,
		c.Logger,
	), c.AdminHandler.ChangePassword)
}

func (c *RouteConfig) SetupWebsocketRoute() {
	r := c.App.Group("/ws/v1")
	r.Use(middlewares.WSAuthenticationMiddleware(c.TokenUtil, constants.ForLogin))
	r.GET("/chatrooms/:id/messages", c.ChatHandler.ConnectChatChannel)
}
