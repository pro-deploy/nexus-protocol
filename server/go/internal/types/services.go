package types

import (
	"github.com/go-redis/redis/v8"
	"github.com/nexus-protocol/server/internal/ai"
	"github.com/nexus-protocol/server/internal/analytics"
	"github.com/nexus-protocol/server/internal/auth"
	"github.com/nexus-protocol/server/internal/batch"
	"github.com/nexus-protocol/server/internal/conversation"
	"github.com/nexus-protocol/server/internal/db"
	"github.com/nexus-protocol/server/internal/webhook"
	"github.com/nexus-protocol/server/pkg/config"
	"go.uber.org/zap"
)

// Services contains all service dependencies
type Services struct {
	Config       *config.Config
	Logger       *zap.Logger
	Database     *db.DB
	Redis        *redis.Client
	AIService    *ai.Service
	Auth         *auth.Service
	Templates    *ai.Service // Templates service is part of AI service
	Batch        *batch.Service
	Webhooks     *webhook.Service
	Analytics    *analytics.Service
	Conversations *conversation.Service
}
