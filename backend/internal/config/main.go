package config

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/upgrader"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func Run() {
	// add production environment option
	var isProd bool
	flag.BoolVar(&isProd, "release", false, "Run production environemnt")
	flag.Parse()
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	upg := upgrader.CreateUpgrader()

	zl, err := logger.CreateNewLogger(isProd)
	if err != nil {
		log.Fatalln(err.Error())
	}

	driver, err := db.ConnectDB(zl)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer driver.Close()

	app := gin.New()
	app.ContextWithFallback = true

	hub := entity.CreateChatHub()
	go hub.Run()

	Bootstrap(driver, zl, app, upg, hub)

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

	log.Println("Shutdown Server ...")
	timeoutEnv := os.Getenv("GRACEFUL_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutEnv)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Printf("timeout of %d seconds.\n", timeout)
	log.Println("Server exiting")
}
