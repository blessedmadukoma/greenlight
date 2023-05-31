package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/blessedmadukoma/greenlight/internal/data"
	"github.com/blessedmadukoma/greenlight/internal/jsonlog"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

// limiValues retreives the values for the rate limiter from the env
func limitValues() (int, int, bool) {

	rps, err := strconv.Atoi(os.Getenv("LIMITER_RPS"))
	if err != nil {
		log.Fatal("Error retrieving rps value:", err)
	}
	burst, err := strconv.Atoi(os.Getenv("LIMITER_BURST"))
	if err != nil {
		log.Fatal("Error retrieving burst value:", err)
	}
	enabled, err := strconv.ParseBool(os.Getenv("LIMITER_ENABLED"))
	if err != nil {
		log.Fatal("Error retrieving enabled value:", err)
	}

	fmt.Println(rps, burst, enabled)

	return rps, burst, enabled
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// set the default values for the database connection pool i.e. maxOpenConns = 25 open connections, maxIdleConns = 25 open connections, maxIdleTime = 15 minutes
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL maximum open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL maximum idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL maximum idle time")

	rps, burst, enabled := limitValues()

	// set the values for limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", float64(rps), "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", burst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", enabled, "Enable rate limiter")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("Starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of idle connections in the pool.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Set the maximum idle time for a connection.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check that the connection is successful.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
