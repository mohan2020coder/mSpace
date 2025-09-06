// cmd/api/main.go
// package main

// import (
// 	"fmt"

// 	"github.com/mohan2020coder/mSpace/internal/api"
// 	"github.com/mohan2020coder/mSpace/internal/config"
// 	"github.com/mohan2020coder/mSpace/internal/logger"
// )

// func main() {
// 	cfg := config.LoadConfig("./config.yaml")
// 	log := logger.NewLogger(cfg.Logging.Level)

//		r := api.SetupRouter(cfg, log)
//		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
//		log.Sugar().Infof("🚀 API running at %s", addr)
//		r.Run(addr)
//	}
//
// cmd/api/main.go
package main

import (
	"fmt"
	"log"

	"github.com/mohan2020coder/mSpace/internal/api"
	"github.com/mohan2020coder/mSpace/internal/config"
	"github.com/mohan2020coder/mSpace/internal/db"
	"github.com/mohan2020coder/mSpace/internal/logger"
	"github.com/mohan2020coder/mSpace/internal/models"
	"github.com/mohan2020coder/mSpace/internal/search"
	"github.com/mohan2020coder/mSpace/internal/storage"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zl, err := logger.New(cfg.Logging.Level)
	if err != nil {
		log.Fatalf("failed to construct logger: %v", err)
	}
	defer zl.Sync()
	zsugar := zl.Sugar()
	zsugar.Infof("starting with config: %+v", cfg.Server)

	// Init DB
	gdb := db.Init(cfg.Database.DSN)

	// AutoMigrate models
	if err := gdb.AutoMigrate(
		&models.Community{},
		&models.Collection{},
		&models.Item{},
		&models.Metadata{},
	); err != nil {
		zl.Fatal("AutoMigrate failed", zap.Error(err))
	}

	// Init MinIO
	minioClient := storage.NewMinio(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey, cfg.Storage.Bucket, cfg.Storage.SSL)

	app := &api.App{
		Cfg:    cfg,
		DB:     gdb,
		Minio:  minioClient,
		Logger: zl,
	}

	index, err := search.NewIndex("./bleve_index")
	if err != nil {
		zl.Fatal("failed to init search index", zap.Error(err))
	}
	
	r := api.SetupRouter(app,index)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	zsugar.Infof("API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		zl.Fatal("failed to run server", zap.Error(err))
	}
}
