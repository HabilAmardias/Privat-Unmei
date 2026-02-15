package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/db"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/redis"
	"privat-unmei/internal/upgrader"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func Run() {
	// add production environment option
	var isProd bool = os.Getenv("ENVIRONMENT_OPTION") == constants.Production
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	upg := upgrader.CreateUpgrader()

	zl, err := logger.CreateNewLogger(isProd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	rc := redis.ConnectRedis()
	driver, err := db.ConnectDB(zl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer driver.Close()

	app := gin.New()
	app.ContextWithFallback = true

	Bootstrap(driver, rc, zl, app, upg)

	port := ":" + os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr:    port,
		Handler: app.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zl.Infoln("Shutdown Server....")
	timeoutEnv := os.Getenv("GRACEFUL_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutEnv)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zl.Infoln("Server Shutdown:", err)
	}

	<-ctx.Done()
	zl.Infof("timeout of %d seconds.\n", timeout)
	zl.Infoln("Server exiting")
}
