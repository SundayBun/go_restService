package middleware

import (
	"go.uber.org/zap"
	"goWebService/config"
)

type MiddlewareManager struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, logger *zap.SugaredLogger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, logger: logger}
}
