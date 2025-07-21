package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AmiyoKm/green_light/internal/auth"
	"github.com/AmiyoKm/green_light/internal/env"
	"github.com/AmiyoKm/green_light/internal/jsonlog"
	"github.com/AmiyoKm/green_light/internal/mailer"
	"github.com/AmiyoKm/green_light/internal/store"
	"github.com/AmiyoKm/green_light/internal/vcs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

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
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
	jwt struct {
		secret string
		exp    time.Duration
		iss    string
	}
}
type application struct {
	config        config
	logger        *jsonlog.Logger
	store         store.Storage
	mailer        mailer.Mailer
	wg            sync.WaitGroup
	authenticator auth.Authenticator
}

func main() {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	if err := godotenv.Load(); err != nil {
		logger.PrintInfo("No .env file found, using environment variables", nil)
	}

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", env.GetString("ENV", "development"), "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", env.GetString("PROD_DB_DSN", ""), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limit requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limit burst size")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiting")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "33d6aac1de496d", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "cc527b84e26fd8", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight<no-reply@greenlight.net>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", env.GetString("JWT_SECRET", ""), "JWT SECRET")
	flag.StringVar(&cfg.jwt.iss, "jwt-iss", env.GetString("JWT_ISS", "greenlight"), "JWT SECRET")
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	portStr := env.GetString("PORT", "")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			cfg.port = port
		}
	}

	cfg.jwt.exp = 24 * time.Hour

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("Database connection pool established", nil)

	storage := store.NewStorage(db)

	expvar.NewString(version).Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config:        cfg,
		logger:        logger,
		store:         storage,
		authenticator: auth.NewJWTAuthenticator(cfg.jwt.secret, cfg.jwt.iss, cfg.jwt.iss),
		mailer:        mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
