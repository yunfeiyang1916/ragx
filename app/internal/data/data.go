package data

import (
	"ragx/app/internal/biz"
	"ragx/app/internal/biz/entity"
	"ragx/app/internal/conf"
	"ragx/app/internal/data/repo"
	logging "ragx/app/pkg/logger"

	redisHelper "ragx/app/pkg/cache/redis"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	repo.NewKnowledgeBaseRepo,
	repo.NewKnowledgeDocumentRepo,
	repo.NewKnowledgeChunkRepo,
)

// Data .
type data struct {
	db *gorm.DB
	// clickhouse的db连接
	chDB *gorm.DB
	rdb  *redisHelper.Client
}

func (d *data) DB() *gorm.DB {
	return d.db
}
func (d *data) ChDB() *gorm.DB {
	return d.chDB
}
func (d *data) Rdb() *redisHelper.Client { return d.rdb }

// NewData .
func NewData(c *conf.Data, logger log.Logger) (biz.Data, func(), error) {
	logHelper := log.NewHelper(logger)
	cleanup := func() {
		logHelper.Info("closing the data resources")
	}
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		logHelper.Fatalf("Got error when connect database, the error is '%+v'", gerror.Wrap(err, ""))
	}
	db = db.Debug()
	db.Logger = logging.DefaultGormLogger
	if err := db.AutoMigrate(&entity.KnowledgeBase{}, &entity.KnowledgeDocument{}, &entity.KnowledgeChunk{}); err != nil {
		logHelper.Fatalf("Got error when auto migrate database, the error is '%+v'", gerror.Wrap(err, ""))
	}
	return &data{db: db}, cleanup, nil
}
