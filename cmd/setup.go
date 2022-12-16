package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup() {
	r := gin.New()
	v1 := r.Group("/api/v1")

	httpServer := http.Server{
		Addr:        ":3000",
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	go func() {
		log.Printf("Server Starting on Port: %v", 3000)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("ERROR STARTING SERVER: %v", "PORT")
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Printf("CLOSING SERVER, SIGNAL %v GOTTEN", sig)

	ctx := context.Background()
	httpServer.Shutdown(ctx)
}
