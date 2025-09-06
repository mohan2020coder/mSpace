// cmd/api/main.go
package main

import (
	"fmt"

	"github.com/mohan2020coder/mSpace/internal/api"
	"github.com/mohan2020coder/mSpace/internal/config"
	"github.com/mohan2020coder/mSpace/internal/logger"
)

func main() {
	cfg := config.LoadConfig("./config.yaml")
	log := logger.NewLogger(cfg.Logging.Level)

	r := api.SetupRouter(cfg, log)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Sugar().Infof("🚀 API running at %s", addr)
	r.Run(addr)
}
