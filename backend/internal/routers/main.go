package routers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/handlers"
	"privat-unmei/internal/middlewares"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	RBACRepository        *repositories.RBACRepository
	TokenUtil             *utils.JWTUtil
}

func (c *RouteConfig) Setup() {
	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}
	c.App.Use(cors.New(config))
	c.App.Use(middlewares.ErrorMiddleware())

	c.SetupPublicRoute()
	c.SetupPrivateRoute()
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
	v1.POST("/verify", c.StudentHandler.Verify)
	v1.POST("/reset-password/send", c.StudentHandler.SendResetTokenEmail)
	v1.POST("/reset-password/reset", c.StudentHandler.ResetPassword)
	v1.POST("/admin/login", c.AdminHandler.Login)
	v1.POST("/mentor/login", c.MentorHandler.Login)
	v1.GET("/courses/categories", c.CourseCategoryHandler.GetCategoriesList)
	v1.GET("/courses/most-bought", c.CourseHandler.MostBoughtCourses)
	v1.GET("/auth/google", c.StudentHandler.GoogleLogin)
	v1.GET("/auth/google/callback", c.StudentHandler.GoogleLoginCallback)
	v1.GET("/courses", c.CourseHandler.ListCourse)
	v1.GET("/courses/:id", c.CourseHandler.CourseDetail)
	v1.GET("/courses/:id/reviews", c.CourseRatingHandler.GetCourseReview)
}

func (c *RouteConfig) SetupPrivateRoute() {
	v1 := c.App.Group("/api/v1")
	v1.Use(middlewares.AuthenticationMiddleware(c.TokenUtil))
	v1.POST("/courses/:id/reviews", c.CourseRatingHandler.AddReview)
	v1.GET("/verify/send", c.StudentHandler.SendVerificationEmail)
	v1.PATCH("/courses/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
	), c.CourseHandler.UpdateCourse)
	v1.GET("/students", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.StudentResource,
		c.RBACRepository,
	), c.StudentHandler.GetStudentList)
	v1.GET("/mentors", middlewares.AuthorizationMiddleware(
		constants.ReadAllPermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.GetMentorList)
	v1.POST("/mentors", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.AddNewMentor)
	v1.GET("/mentors/password", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.AdminHandler.GenerateRandomPassword)
	v1.PATCH("/admin/mentors/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.UpdateMentorForAdmin)
	v1.DELETE("/mentors/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteAllPermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.DeleteMentor)
	v1.POST("/courses/categories", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseCategoryResource,
		c.RBACRepository,
	), c.CourseCategoryHandler.CreateCategory)
	v1.PATCH("/courses/categories/:id", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.CourseCategoryResource,
		c.RBACRepository,
	), c.CourseCategoryHandler.UpdateCategory)
	v1.POST("/mentors/me/change-password", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.ChangePassword)
	v1.PATCH("/students/me", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.StudentResource,
		c.RBACRepository,
	), c.StudentHandler.UpdateStudentProfile)
	v1.POST("/courses", middlewares.AuthorizationMiddleware(
		constants.CreatePermission,
		constants.CourseResource,
		c.RBACRepository,
	), c.CourseHandler.AddNewCourse)
	v1.DELETE("/courses/:id", middlewares.AuthorizationMiddleware(
		constants.DeleteOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
	), c.CourseHandler.DeleteCourse)
	v1.GET("/mentor/courses", middlewares.AuthorizationMiddleware(
		constants.ReadOwnPermission,
		constants.CourseResource,
		c.RBACRepository,
	), c.CourseHandler.MentorListCourse)
	v1.PATCH("/mentors/me", middlewares.AuthorizationMiddleware(
		constants.UpdateOwnPermission,
		constants.MentorResource,
		c.RBACRepository,
	), c.MentorHandler.UpdateMentor)
	v1.POST("/courses/:id/course-requests", c.CourseRequestHandler.CreateReservation)
	v1.GET("/course-requests/:id/approve", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
	), c.CourseRequestHandler.AcceptCourseRequest)
	v1.GET("/course-requests/:id/reject", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
	), c.CourseRequestHandler.RejectCourseRequest)
	v1.GET("/course-requests/:id/confirm-payment", middlewares.AuthorizationMiddleware(
		constants.UpdateAllPermission,
		constants.CourseRequestResource,
		c.RBACRepository,
	), c.CourseRequestHandler.ConfirmPayment)
}
