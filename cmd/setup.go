package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"tosinjs/cloud-backup/cmd/api/v1/routes"
	mySqlFileRepo "tosinjs/cloud-backup/internal/repository/fileRepo/mySqlRepo"
	mySqlUserRepo "tosinjs/cloud-backup/internal/repository/userRepo/mySqlRepo"
	"tosinjs/cloud-backup/internal/service/authService"
	"tosinjs/cloud-backup/internal/service/awsService"
	"tosinjs/cloud-backup/internal/service/cryptoService"
	"tosinjs/cloud-backup/internal/service/fileService"
	"tosinjs/cloud-backup/internal/service/userService"
	"tosinjs/cloud-backup/internal/service/validationService"
	"tosinjs/cloud-backup/internal/setup/aws"
	"tosinjs/cloud-backup/internal/setup/database/mySql"
	utils "tosinjs/cloud-backup/utils/config"

	"github.com/gin-gonic/gin"
)

func Setup() {
	config, err := utils.LoadConfig("./", "app", "env")
	if err != nil {
		log.Fatalf("Error GETIING CONFIG: %v", err)
		os.Exit(1)
	}

	r := gin.New()
	v1 := r.Group("/api/v1")

	connFactory, err := mySql.NewMySQLServer(config.DSN)
	if err != nil {
		log.Fatalf("MYSQL SETUP ERROR %v", err)
		os.Exit(1)
	}

	defer connFactory.Close()
	conn := connFactory.GetConn()

	if err := conn.Ping(); err != nil {
		log.Fatalf("MYSQL SETUP ERROR %v", err)
		os.Exit(1)
	}

	s3, err := aws.NewS3Service(
		config.AWS_ID, config.AWS_SECRET,
		config.AWS_TOKEN, config.AWS_REGION,
	)
	if err != nil {
		log.Fatalf("AWS SETUP ERROR %v", err)
		os.Exit(1)
	}

	//Repository Setups
	userRepo := mySqlUserRepo.New(conn)
	fileRepo := mySqlFileRepo.New(conn)

	//Service Setup
	awsSVC := awsService.New(s3)
	fileSVC := fileService.New(awsSVC, fileRepo)
	cryptoSVC := cryptoService.New()
	validationSVC := validationService.New()
	authSVC := authService.New(config.JWTSECRET)
	userSVC := userService.New(userRepo, cryptoSVC, authSVC)

	httpServer := http.Server{
		Addr:        fmt.Sprintf(":%v", config.PORT),
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}

	//SETUP ROUTES

	//File Routes
	routes.FileRoutes(v1, fileSVC, authSVC)

	//User Routes
	routes.UserRoutes(v1, userSVC, validationSVC)

	//Ping Route
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "pong")
	})

	go func() {
		log.Printf("Server Starting on Port: %v", config.PORT)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("ERROR STARTING SERVER: %v", config.PORT)
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
