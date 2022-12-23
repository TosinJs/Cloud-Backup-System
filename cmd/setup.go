package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"tosinjs/cloud-backup/cmd/api/v1/routes"
	"tosinjs/cloud-backup/internal/service/awsService"
	"tosinjs/cloud-backup/internal/service/fileService"
	"tosinjs/cloud-backup/internal/setup/aws"
)

func Setup() {
	r := gin.New()
	v1 := r.Group("/api/v1")

	s3, err := aws.NewS3Service()

	if err != nil {
		log.Fatalf("AWS SETUP ERROR %v", err)
		os.Exit(1)
	}

	//Service Setup
	awsSVC := awsService.New(s3)
	fileSVC := fileService.New(awsSVC)

	httpServer := http.Server{
		Addr:        ":3000",
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}

	//SETUP ROUTES

	//File Routes
	routes.FileRoutes(v1, fileSVC)

	//Ping Route
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
