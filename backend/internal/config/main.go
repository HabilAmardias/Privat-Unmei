package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"privat-unmei/internal/db"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func Run() {
	driver, err := db.ConnectDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer driver.Close()

	app := gin.New()
	app.ContextWithFallback = true
	Bootstrap(driver, app)

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
