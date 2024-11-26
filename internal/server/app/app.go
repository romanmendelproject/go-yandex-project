package app

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/romanmendelproject/go-yandex-project/internal/server/config"
	"github.com/romanmendelproject/go-yandex-project/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-project/internal/server/jwt"
	"github.com/romanmendelproject/go-yandex-project/internal/server/middlewares/logger"
	"github.com/romanmendelproject/go-yandex-project/internal/server/router"
	"github.com/romanmendelproject/go-yandex-project/internal/server/storage"
	"github.com/romanmendelproject/go-yandex-project/internal/server/user"

	"github.com/pressly/goose/v3"
	_ "github.com/romanmendelproject/go-yandex-project/internal/server/storage/migrations"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.ReadConfig()
	logger.SetLogLevel(cfg.App.LogLevel)
	var handler *handlers.ServiceHandlers

	database := dbInit(ctx, cfg)
	defer database.Close()

	token := jwt.NewJWT()
	userData := user.NewUserAuth(database, token)
	handler = handlers.NewHandlers(cfg, database, userData)

	r := router.NewRouter(handler, cfg, token)
	func() {
		err := http.ListenAndServe(cfg.DB.DBIP, r)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func dbInit(ctx context.Context, cfg config.Config) *storage.PostgresStorage {
	// ps := "postgres://username:userpassword@localhost:5432/dbname"

	database := storage.NewPostgresStorage(ctx, cfg.DB.DNDSN)

	db, err := sql.Open("postgres", cfg.DB.DNDSN)
	if err != nil {
		log.Error("Failed to open DB", "error", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./internal/server/storage/migrations"); err != nil {
		log.Error("Failed to run migrations", "error", err)
	}
	return database
}
