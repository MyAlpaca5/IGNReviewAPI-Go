package router

import (
	"reflect"
	"strings"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/controllers"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/middlewares"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/repositories"
	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/pasetotoken"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func New(pool *pgxpool.Pool) *gin.Engine {
	// Default With the Logger and Recovery middleware already attached
	router := gin.Default()

	// --- Set Custom Validator Methods for JSON data ---
	// https://pkg.go.dev/github.com/go-playground/validator/v10#Validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// use the names which have been specified for JSON representations of structs, rather than normal Go field names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	// --- Set Global Middlewares ---
	router.Use(middlewares.RequestRateLimiter(2, 4))
	router.Use(middlewares.BodySizeLimiter(65_536))
	registry := prometheus.NewRegistry()
	router.Use(middlewares.MonitorMetrics(registry))

	// --- Set Error Handlers ---
	// router.NoRoute(notFoundErrorHandler)
	router.HandleMethodNotAllowed = true
	// router.NoMethod(replyWithUnsupportedHTTPMethodError)

	// --- Create Controllers ---
	var healthcheckController = controllers.HealthcheckController{}
	var reviewRepo = repositories.NewReview(pool)
	var reviewController = controllers.ReviewController{Repo: reviewRepo}
	var userRepo = repositories.NewUser(pool)
	var userController = controllers.UserController{Repo: userRepo}
	var tokenRepo = repositories.NewToken(pool)
	tokenMaker, err := pasetotoken.NewPasetoMaker(viper.GetString("secret.symmetric_key"))
	if err != nil {
		panic(err.Error())
	}
	var tokenController = controllers.TokenController{Repo: tokenRepo, UserRepo: userRepo, TokenMaker: tokenMaker}

	// --- Set Routes, Handlers, and per-request Middlewares ---
	router.GET("/healthcheck", healthcheckController.HealthcheckHandler)
	router.POST("/api/tokens/authentication", tokenController.CreateAuthenticationTokenHandler)
	router.POST("/api/users", userController.CreateUserHandler)

	authorized := router.Group("api", middlewares.Authenticate(tokenMaker))
	{
		simpleRole := authorized.Group("", middlewares.Authorize(models.RoleSimple))
		{
			simpleRole.GET("/reviews/:id", reviewController.ShowReviewHandler)
			simpleRole.DELETE("/reviews/:id", reviewController.DeleteReviewHandler)
			simpleRole.PATCH("/reviews/:id", reviewController.UpdateReviewHandler)
			simpleRole.GET("/reviews", reviewController.ListReviewsHandler)
			simpleRole.POST("/reviews", reviewController.CreateReviewHandler)
		}

		adminRole := authorized.Group("admin", middlewares.Authorize(models.RoleAdmin))
		{
			adminRole.GET("/metrics", func(c *gin.Context) {
				promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
			})
		}
	}

	return router
}
