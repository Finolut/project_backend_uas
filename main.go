package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	mongorepo "clean-arch/app/repository/mongo"
	pgrepo "clean-arch/app/repository/postgre"
	service "clean-arch/app/service"
	config "clean-arch/config"
	db "clean-arch/database"
	route "clean-arch/route"

	"go.mongodb.org/mongo-driver/mongo"
)

// @title Auth API v1
// @version 1.0
// @description API untuk autentikasi dan manajemen user menggunakan Clean Architecture (Support MongoDB & PostgreSQL)
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	// Load env & config
	if err := config.LoadEnv(); err != nil {
		log.Printf("warning loading env: %v", err)
	}
	conf := config.Get()

	// Init logger
	logWriter, stdLogger := config.InitLogger(conf.LogPath)
	_ = stdLogger // use if you want

	// Create Fiber app
	app := config.NewFiberApp(logWriter)

	// DB driver selection: "mongo", "postgres", or "both"
	driver := conf.DBDriver
	if driver == "" {
		driver = "mongo"
	}

	// Holders
	var pgDB *sql.DB
	var mongoClient *mongo.Client
	var mongoDB *mongo.Database

	// Connect to Postgres if required or available
	if driver == "postgres" || driver == "both" {
		// Use provided DSN or build one from environment variables
		psqlDsn := conf.PostgresDsn
		if psqlDsn == "" {
			// try to build from env vars (fallback)
			host := os.Getenv("PG_HOST")
			if host == "" {
				host = "localhost"
			}
			port := os.Getenv("PG_PORT")
			if port == "" {
				port = "5432"
			}
			user := os.Getenv("PG_USER")
			if user == "" {
				user = "postgres"
			}
			pass := os.Getenv("PG_PASSWORD")
			dbname := os.Getenv("PG_DB")
			if dbname == "" {
				dbname = "alumni_db2" // your DB name as you said
			}
			psqlDsn = "postgres://"
			// format: postgres://user:pass@host:port/dbname?sslmode=disable
			if pass != "" {
				psqlDsn = psqlDsn + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
			} else {
				psqlDsn = psqlDsn + user + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
			}
		}

		var err error
		pgDB, err = db.ConnectPostgres(psqlDsn)
		if err != nil {
			log.Fatalf("failed connect to postgres: %v", err)
		}
		log.Println("connected to postgres")
	}

	// Connect to Mongo if required or available
	if driver == "mongo" || driver == "both" {
		// Determine mongo uri and db name
		mongoURI := conf.MongoURI
		if mongoURI == "" {
			mongoURI = os.Getenv("MONGO_URI")
			if mongoURI == "" {
				mongoURI = "mongodb://localhost:27017"
			}
		}
		mongoDBName := os.Getenv("MONGO_DB")
		if mongoDBName == "" {
			mongoDBName = "alumni_db" // default name if not set
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var err error
		mongoClient, mongoDB, err = db.ConnectMongo(ctx, mongoURI, mongoDBName)
		if err != nil {
			log.Fatalf("failed connect to mongo: %v", err)
		}
		log.Printf("connected to mongo (db=%s)", mongoDBName)
	}

	// Build repositories (only if DB connections exist)
	var userRepo pgrepo.UserRepository
	var studentRepo pgrepo.StudentRepository
	var lecturerRepo pgrepo.LecturerRepository
	var roleRepo pgrepo.RoleRepository
	var permissionRepo pgrepo.PermissionRepository
	var rolePermRepo pgrepo.RolePermissionRepository
	var achRefRepo pgrepo.AchievementRefRepository
	var achRepo mongorepo.AchievementRepository

	if pgDB != nil {
		userRepo = pgrepo.NewUserRepository(pgDB)
		studentRepo = pgrepo.NewStudentRepository(pgDB)
		lecturerRepo = pgrepo.NewLecturerRepository(pgDB)
		roleRepo = pgrepo.NewRoleRepository(pgDB)
		permissionRepo = pgrepo.NewPermissionRepository(pgDB)
		rolePermRepo = pgrepo.NewRolePermissionRepository(pgDB)
		achRefRepo = pgrepo.NewAchievementRefRepository(pgDB)
	}

	if mongoDB != nil {
		achRepo = mongorepo.NewAchievementRepository(mongoDB, "achievements")
	}

	// Build service repos struct
	repos := &service.Repos{
		UserRepo:           userRepo,
		RoleRepo:           roleRepo,
		PermissionRepo:     permissionRepo,
		RolePermissionRepo: rolePermRepo,
		StudentRepo:        studentRepo,
		LecturerRepo:       lecturerRepo,
		AchievementRefRepo: achRefRepo,
		AchievementRepo:    achRepo,
	}

	// Create services
	services := service.NewServices(pgDB, mongoDB, repos)

	// Register routes (assumes route.RegisterRoutes accepts app and services)
	// You may need to adapt if your route.RegisterRoutes signature is different.
	route.RegisterRoutes(app, services)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server with graceful shutdown
	port := conf.AppPort
	if port == "" {
		port = "3000"
	}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server starting on :%s (driver=%s)\n", port, driver)
		serverErr <- app.Listen(":" + port)
	}()

	// Wait for interrupt or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatalf("server error: %v", err)
	case sig := <-quit:
		log.Printf("signal %v received, shutting down...", sig)
		// Graceful shutdown sequence
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.Shutdown(); err != nil {
			log.Printf("error during app shutdown: %v", err)
		}

		// Close Mongo
		if mongoClient != nil {
			_ = mongoClient.Disconnect(ctxShutdown)
			log.Println("mongo client disconnected")
		}
		// Close Postgres
		if pgDB != nil {
			_ = pgDB.Close()
			log.Println("postgres closed")
		}
	}
}
