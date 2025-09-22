package biz

import (
	"ragx/app/pkg/cache/redis"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// 采用依赖倒置的方式，在biz定义data接口，防止循环依赖
type Data interface {
	// pg的db连接
	DB() *gorm.DB
	Rdb() *redis.Client
}

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewChatUsecase,
	NewKnowledgeBaseUsecase,
	NewKnowledgeDocumentUsecase,
)
